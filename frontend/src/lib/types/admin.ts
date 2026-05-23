export type AdminOrder = {
  id: string;
  external_reference: string;
  status: string;
  customer_name: string;
  customer_email: string;
  customer_phone?: string;
  shipping_address?: Record<string, string | number | boolean>;
  shipping_cost_ars_cents: number;
  subtotal_ars_cents: number;
  total_ars_cents: number;
  discount_code?: string;
  discount_ars_cents: number;
  currency: string;
  tracking_number?: string;
  shipping_carrier?: string;
  shipped_at?: string;
  internal_notes?: string;
  created_at: string;
  payment?: {
    mp_payment_id?: string;
    mp_preference_id?: string;
    mp_status?: string;
    mp_status_detail?: string;
    processed_at?: string;
  };
  items?: {
    product_name: string;
    size_ml: number;
    quantity: number;
    unit_price_ars_cents: number;
    subtotal_ars_cents: number;
  }[];
};

export type AdminUser = {
  username: string;
  created_at: string;
  last_login_at?: string;
};

export type AdminShippingZone = {
  id: string;
  zone_name: string;
  province_codes: string[];
  base_cost_cents: number;
  per_kg_cents: number;
  estimated_days_min: number;
  estimated_days_max: number;
  active: boolean;
};
