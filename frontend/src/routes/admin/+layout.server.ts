import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
  if (url.pathname === '/admin/login') {
    return {};
  }
  const res = await fetch('/api/admin/me');
  if (res.status === 401 || res.status === 403) {
    redirect(303, '/admin/login');
  }
  if (!res.ok) {
    throw new Error('No se pudo validar la sesión administrativa');
  }
  return {};
};
