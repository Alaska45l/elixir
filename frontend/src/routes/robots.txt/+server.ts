import { PUBLIC_SITE_URL } from '$env/static/public';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async () => new Response(`User-agent: *
Disallow: /admin
Disallow: /api/
Sitemap: ${PUBLIC_SITE_URL || 'http://localhost:5173'}/sitemap.xml
`, { headers: { 'Content-Type': 'text/plain' } });
