import { env } from '$env/dynamic/public';

export type Variant = {
  id: string;
  product_id: string;
  size_ml: number;
  price_ars_cents: number;
  stock: number;
  sku?: string;
  active: boolean;
  weight_grams: number;
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
  per_kg_cents: number;
  estimated_days_min: number;
  estimated_days_max: number;
  active: boolean;
};

export type FAQItem = { question: string; answer: string };
export type NavItem = { label: string; href: string };
export type SiteSettings = {
  footer_instagram_url: string;
  footer_tiktok_url: string;
  footer_whatsapp_url: string;
  announcement_bar_text: string;
  announcement_bar_active: boolean;
  about_title: string;
  about_description: string;
  about_location: string;
  about_phone: string;
  faq_items: FAQItem[];
  return_policy_html: string;
  navbar_product_categories: NavItem[];
  low_stock_threshold: number;
};

export type ShippingQuoteRequest = {
  destination_postal_code: string;
  province_code: string;
  weight_grams: number;
  dimensions: { length_cm: number; width_cm: number; height_cm: number };
};

export type ShippingQuoteOption = {
  id: string;
  carrier_name: string;
  service_name: string;
  price_cents: number;
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
    const text = await res.text();
    let message = '';
    try {
      const body = JSON.parse(text) as { error?: string };
      message = body.error ?? '';
    } catch {
      message = text;
    }
    throw new Error(message || 'No se pudo completar la operación');
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

export async function getSettings(fetcher: typeof fetch = fetch): Promise<SiteSettings> {
  try {
    const data = await apiFetch<SiteSettings>('/api/settings', undefined, fetcher);
    return { ...defaultSettings, ...data };
  } catch {
    return defaultSettings;
  }
}

export async function getShippingZones(fetcher: typeof fetch = fetch): Promise<ShippingZone[]> {
  try {
    const data = await apiFetch<ListResponse<ShippingZone>>('/api/shipping/zones', undefined, fetcher);
    return data.items;
  } catch {
    return [
      { id: 'caba', zone_name: 'CABA', province_codes: ['CF'], base_cost_cents: 0, per_kg_cents: 0, estimated_days_min: 0, estimated_days_max: 3, active: true },
      { id: 'gba', zone_name: 'Gran Buenos Aires', province_codes: ['BA'], base_cost_cents: 250000, per_kg_cents: 0, estimated_days_min: 1, estimated_days_max: 4, active: true },
      { id: 'interior', zone_name: 'Interior', province_codes: [], base_cost_cents: 420000, per_kg_cents: 0, estimated_days_min: 3, estimated_days_max: 7, active: true }
    ];
  }
}

export async function quoteShipping(req: ShippingQuoteRequest): Promise<ShippingQuoteOption[]> {
  const data = await apiFetch<ListResponse<ShippingQuoteOption>>('/api/shipping/quote', {
    method: 'POST',
    body: JSON.stringify(req)
  });
  return data.items;
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

export const defaultSettings: SiteSettings = {
  footer_instagram_url: 'https://www.instagram.com/',
  footer_tiktok_url: 'https://www.tiktok.com/',
  footer_whatsapp_url: '',
  announcement_bar_text: 'Envíos a todo el país · Empaque discreto · Seguimiento personalizado',
  announcement_bar_active: true,
  about_title: 'ELIXIR Exclusive',
  about_description: 'Perfumería argentina de lujo discreto. Fragancias intensas, envíos nacionales y atención privada.',
  about_location: 'Buenos Aires, Argentina',
  about_phone: '',
  faq_items: [
    { question: '¿Los perfumes son originales?', answer: 'Sí. ELIXIR Exclusive comercializa fragancias seleccionadas y documentadas.' },
    { question: '¿Qué medios de pago aceptan?', answer: 'El checkout opera en ARS mediante MercadoPago.' },
    { question: '¿Hacen envíos?', answer: 'Sí, a CABA, GBA e Interior con seguimiento.' },
    { question: '¿Puedo consultar por WhatsApp?', answer: 'Sí. Recomendamos WhatsApp para asesoramiento rápido.' }
  ],
  return_policy_html: '<p>Los cambios se revisan caso por caso con el producto cerrado, sin uso y dentro de los plazos informados por atención al cliente.</p>',
  navbar_product_categories: [
    { label: 'Fragancias Masculinas', href: '/fragrances?gender=Masculino' },
    { label: 'Fragancias Femeninas', href: '/fragrances?gender=Femenino' },
    { label: 'Línea Oriental', href: '/fragrances?family=Oriental' },
    { label: 'Línea Amaderada', href: '/fragrances?family=Amaderado' }
  ],
  low_stock_threshold: 5
};

export const demoProducts: Product[] = [
  makeProduct('nocturno-oud', 'Nocturno Oud', 'Oud seco, rosa negra y cuero limpio', 'Amaderado', 8900000, 4, 'https://images.unsplash.com/photo-1541643600914-78b084683702?auto=format&fit=crop&w=900&q=90'),
  makeProduct('ambar-de-recoleta', 'Ámbar de Recoleta', 'Ámbar cálido, vainilla sobria y incienso', 'Oriental', 7600000, 12, 'https://images.unsplash.com/photo-1588405748880-12d1d2a59f75?auto=format&fit=crop&w=900&q=90'),
  makeProduct('flor-de-noche', 'Flor de Noche', 'Jazmín oscuro, iris y almizcle limpio', 'Floral', 8200000, 3, 'https://images.unsplash.com/photo-1616604426203-b9baf4ac29d2?auto=format&fit=crop&w=900&q=90'),
  makeProduct('citrino-frio', 'Citrino Frío', 'Bergamota helada, neroli y cedro blanco', 'Cítrico', 6900000, 8, 'https://images.unsplash.com/photo-1609541657971-7a22e8e76219?auto=format&fit=crop&w=900&q=90'),
  makeProduct('gourmand-reserva', 'Gourmand Reserva', 'Tonka, cacao amargo y sándalo', 'Gourmand', 9300000, 2, 'https://images.unsplash.com/photo-1590736969955-71cc94901144?auto=format&fit=crop&w=900&q=90'),
  makeProduct('fresco-sur', 'Fresco Sur', 'Mate verde, pomelo y vetiver', 'Fresco', 7100000, 0, 'https://images.unsplash.com/photo-1615634260167-c8cdede054de?auto=format&fit=crop&w=900&q=90')
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
      { id: `${slug}-50`, product_id: id, size_ml: 50, price_ars_cents: price, stock, sku: `${slug.toUpperCase()}-50`, active: true, weight_grams: 200 },
      { id: `${slug}-100`, product_id: id, size_ml: 100, price_ars_cents: Math.round(price * 1.65), stock, sku: `${slug.toUpperCase()}-100`, active: true, weight_grams: 320 }
    ],
    images: [
      { id: `${slug}-1`, product_id: id, url: image, alt_text: name, is_primary: true, sort_order: 0 },
      { id: `${slug}-2`, product_id: id, url: image.replace('w=900', 'w=901'), alt_text: `${name} detalle`, is_primary: false, sort_order: 1 }
    ],
    min_price_ars_cents: price,
    total_stock: stock
  };
}
