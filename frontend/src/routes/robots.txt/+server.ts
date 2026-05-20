import { env } from '$env/dynamic/public';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async () => new Response(`User-agent: *
Disallow: /admin
Disallow: /api/
Sitemap: ${env.PUBLIC_SITE_URL ?? 'http://localhost:5173'}/sitemap.xml
`, { headers: { 'Content-Type': 'text/plain' } });
