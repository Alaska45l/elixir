import { error } from '@sveltejs/kit';
import { getProduct, getProducts } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, params }) => {
  try {
    const product = await getProduct(params.slug, fetch);
    if (!product) error(404, 'Fragancia no encontrada');
    const related = (await getProducts(fetch, `?family=${encodeURIComponent(product.scent_family ?? '')}&limit=4`)).filter((item) => item.slug !== product.slug).slice(0, 4);
    return { product, related };
  } catch (err) {
    throw err;
  }
};
