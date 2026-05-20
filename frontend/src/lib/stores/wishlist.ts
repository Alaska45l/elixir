import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const key = 'elixir_wishlist';
const initial = browser ? JSON.parse(localStorage.getItem(key) ?? '[]') as string[] : [];
const store = writable<string[]>(initial);

if (browser) {
  store.subscribe((items) => localStorage.setItem(key, JSON.stringify(items)));
}

export const wishlist = {
  subscribe: store.subscribe,
  toggle: (slug: string) => store.update((items) => items.includes(slug) ? items.filter((x) => x !== slug) : [...items, slug])
};
