import { env as privateEnv } from '$env/dynamic/private';
import { PUBLIC_API_URL } from '$env/static/public';
import type { RequestHandler } from './$types';

const hopByHopHeaders = new Set([
  'connection',
  'content-encoding',
  'content-length',
  'host',
  'keep-alive',
  'transfer-encoding',
  'upgrade'
]);

const proxy: RequestHandler = async ({ request, url, fetch }) => {
  const backend = (privateEnv.PRIVATE_API_URL || PUBLIC_API_URL || 'http://localhost:8080').replace(/\/$/, '');
  const path = url.pathname.replace(/^\/api\/?/, '');
  const target = `${backend}/api/${path}${url.search}`;
  const headers = new Headers(request.headers);
  for (const header of hopByHopHeaders) headers.delete(header);

  const method = request.method.toUpperCase();
  const body = method === 'GET' || method === 'HEAD' ? undefined : await request.arrayBuffer();
  const upstream = await fetch(target, { method, headers, body });
  const responseHeaders = new Headers(upstream.headers);
  for (const header of hopByHopHeaders) responseHeaders.delete(header);

  return new Response(upstream.body, {
    status: upstream.status,
    statusText: upstream.statusText,
    headers: responseHeaders
  });
};

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const DELETE = proxy;
export const OPTIONS = proxy;
