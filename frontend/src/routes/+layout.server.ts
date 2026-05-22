import { getHomepage, getSettings } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
  const [homepage, settings] = await Promise.all([getHomepage(fetch), getSettings(fetch)]);
  return { homepage, settings };
};
