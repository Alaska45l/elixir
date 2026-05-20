import { getHomepage } from '$lib/api/client';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
  return { homepage: await getHomepage(fetch) };
};
