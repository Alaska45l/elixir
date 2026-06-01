import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const key = 'elixir_wishlist';
const initial = readWishlist();
const store = writable<string[]>(initial);

function readWishlist(): string[] {
  if (!browser) return [];
  try {
    const parsed = JSON.parse(localStorage.getItem(key) ?? '[]') as unknown;
    return Array.isArray(parsed)
      ? Array.from(new Set(parsed.filter((item): item is string => typeof item === 'string' && item !== '')))
      : [];
  } catch {
    return [];
  }
}

if (browser) {
  store.subscribe((items) => {
    try {
      localStorage.setItem(key, JSON.stringify(Array.from(new Set(items))));
    } catch {
      // Storage can be unavailable in private browsing contexts.
    }
  });
  window.addEventListener('storage', (event) => {
    if (event.key === key) store.set(readWishlist());
  });
}

export const wishlist = {
  subscribe: store.subscribe,
  add: (slug: string) => store.update((items) => slug && !items.includes(slug) ? [...items, slug] : items),
  remove: (slug: string) => store.update((items) => items.filter((x) => x !== slug)),
  toggle: (slug: string) => store.update((items) => !slug ? items : items.includes(slug) ? items.filter((x) => x !== slug) : [...items, slug]),
  clear: () => store.set([])
};
