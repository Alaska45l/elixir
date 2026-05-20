import { env } from '$env/dynamic/public';
import { getProducts } from '$lib/api/client';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ fetch }) => {
  const products = await getProducts(fetch);
  const urls = ['/', '/fragrances', '/envios', '/contacto', ...products.map((p) => `/fragrances/${p.slug}`)];
  const site = env.PUBLIC_SITE_URL ?? 'http://localhost:5173';
  const body = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls.map((url) => `<url><loc>${site}${url}</loc></url>`).join('')}</urlset>`;
  return new Response(body, { headers: { 'Content-Type': 'application/xml' } });
};
