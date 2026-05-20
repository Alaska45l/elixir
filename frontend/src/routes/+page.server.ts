import { getProducts } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, parent }) => {
  const { homepage } = await parent();
  return { homepage, featured: (await getProducts(fetch, '?featured=true&limit=3')).slice(0, 3) };
};
