import { defaultHomepage, defaultSettings, getHomepage, getSettings } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
  const isAdmin = url.pathname.startsWith('/admin');
  if (isAdmin) {
    return { homepage: defaultHomepage, settings: defaultSettings, isAdmin };
  }

  const [homepage, settings] = await Promise.all([getHomepage(fetch), getSettings(fetch)]);
  return { homepage, settings, isAdmin };
};
