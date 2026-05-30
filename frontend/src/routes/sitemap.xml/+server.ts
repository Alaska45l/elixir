import { PUBLIC_SITE_URL } from '$env/static/public';
import { getProductsResponse } from '$lib/api/client';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ fetch }) => {
  const products = [];
  const limit = 60;
  let offset = 0;
  let total = 0;
  do {
    const page = await getProductsResponse(fetch, `?limit=${limit}&offset=${offset}`);
    products.push(...page.items);
    total = page.total ?? products.length;
    offset += limit;
  } while (products.length < total);
  const urls = ['/', '/fragrances', '/envios', '/contacto', '/preguntas-frecuentes', '/politica-de-devolucion', ...products.map((p) => `/fragrances/${p.slug}`)];
  const site = PUBLIC_SITE_URL || 'http://localhost:5173';
  const body = `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${urls.map((url) => `<url><loc>${site}${url}</loc></url>`).join('')}</urlset>`;
  return new Response(body, { headers: { 'Content-Type': 'application/xml' } });
};
