import { writable } from 'svelte/store';

export type ToastMessage = { id: number; text: string; tone: 'ok' | 'error' };
const store = writable<ToastMessage[]>([]);
let nextID = 1;

export const toast = {
  subscribe: store.subscribe,
  push: (text: string, tone: ToastMessage['tone'] = 'ok') => {
    const id = nextID++;
    store.update((items) => [...items, { id, text, tone }]);
    setTimeout(() => store.update((items) => items.filter((x) => x.id !== id)), 3800);
  }
};
