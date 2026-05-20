import { env } from '$env/dynamic/public';

export type Variant = {
  id: string;
  product_id: string;
  size_ml: number;
  price_ars_cents: number;
  stock: number;
  sku?: string;
  active: boolean;
};

export type ProductImage = {
  id: string;
  product_id: string;
  url: string;
  alt_text?: string;
  is_primary: boolean;
  sort_order: number;
};

export type Product = {
  id: string;
  slug: string;
  name: string;
  tagline?: string;
  description?: string;
  scent_family?: string;
  gender_tag?: string;
  concentration?: string;
  top_notes: string[];
  heart_notes: string[];
  base_notes: string[];
  featured: boolean;
  active: boolean;
  display_order: number;
  variants: Variant[];
  images: ProductImage[];
  min_price_ars_cents: number;
  total_stock: number;
};

export type HomepageSettings = {
  hero_heading: string;
  hero_subheading: string;
  hero_image_url: string;
  hero_cta_label: string;
  hero_cta_url: string;
  editorial_heading: string;
  editorial_body: string;
  editorial_image_url: string;
};

export type CartLine = { variant_id: string; quantity: number; unit_price_ars_cents?: number };
export type ValidatedItem = {
  product_id: string;
  variant_id: string;
  product_name: string;
  slug: string;
  primary_image?: string;
  size_ml: number;
  quantity: number;
  unit_price_ars_cents: number;
  subtotal_ars_cents: number;
  available_stock: number;
  corrected_price: boolean;
};
export type CartValidation = { valid: boolean; items: ValidatedItem[]; subtotal_ars_cents: number; errors: string[] };
export type DiscountValidation = { valid: boolean; code?: string; discount_ars_cents: number; message?: string };
export type Order = {
  id: string;
  external_reference: string;
  status: 'pending' | 'paid' | 'failed' | 'cancelled' | 'shipped' | 'delivered';
  customer_name: string;
  customer_email: string;
  shipping_cost_ars_cents: number;
  subtotal_ars_cents: number;
  total_ars_cents: number;
  discount_ars_cents: number;
  items?: ValidatedItem[];
};
export type ShippingZone = {
  id: string;
  zone_name: string;
  province_codes: string[];
  base_cost_cents: number;
  estimated_days_min: number;
  estimated_days_max: number;
};

export type ListResponse<T> = { items: T[]; total?: number; limit?: number; offset?: number };

export async function apiFetch<T>(path: string, init?: RequestInit, fetcher: typeof fetch = fetch): Promise<T> {
  const base = env.PUBLIC_API_URL ?? 'http://localhost:8080';
  const res = await fetcher(`${base}${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
    ...init
  });
  if (!res.ok) {
    throw new Error((await res.text()) || 'Error de API');
  }
  return (await res.json()) as T;
}

export async function getProducts(fetcher: typeof fetch = fetch, query = ''): Promise<Product[]> {
  return (await getProductsResponse(fetcher, query)).items;
}

export async function getProductsResponse(fetcher: typeof fetch = fetch, query = ''): Promise<ListResponse<Product>> {
  try {
    const data = await apiFetch<ListResponse<Product>>(`/api/products${query}`, undefined, fetcher);
    return data;
  } catch {
    return { items: demoProducts, total: demoProducts.length, limit: demoProducts.length, offset: 0 };
  }
}

export async function getProduct(slug: string, fetcher: typeof fetch = fetch): Promise<Product | null> {
  try {
    return await apiFetch<Product>(`/api/products/${slug}`, undefined, fetcher);
  } catch {
    return demoProducts.find((p) => p.slug === slug) ?? null;
  }
}

export async function getHomepage(fetcher: typeof fetch = fetch): Promise<HomepageSettings> {
  try {
    const data = await apiFetch<HomepageSettings>('/api/homepage', undefined, fetcher);
    return data.hero_heading ? data : defaultHomepage;
  } catch {
    return defaultHomepage;
  }
}

export async function getShippingZones(fetcher: typeof fetch = fetch): Promise<ShippingZone[]> {
  try {
    const data = await apiFetch<ListResponse<ShippingZone>>('/api/shipping/zones', undefined, fetcher);
    return data.items;
  } catch {
    return [
      { id: 'caba', zone_name: 'CABA', province_codes: ['CF'], base_cost_cents: 0, estimated_days_min: 0, estimated_days_max: 3 },
      { id: 'gba', zone_name: 'Gran Buenos Aires', province_codes: ['BA'], base_cost_cents: 250000, estimated_days_min: 1, estimated_days_max: 4 },
      { id: 'interior', zone_name: 'Interior', province_codes: [], base_cost_cents: 420000, estimated_days_min: 3, estimated_days_max: 7 }
    ];
  }
}

export const defaultHomepage: HomepageSettings = {
  hero_heading: 'Perfumería argentina de gesto privado',
  hero_subheading: 'Fragancias intensas, precisas y comerciales, curadas para noches largas, hoteles silenciosos y piel con presencia.',
  hero_image_url: 'https://images.unsplash.com/photo-1619994403073-2cec844b8e63?auto=format&fit=crop&w=1200&q=85',
  hero_cta_label: 'Catálogo',
  hero_cta_url: '/fragrances',
  editorial_heading: 'Una firma de baja voz',
  editorial_body: 'ELIXIR Exclusive trabaja familias olfativas densas y modernas: maderas limpias, ámbar seco, flores oscuras y cítricos fríos. Cada compra se prepara con empaque sobrio y seguimiento personalizado.',
  editorial_image_url: 'https://images.unsplash.com/photo-1595425970377-c9703cf48b6f?auto=format&fit=crop&w=1000&q=85'
};

export const demoProducts: Product[] = [
  makeProduct('nocturno-oud', 'Nocturno Oud', 'Oud seco, rosa negra y cuero limpio', 'Amaderado', 8900000, 4, 'https://images.unsplash.com/photo-1592945403244-b3fbafd7f539?auto=format&fit=crop&w=900&q=85'),
  makeProduct('ambar-de-recoleta', 'Ámbar de Recoleta', 'Ámbar cálido, vainilla sobria y incienso', 'Oriental', 7600000, 12, 'https://images.unsplash.com/photo-1615634260167-c8cdede054de?auto=format&fit=crop&w=900&q=85'),
  makeProduct('flor-de-noche', 'Flor de Noche', 'Jazmín oscuro, iris y almizcle limpio', 'Floral', 8200000, 3, 'https://images.unsplash.com/photo-1547887538-e3a2f32cb1cc?auto=format&fit=crop&w=900&q=85'),
  makeProduct('citrino-frio', 'Citrino Frío', 'Bergamota helada, neroli y cedro blanco', 'Cítrico', 6900000, 8, 'https://images.unsplash.com/photo-1600612253971-422e7f7faeb6?auto=format&fit=crop&w=900&q=85'),
  makeProduct('gourmand-reserva', 'Gourmand Reserva', 'Tonka, cacao amargo y sándalo', 'Gourmand', 9300000, 2, 'https://images.unsplash.com/photo-1594035910387-fea47794261f?auto=format&fit=crop&w=900&q=85'),
  makeProduct('fresco-sur', 'Fresco Sur', 'Mate verde, pomelo y vetiver', 'Fresco', 7100000, 0, 'https://images.unsplash.com/photo-1585386959984-a4155223168f?auto=format&fit=crop&w=900&q=85')
];

function makeProduct(slug: string, name: string, tagline: string, family: string, price: number, stock: number, image: string): Product {
  const id = slug;
  return {
    id,
    slug,
    name,
    tagline,
    description: `${tagline}. Una composición de alta permanencia pensada para piel y clima urbano argentino.`,
    scent_family: family,
    gender_tag: 'Unisex',
    concentration: 'EDP',
    top_notes: ['Bergamota', 'Pimienta rosa', 'Azafrán'],
    heart_notes: ['Rosa', 'Iris', 'Incienso'],
    base_notes: ['Sándalo', 'Ámbar', 'Almizcle'],
    featured: true,
    active: true,
    display_order: 0,
    variants: [
      { id: `${slug}-50`, product_id: id, size_ml: 50, price_ars_cents: price, stock, sku: `${slug.toUpperCase()}-50`, active: true },
      { id: `${slug}-100`, product_id: id, size_ml: 100, price_ars_cents: Math.round(price * 1.65), stock, sku: `${slug.toUpperCase()}-100`, active: true }
    ],
    images: [
      { id: `${slug}-1`, product_id: id, url: image, alt_text: name, is_primary: true, sort_order: 0 },
      { id: `${slug}-2`, product_id: id, url: image.replace('w=900', 'w=901'), alt_text: `${name} detalle`, is_primary: false, sort_order: 1 }
    ],
    min_price_ars_cents: price,
    total_stock: stock
  };
}
