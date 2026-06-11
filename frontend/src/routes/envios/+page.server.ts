import { getShippingZones } from '$lib/api/client';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
  try {
    return { zones: await getShippingZones(fetch) };
  } catch {
    error(503, 'No pudimos cargar las opciones de envío');
  }
};
