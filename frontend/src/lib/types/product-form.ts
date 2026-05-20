export type VariantFormValue = { size_ml: number; price_ars_cents: number; stock: number; sku: string };
export type ImageFormValue = { url: string; alt_text: string; is_primary: boolean; sort_order: number };

export type ProductFormValue = {
  slug: string;
  name: string;
  tagline: string;
  description: string;
  scent_family: string;
  gender_tag: string;
  concentration: string;
  top_notes: string[];
  heart_notes: string[];
  base_notes: string[];
  featured: boolean;
  display_order: number;
  variants: VariantFormValue[];
  images: ImageFormValue[];
};
