import { apiFetch } from '$lib/api/client';
import type { Order } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
  try {
    return { order: await apiFetch<Order>(`/api/orders/${params.ref}`, undefined, fetch), ref: params.ref };
  } catch {
    return { order: null, ref: params.ref };
  }
};
