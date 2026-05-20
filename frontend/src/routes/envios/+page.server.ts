import { getShippingZones } from '$lib/api/client';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => ({ zones: await getShippingZones(fetch) });
