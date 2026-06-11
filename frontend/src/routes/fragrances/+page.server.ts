import { getProductsResponse } from '$lib/api/client';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
  try {
    const queryParams = new URLSearchParams(url.searchParams);
    queryParams.delete('concentration');
    queryParams.delete('min_price');
    queryParams.delete('max_price');
    if (!queryParams.has('limit')) queryParams.set('limit', '24');
    const query = `?${queryParams.toString()}`;
    const result = await getProductsResponse(fetch, query);
    const limit = Number(queryParams.get('limit') ?? 24);
    const offset = Number(queryParams.get('offset') ?? 0);

    return {
      products: result.items,
      total: result.total ?? result.items.length,
      limit: result.limit ?? limit,
      offset: result.offset ?? offset,
      search: queryParams.toString(),
      loadError: result.error ?? null
    };
  } catch {
    error(503, 'No pudimos cargar el catálogo');
  }
};
