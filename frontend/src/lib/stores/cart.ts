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
const initial: CartItem[] = browser ? JSON.parse(localStorage.getItem(key) ?? '[]') as CartItem[] : [];

function createCart() {
  const store = writable<CartItem[]>(initial);
  if (browser) {
    store.subscribe((items) => localStorage.setItem(key, JSON.stringify(items)));
  }
  return {
    subscribe: store.subscribe,
    add: (item: CartItem) => store.update((items) => {
      const existing = items.find((x) => x.variantId === item.variantId);
      if (existing) {
        return items.map((x) => x.variantId === item.variantId ? { ...x, quantity: x.quantity + item.quantity } : x);
      }
      return [...items, item];
    }),
    setQuantity: (variantId: string, quantity: number) => store.update((items) => items.map((x) => x.variantId === variantId ? { ...x, quantity: Math.max(1, quantity) } : x)),
    remove: (variantId: string) => store.update((items) => items.filter((x) => x.variantId !== variantId)),
    clear: () => store.set([])
  };
}

export const cart = createCart();
export const cartCount = derived(cart, ($cart) => $cart.reduce((sum, item) => sum + item.quantity, 0));
export const cartSubtotal = derived(cart, ($cart) => $cart.reduce((sum, item) => sum + item.quantity * item.unitPriceCents, 0));
