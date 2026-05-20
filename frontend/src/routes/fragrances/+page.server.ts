import { getProductsResponse } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
  const query = url.search ? url.search : '?limit=24';
  const result = await getProductsResponse(fetch, query);
  return {
    products: result.items,
    total: result.total ?? result.items.length,
    limit: result.limit ?? Number(url.searchParams.get('limit') ?? 24),
    offset: result.offset ?? Number(url.searchParams.get('offset') ?? 0),
    search: url.searchParams.toString()
  };
};
