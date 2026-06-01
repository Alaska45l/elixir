import { browser } from '$app/environment';
import { derived, writable } from 'svelte/store';

export type CartItem = {
  variantId: string;
  productSlug: string;
  productName: string;
  image: string;
  sizeML: number;
  unitPriceCents: number;
  weightGrams?: number;
  quantity: number;
};

const key = 'elixir_cart';
const initial: CartItem[] = readCart();

function isCartItem(value: unknown): value is CartItem {
  if (!value || typeof value !== 'object') return false;
  const item = value as Partial<CartItem>;
  return typeof item.variantId === 'string'
    && typeof item.productSlug === 'string'
    && typeof item.productName === 'string'
    && typeof item.image === 'string'
    && typeof item.sizeML === 'number'
    && typeof item.unitPriceCents === 'number'
    && typeof item.quantity === 'number';
}

function readCart(): CartItem[] {
  if (!browser) return [];
  try {
    const parsed = JSON.parse(localStorage.getItem(key) ?? '[]') as unknown;
    return Array.isArray(parsed)
      ? parsed.filter(isCartItem).map((item) => ({ ...item, quantity: Math.max(1, Math.min(99, item.quantity)) }))
      : [];
  } catch {
    return [];
  }
}

function createCart() {
  const store = writable<CartItem[]>(initial);
  if (browser) {
    store.subscribe((items) => {
      try {
        localStorage.setItem(key, JSON.stringify(items));
      } catch {
        // Storage can be unavailable in private browsing contexts.
      }
    });
    window.addEventListener('storage', (event) => {
      if (event.key === key) store.set(readCart());
    });
  }
  return {
    subscribe: store.subscribe,
    add: (item: CartItem) => store.update((items) => {
      item.quantity = Math.max(1, Math.min(99, item.quantity));
      const existing = items.find((x) => x.variantId === item.variantId);
      if (existing) {
        return items.map((x) => x.variantId === item.variantId ? { ...x, quantity: Math.min(99, x.quantity + item.quantity) } : x);
      }
      return [...items, item];
    }),
    setQuantity: (variantId: string, quantity: number) => store.update((items) => items.map((x) => x.variantId === variantId ? { ...x, quantity: Math.max(1, Math.min(99, quantity)) } : x)),
    remove: (variantId: string) => store.update((items) => items.filter((x) => x.variantId !== variantId)),
    clear: () => store.set([])
  };
}

export const cart = createCart();
export const cartCount = derived(cart, ($cart) => $cart.reduce((sum, item) => sum + item.quantity, 0));
export const cartSubtotal = derived(cart, ($cart) => $cart.reduce((sum, item) => sum + item.quantity * item.unitPriceCents, 0));
