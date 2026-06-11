import { PUBLIC_SITE_URL } from '$env/static/public';
import { defaultHomepage, defaultSettings, getHomepage, getSettings } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
  const isAdmin = url.pathname.startsWith('/admin');
  const siteUrl = (PUBLIC_SITE_URL || url.origin).replace(/\/$/, '');
  if (isAdmin) {
    return { homepage: defaultHomepage, settings: defaultSettings, isAdmin, siteUrl };
  }

  const [homepage, settings] = await Promise.all([getHomepage(fetch), getSettings(fetch)]);
  return { homepage, settings, isAdmin, siteUrl };
};
