import { PUBLIC_SITE_URL } from '$env/static/public';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url }) => {
  const siteUrl = (PUBLIC_SITE_URL || url.origin).replace(/\/$/, '');
  return new Response(`User-agent: *
Disallow: /admin
Disallow: /api/
Sitemap: ${siteUrl}/sitemap.xml
`, { headers: { 'Content-Type': 'text/plain' } });
};
