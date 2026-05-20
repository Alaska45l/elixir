export type AdminOrder = {
  id: string;
  external_reference: string;
  status: string;
  customer_name: string;
  customer_email: string;
  customer_phone?: string;
  shipping_address?: Record<string, string>;
  total_ars_cents: number;
  tracking_number?: string;
  items?: {
    product_name: string;
    size_ml: number;
    quantity: number;
    unit_price_ars_cents: number;
    subtotal_ars_cents: number;
  }[];
};
