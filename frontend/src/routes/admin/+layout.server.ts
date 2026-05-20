import { redirect } from '@sveltejs/kit';
import { apiFetch } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
  if (url.pathname === '/admin/login') return {};
  try {
    await apiFetch('/api/admin/me', undefined, fetch);
    return {};
  } catch {
    redirect(302, '/admin/login');
  }
};
