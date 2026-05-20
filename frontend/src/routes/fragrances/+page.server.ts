import { getProducts } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
  const query = url.search ? url.search : '?limit=24';
  return { products: await getProducts(fetch, query), search: url.searchParams.toString() };
};
