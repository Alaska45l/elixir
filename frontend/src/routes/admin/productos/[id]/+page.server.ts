import { env } from '$env/dynamic/public';
import type { ProductFormValue } from '$lib/types/product-form';
import type { PageServerLoad } from './$types';

type AdminProductResponse = {
  id: string;
  active: boolean;
  product: ProductFormValue;
};

export const load: PageServerLoad = async ({ fetch, params, request }) => {
  const cookie = request.headers.get('cookie') ?? '';
  const base = env.PUBLIC_API_URL ?? 'http://localhost:8080';
  const res = await fetch(`${base}/api/admin/products/${params.id}`, {
    headers: { cookie },
    credentials: 'include'
  });
  if (!res.ok) return { id: params.id, product: undefined };
  const data = (await res.json()) as AdminProductResponse;
  return { id: data.id, product: data.product };
};
