import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/public';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, request, url }) => {
  if (url.pathname === '/admin/login') return {};
  const cookie = request.headers.get('cookie') ?? '';
  const base = env.PUBLIC_API_URL ?? 'http://localhost:8080';
  const res = await fetch(`${base}/api/admin/me`, {
    headers: { cookie },
    credentials: 'include'
  });
  if (!res.ok) {
    redirect(302, '/admin/login');
  }
  return {};
};
