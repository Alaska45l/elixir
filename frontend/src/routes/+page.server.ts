import { getProducts } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, parent }) => {
  const { homepage } = await parent();
  const featured = await getProducts(fetch, '?featured=true&limit=3');
  return { homepage, featured: featured.slice(0, 3) };
};
