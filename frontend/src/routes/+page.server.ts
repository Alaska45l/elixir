import { getProducts } from '$lib/api/client';
import type { Product } from '$lib/api/client';
import type { PageServerLoad } from './$types';

const FAMILIES = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
const COLLECTION_LIMIT = 4;

type CollectionData = {
  family: string;
  images: string[];
  href: string;
};

export const load: PageServerLoad = async ({ fetch, parent }) => {
  const { homepage } = await parent();
  const featured = await getProducts(fetch, '?featured=true&limit=3');
  const collections = await buildCollections(fetch);

  return { homepage, featured: featured.slice(0, 3), collections };
};

async function buildCollections(fetcher: typeof fetch): Promise<CollectionData[]> {
  const shuffled = [...FAMILIES].sort(() => Math.random() - 0.5);
  const result: CollectionData[] = [];

  for (const family of shuffled) {
    if (result.length >= COLLECTION_LIMIT) break;

    const products = await getProducts(fetcher, `?family=${encodeURIComponent(family)}&limit=6`);
    const withImages = products.filter((product: Product) => {
      return product.scent_family === family && product.images?.some((image) => image.url);
    });

    if (withImages.length === 0) continue;

    const images = withImages
      .map((product) => {
        const primary = product.images.find((image) => image.is_primary && image.url);
        return primary?.url ?? product.images.find((image) => image.url)?.url ?? '';
      })
      .filter(Boolean);

    result.push({
      family,
      images,
      href: `/fragrances?family=${encodeURIComponent(family)}`
    });
  }

  return result;
}
