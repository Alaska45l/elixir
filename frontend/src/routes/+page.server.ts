import { getProducts } from '$lib/api/client';
import type { Product } from '$lib/api/client';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

const FAMILIES = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
const COLLECTION_LIMIT = 4;
const HERO_IMAGE_LIMIT = 16;

type CollectionData = {
  family: string;
  images: string[];
  href: string;
};

type HeroImageData = {
  url: string;
  alt: string;
};

export const load: PageServerLoad = async ({ fetch, parent }) => {
  try {
    const { homepage } = await parent();
    const [featured, collections, heroImages] = await Promise.all([
      getProducts(fetch, '?featured=true&limit=3'),
      buildCollections(fetch),
      homepage.hero_image_mode === 'product_covers' ? buildHeroImages(fetch) : Promise.resolve([])
    ]);

    return { homepage, featured: featured.slice(0, 3), collections, heroImages };
  } catch {
    error(503, 'No pudimos cargar la página de inicio');
  }
};

async function buildHeroImages(fetcher: typeof fetch): Promise<HeroImageData[]> {
  const products = await getProducts(fetcher, `?limit=${HERO_IMAGE_LIMIT}`);
  const seen = new Set<string>();
  const images: HeroImageData[] = [];

  for (const product of products) {
    const primary = product.images.find((image) => image.is_primary && image.url);
    const image = primary ?? product.images.find((item) => item.url);
    if (!image?.url || seen.has(image.url)) continue;

    seen.add(image.url);
    images.push({ url: image.url, alt: image.alt_text || product.name });
  }

  return images;
}

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
