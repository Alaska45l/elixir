import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';

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
  hero_image_mode: 'static' | 'product_covers';
  hero_rotation_interval_ms: number;
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

export type ListResponse<T> = { items: T[]; total?: number; limit?: number; offset?: number; error?: string };

const DEFAULT_TIMEOUT_MS = 15000;

async function responseError(res: Response): Promise<Error> {
  const text = await res.text();
  let message = '';
  try {
    const body = JSON.parse(text) as { error?: string };
    message = body.error ?? '';
  } catch {
    message = text;
  }
  if (!message) {
    switch (res.status) {
      case 401:
        message = 'Sesión expirada. Iniciá sesión nuevamente.';
        break;
      case 403:
        message = 'No tenés permisos para realizar esta acción.';
        break;
      case 404:
        message = 'No encontramos el recurso solicitado.';
        break;
      default:
        message = res.status >= 500 ? 'Servicio temporalmente no disponible' : 'No se pudo completar la operación';
    }
  }
  return new Error(message);
}

function apiBase(path: string): string {
  const configuredBase = PUBLIC_API_URL.trim();
  const shouldUseSameOriginProxy = browser && path.startsWith('/api/');
  return shouldUseSameOriginProxy || path.startsWith('/api/admin') ? '' : configuredBase;
}

export async function apiFetch<T>(path: string, init?: RequestInit, fetcher: typeof fetch = fetch): Promise<T> {
  const controller = init?.signal ? null : new AbortController();
  const timeout = controller ? setTimeout(() => controller.abort(), DEFAULT_TIMEOUT_MS) : undefined;
  let res: Response;
  try {
    res = await fetcher(`${apiBase(path)}${path}`, {
      credentials: 'include',
      headers: { 'Content-Type': 'application/json', ...(init?.headers ?? {}) },
      ...init,
      signal: init?.signal ?? controller?.signal
    });
  } catch (err) {
    if (err instanceof DOMException && err.name === 'AbortError') {
      throw new Error('La solicitud tardó demasiado. Intentá nuevamente.');
    }
    throw new Error('No se pudo conectar con el servicio. Intentá nuevamente.');
  } finally {
    if (timeout) clearTimeout(timeout);
  }
  if (!res.ok) {
    const err = await responseError(res);
    if (browser && res.status === 401 && path.startsWith('/api/admin') && !path.includes('/login')) {
      window.location.assign('/admin/login');
    }
    throw err;
  }
  if (res.status === 204) {
    return undefined as T;
  }
  const text = await res.text();
  try {
    return (text ? JSON.parse(text) : undefined) as T;
  } catch {
    throw new Error('La respuesta del servicio no es válida.');
  }
}

export async function uploadAdminImage(file: File, folder = 'products', fetcher: typeof fetch = fetch): Promise<string> {
  const form = new FormData();
  form.append('file', file);
  form.append('folder', folder);

  const path = '/api/admin/upload';
  const res = await fetcher(`${apiBase(path)}${path}`, {
    method: 'POST',
    credentials: 'include',
    body: form
  });
  if (!res.ok) {
    throw await responseError(res);
  }
  const data = (await res.json()) as { url?: string };
  if (!data.url) {
    throw new Error('La subida no devolvió una URL');
  }
  return data.url;
}

export async function getProducts(fetcher: typeof fetch = fetch, query = ''): Promise<Product[]> {
  return (await getProductsResponse(fetcher, query)).items;
}

export async function getProductsResponse(fetcher: typeof fetch = fetch, query = ''): Promise<ListResponse<Product>> {
  try {
    const data = await apiFetch<ListResponse<Product>>(`/api/products${query}`, undefined, fetcher);
    return data;
  } catch {
    return { items: [], total: 0, error: 'No pudimos cargar el catálogo.' };
  }
}

export async function getProduct(slug: string, fetcher: typeof fetch = fetch): Promise<Product | null> {
  try {
    return await apiFetch<Product>(`/api/products/${slug}`, undefined, fetcher);
  } catch {
    return null;
  }
}

export async function getHomepage(fetcher: typeof fetch = fetch): Promise<HomepageSettings> {
  try {
    const data = await apiFetch<HomepageSettings>('/api/homepage', undefined, fetcher);
    return data.hero_heading ? { ...defaultHomepage, ...data } : defaultHomepage;
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
  hero_heading: 'Perfumería selecta con entrega nacional',
  hero_subheading: 'Fragancias originales, asesoramiento privado y despacho cuidado desde Buenos Aires.',
  hero_image_url: 'https://images.unsplash.com/photo-1619994403073-2cec844b8e63?auto=format&fit=crop&w=1200&q=85',
  hero_image_mode: 'product_covers',
  hero_rotation_interval_ms: 8000,
  hero_cta_label: 'Catálogo',
  hero_cta_url: '/fragrances',
  editorial_heading: 'Selección curada',
  editorial_body: 'Elegimos perfumes con trazabilidad, buena performance y una experiencia de compra sobria de punta a punta.',
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
    { label: 'Fragancias Unisex', href: '/fragrances?gender=Unisex' },
    { label: 'Fragancias Masculinas', href: '/fragrances?gender=Masculino' },
    { label: 'Fragancias Femeninas', href: '/fragrances?gender=Femenino' }
  ],
  low_stock_threshold: 5
};
