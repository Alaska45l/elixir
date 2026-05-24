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
  const configuredBase = env.PUBLIC_API_URL?.trim() ?? '';
  const base = path.startsWith('/api/admin') ? '' : configuredBase;
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
  hero_heading: 'Lorem ipsum dolor sit amet',
  hero_subheading: 'Consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  hero_image_url: 'https://images.unsplash.com/photo-1619994403073-2cec844b8e63?auto=format&fit=crop&w=1200&q=85',
  hero_cta_label: 'Catálogo',
  hero_cta_url: '/fragrances',
  editorial_heading: 'Lorem ipsum',
  editorial_body: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  editorial_image_url: 'https://images.unsplash.com/photo-1595425970377-c9703cf48b6f?auto=format&fit=crop&w=1000&q=85'
};

export const defaultSettings: SiteSettings = {
  footer_instagram_url: 'https://www.instagram.com/',
  footer_tiktok_url: 'https://www.tiktok.com/',
  footer_whatsapp_url: '',
  announcement_bar_text: 'Lorem ipsum dolor sit amet · Consectetur adipiscing elit · Sed do eiusmod tempor',
  announcement_bar_active: true,
  about_title: 'ELIXIR Exclusive',
  about_description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
  about_location: 'Buenos Aires, Argentina',
  about_phone: '',
  faq_items: [
    { question: '¿Los perfumes son originales?', answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.' },
    { question: '¿Qué medios de pago aceptan?', answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.' },
    { question: '¿Hacen envíos?', answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.' },
    { question: '¿Puedo consultar por WhatsApp?', answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.' }
  ],
  return_policy_html: '<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>',
  navbar_product_categories: [
    { label: 'Fragancias Unisex', href: '/fragrances?gender=Unisex' },
    { label: 'Fragancias Masculinas', href: '/fragrances?gender=Masculino' },
    { label: 'Fragancias Femeninas', href: '/fragrances?gender=Femenino' },
    { label: 'Línea Oriental', href: '/fragrances?family=Oriental' },
    { label: 'Línea Amaderada', href: '/fragrances?family=Amaderado' }
  ],
  low_stock_threshold: 5
};

type DemoProductSeed = {
  slug: string;
  name: string;
  tagline: string;
  description: string;
  family: string;
  gender: string;
  concentration: string;
  top: string[];
  heart: string[];
  base: string[];
  price: number;
  stock: number;
  image: string;
};

export const demoProducts: Product[] = [
  {
    slug: 'miss-armaf-chic',
    name: 'Miss Armaf Chic',
    tagline: 'Frutas brillantes, cítricos dulces y flores limpias',
    description: 'Apertura de frutilla, frambuesa, pera y cítricos sobre jazmín, peonía y azahar; secado de vainilla, musk, cedro, ambroxan y musgo.',
    family: 'Cítrico',
    gender: 'Femenino',
    concentration: 'EDP',
    top: ['Frutilla', 'Frambuesa', 'Pera', 'Naranja', 'Mandarina', 'Bergamota', 'Calone'],
    heart: ['Jazmín', 'Peonía', 'Azahar'],
    base: ['Patchouli', 'Musk', 'Vainilla', 'Ambroxan', 'Cedro', 'Musgo'],
    price: 2950000,
    stock: 10,
    image: 'https://armaf.com/cdn/shop/files/armaf-1design-2025-11-18T235013.146.png?v=1763491870&width=1080'
  },
  {
    slug: 'miss-armaf-catwalk',
    name: 'Miss Armaf Catwalk',
    tagline: 'Mandarina fresca, durazno floral y musk dulce',
    description: 'Un perfil cítrico floral con salida de mandarina, naranja y cítricos, corazón de durazno, jazmín y lirio, y fondo de vainilla, musk y ambroxan.',
    family: 'Cítrico',
    gender: 'Femenino',
    concentration: 'EDP',
    top: ['Mandarina', 'Naranja', 'Cítricos'],
    heart: ['Durazno', 'Jazmín', 'Lirio'],
    base: ['Vainilla', 'Musk', 'Ambroxan'],
    price: 2950000,
    stock: 10,
    image: 'https://armafperfume.us/cdn/shop/files/9ad438dce22642d7b50a7bd7a9c81743_tplv-omjb5zjo8w-resize-jpeg_800_800.jpg?v=1751506753&width=900'
  },
  {
    slug: 'creme-of-clouds',
    name: 'Creme of Clouds',
    tagline: 'Crema batida, azúcar quemada y vainilla suave',
    description: 'Gourmand cremoso con crema batida, azúcar tostada, leche de coco y vainilla; una estela dulce y lactónica con sensación de caramelo.',
    family: 'Gourmand',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Crema batida', 'Azúcar quemada'],
    heart: ['Leche de coco', 'Vainilla'],
    base: ['Caramelo', 'Musk'],
    price: 2950000,
    stock: 10,
    image: 'https://perfumeoriental.com/cdn/shop/files/creme-of-clouds-fragrance-world-edp-perfume-oriental.webp?v=1771523950&width=900'
  },
  {
    slug: 'lattafa-eclaire',
    name: 'Lattafa Eclaire',
    tagline: 'Caramelo cremoso, leche y vainilla con praliné',
    description: 'Dulce gourmand de caramelo, leche y azúcar con corazón de miel y flores blancas, apoyado en vainilla, praliné y musk.',
    family: 'Gourmand',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Caramelo', 'Leche', 'Azúcar'],
    heart: ['Miel', 'Flores blancas'],
    base: ['Vainilla', 'Praliné', 'Musk'],
    price: 2950000,
    stock: 10,
    image: 'https://www.lattafa-usa.com/cdn/shop/files/Eclaire-1_5803282e-ea5b-4de5-99a5-7d06f5cbae33.png?v=1747415649&width=1200'
  },
  {
    slug: 'odyssey-aqua',
    name: 'Odyssey Aqua',
    tagline: 'Pomelo, naranja y menta sobre una base limpia amaderada',
    description: 'Fragancia fresca de salida cítrica con pomelo, naranja y artemisia; evoluciona hacia lavanda y menta sobre ciprés, patchouli y ambroxan.',
    family: 'Fresco',
    gender: 'Masculino',
    concentration: 'EDP',
    top: ['Pomelo', 'Naranja', 'Artemisia'],
    heart: ['Lavanda', 'Menta'],
    base: ['Ciprés', 'Patchouli', 'Ambroxan'],
    price: 2950000,
    stock: 10,
    image: 'https://armaf.com/cdn/shop/files/image-2023-05-04T112339.859.jpg?v=1739111570&width=1200'
  },
  {
    slug: 'nectar-of-ecstasy-v1',
    name: 'Nectar of Ecstasy (Versión 1)',
    tagline: 'Fruta jugosa, bergamota y vainilla caramelada',
    description: 'Versión dulce y cremosa de Nectar of Ecstasy, con açai, arándano y bergamota sobre flores suaves, vainilla, caramelo y musk.',
    family: 'Gourmand',
    gender: 'Femenino',
    concentration: 'EDP',
    top: ['Açai', 'Arándano', 'Bergamota'],
    heart: ['Fresia', 'Muguet'],
    base: ['Vainilla', 'Caramelo', 'Musk'],
    price: 2950000,
    stock: 10,
    image: 'https://www.french-avenue-parfum.com/wp-content/uploads/2025/01/Eau-de-parfum-Nectar-of-Ecstasy.jpg'
  },
  {
    slug: 'ameer-al-arab-imperium',
    name: 'Ameer Al Arab Imperium',
    tagline: 'Bergamota, jengibre y salvia con fondo amaderado',
    description: 'Woody aromatic con apertura de bergamota, jengibre y salvia, corazón de manzana, cashmeran y geranio, y fondo de musk, ámbar y sándalo.',
    family: 'Amaderado',
    gender: 'Masculino',
    concentration: 'EDP',
    top: ['Bergamota', 'Jengibre', 'Salvia'],
    heart: ['Manzana', 'Cashmeran', 'Geranio'],
    base: ['Musk', 'Ámbar', 'Sándalo'],
    price: 2750000,
    stock: 10,
    image: 'https://fimgs.net/mdimg/perfume-thumbs/375x500.100148.jpg'
  },
  {
    slug: 'spicebomb-night-vision',
    name: 'Spicebomb Night Vision',
    tagline: 'Limón, especias negras y maderas verdes intensas',
    description: 'EDP especiado amaderado con limón y especias negras, corazón verde resinoso con salvia e incienso, y base de abeto balsámico, cedro, patchouli y ládano.',
    family: 'Amaderado',
    gender: 'Masculino',
    concentration: 'EDP',
    top: ['Limón', 'Pimienta negra', 'Chile negro', 'Nuez moscada', 'Clavo'],
    heart: ['Salvia esclarea', 'Resina verde', 'Incienso'],
    base: ['Abeto balsámico', 'Cedro', 'Patchouli', 'Ládano'],
    price: 2750000,
    stock: 10,
    image: 'https://us.viktor-rolf.com/dw/image/v2/AANG_PRD/on/demandware.static/-/Sites-vr-master-catalog/default/dw06340d3b/SB%20NV%20EDP%202024/01_vr_frag_spb_night_vision_edp_perfect_pdp_premium_packshot_90ml_1x1.jpg?q=70&sfrm=jpg&sh=900&sm=cut&sw=900'
  },
  {
    slug: 'nectar-of-ecstasy-v2',
    name: 'Nectar of Ecstasy (Versión 2)',
    tagline: 'Cítricos dulces con firma de cedro y ámbar',
    description: 'Versión amaderada y dulce de Nectar of Ecstasy, con fruta roja, bergamota y flores claras sobre cedro, ámbar y notas dulces.',
    family: 'Amaderado',
    gender: 'Femenino',
    concentration: 'EDP',
    top: ['Açai', 'Arándano', 'Bergamota'],
    heart: ['Fresia', 'Muguet'],
    base: ['Cedro', 'Ámbar', 'Notas dulces'],
    price: 2950000,
    stock: 10,
    image: 'https://www.french-avenue-parfum.com/wp-content/uploads/2025/01/Eau-de-parfum-Nectar-of-Ecstasy.jpg'
  },
  {
    slug: 'odyssey-limoni',
    name: 'Odyssey Limoni',
    tagline: 'Limón dulce, naranja y té con frescura marina',
    description: 'Cítrico fresco con limón, naranja dulce, mandarina y bergamota; corazón de flor de azahar, notas marinas y jengibre, y base de té, musk y ámbar.',
    family: 'Cítrico',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Limón', 'Naranja dulce', 'Mandarina', 'Bergamota'],
    heart: ['Azahar', 'Notas marinas', 'Jengibre'],
    base: ['Té', 'Musk', 'Ámbar'],
    price: 2950000,
    stock: 10,
    image: 'https://fimgs.net/mdimg/perfume-thumbs/375x500.98695.jpg'
  },
  {
    slug: 'odyssey-aqua-v2',
    name: 'Odyssey Aqua (Versión 2)',
    tagline: 'Menta fresca, lavanda y cítricos acuáticos',
    description: 'Segunda lectura de Odyssey Aqua centrada en frescura y menta, con pomelo, naranja, artemisia, lavanda y un cierre de ambroxan, patchouli y ciprés.',
    family: 'Fresco',
    gender: 'Masculino',
    concentration: 'EDP',
    top: ['Pomelo', 'Naranja', 'Artemisia'],
    heart: ['Menta', 'Lavanda'],
    base: ['Ambroxan', 'Patchouli', 'Ciprés'],
    price: 2950000,
    stock: 10,
    image: 'https://armaf.com/cdn/shop/files/IMG_20250710_1558211-ezgif.com-webp-to-jpg-converter.jpg?v=1767894320&width=1200'
  },
  {
    slug: 'paris-corner-taskeen-dia',
    name: 'Paris Corner Taskeen Día',
    tagline: 'Durazno, naranja sanguina y vainilla dulce',
    description: 'Perfil frutal dulce de durazno, naranja sanguina y cardamomo, con corazón de heliotropo, davana, cognac y jazmín, y base de sándalo, vainilla, tonka y patchouli.',
    family: 'Gourmand',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Durazno', 'Naranja sanguina', 'Cardamomo'],
    heart: ['Heliotropo', 'Davana', 'Cognac', 'Jazmín'],
    base: ['Sándalo', 'Benjuí', 'Cashmeran', 'Vainilla', 'Tonka', 'Ládano', 'Patchouli'],
    price: 2950000,
    stock: 10,
    image: 'https://www.pariscornerperfumes.com/cdn/shop/products/s-l1600_8555abe4-8c06-4e28-9072-d21df9f08a1e.jpg?v=1644237892&width=900'
  },
  {
    slug: 'paris-corner-taskeen-noche',
    name: 'Paris Corner Taskeen Noche',
    tagline: 'Fruta dulce, cognac suave y fondo cremoso',
    description: 'Perfil nocturno de Taskeen con durazno y naranja sanguina, un corazón floral con davana y cognac, y una base dulce de sándalo, benjuí, vainilla, tonka, ládano y patchouli.',
    family: 'Gourmand',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Durazno', 'Naranja sanguina', 'Cardamomo'],
    heart: ['Heliotropo', 'Davana', 'Cognac', 'Jazmín'],
    base: ['Sándalo', 'Benjuí', 'Cashmeran', 'Vainilla', 'Tonka', 'Ládano', 'Patchouli'],
    price: 2950000,
    stock: 10,
    image: 'https://www.pariscornerperfumes.com/cdn/shop/products/s-l1600_8555abe4-8c06-4e28-9072-d21df9f08a1e.jpg?v=1644237892&width=900'
  },
  {
    slug: 'paris-corner-khair-confection',
    name: 'Paris Corner Khair Confection',
    tagline: 'Pera, crema batida y vainilla malvavisco',
    description: 'Gourmand dulce y frutal con pera y crema batida, corazón de jazmín, ylang-ylang y cashmeran, y fondo de sándalo, malvavisco y vainilla.',
    family: 'Gourmand',
    gender: 'Unisex',
    concentration: 'EDP',
    top: ['Pera', 'Crema batida'],
    heart: ['Jazmín', 'Ylang-ylang', 'Cashmeran'],
    base: ['Sándalo', 'Malvavisco', 'Vainilla'],
    price: 2950000,
    stock: 10,
    image: 'https://www.pariscornerperfumes.com/cdn/shop/files/KHAIRCONFECTION01.jpg?v=1726730340&width=900'
  }
].map(makeProduct);

function makeProduct(seed: DemoProductSeed, displayOrder: number): Product {
  const id = seed.slug;
  return {
    id,
    slug: seed.slug,
    name: seed.name,
    tagline: seed.tagline,
    description: seed.description,
    scent_family: seed.family,
    gender_tag: seed.gender,
    concentration: seed.concentration,
    top_notes: seed.top,
    heart_notes: seed.heart,
    base_notes: seed.base,
    featured: true,
    active: true,
    display_order: displayOrder,
    variants: [
      { id: `${seed.slug}-100`, product_id: id, size_ml: 100, price_ars_cents: seed.price, stock: seed.stock, sku: `${seed.slug.toUpperCase()}-100`, active: true, weight_grams: 320 }
    ],
    images: [
      { id: `${seed.slug}-1`, product_id: id, url: seed.image, alt_text: seed.name, is_primary: true, sort_order: 0 }
    ],
    min_price_ars_cents: seed.price,
    total_stock: seed.stock
  };
}
