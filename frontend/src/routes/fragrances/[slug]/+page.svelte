<script lang="ts">
  import { browser } from '$app/environment';
  import ImageGallery from '$lib/components/ImageGallery.svelte';
  import NoteCloud from '$lib/components/NoteCloud.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import RecentlyViewed from '$lib/components/RecentlyViewed.svelte';
  import SizeSelector from '$lib/components/SizeSelector.svelte';
  import WishlistButton from '$lib/components/WishlistButton.svelte';
  import { cart } from '$lib/stores/cart';
  import { toast } from '$lib/stores/toast';
  import { formatARS } from '$lib/utils/currency';
  import type { PageData } from './$types';
  export let data: PageData;
  let selected = data.product.variants.find((v) => v.stock > 0) ?? data.product.variants[0];
  let qty = 1;
  if (browser) {
    const key = 'elixir_recent';
    const existing = JSON.parse(localStorage.getItem(key) ?? '[]') as string[];
    localStorage.setItem(key, JSON.stringify([data.product.slug, ...existing.filter((x) => x !== data.product.slug)].slice(0, 6)));
  }
  function add() {
    if (!selected || selected.stock === 0) return;
    cart.add({
      variantId: selected.id,
      productSlug: data.product.slug,
      productName: data.product.name,
      image: data.product.images[0]?.url ?? '',
      sizeML: selected.size_ml,
      unitPriceCents: selected.price_ars_cents,
      quantity: qty
    });
    toast.push('Añadido al carrito');
  }
</script>

<svelte:head>
  <title>{data.product.name} | ELIXIR Exclusive</title>
  <meta name="description" content={data.product.tagline} />
  <meta property="og:title" content={`${data.product.name} | ELIXIR Exclusive`} />
  <meta property="og:description" content={data.product.description} />
  <meta property="og:image" content={data.product.images[0]?.url} />
</svelte:head>

<section class="container page-pad pdp">
  <ImageGallery images={data.product.images} />
  <div class="panel">
    <div class="topline"><span>{data.product.scent_family} · {data.product.concentration}</span><WishlistButton slug={data.product.slug} /></div>
    <h1 class="display">{data.product.name}</h1>
    <div class="gold-rule"></div>
    <p class="tagline">{data.product.tagline}</p>
    <p class="price">{selected ? formatARS(selected.price_ars_cents) : 'Agotado'}</p>
    <SizeSelector variants={data.product.variants} {selected} onSelect={(variant) => selected = variant} />
    <NoteCloud top={data.product.top_notes} heart={data.product.heart_notes} base={data.product.base_notes} />
    <div class="buy">
      <div class="qty-stepper">
        <button type="button" on:click={() => qty = Math.max(1, qty - 1)}>−</button>
        <span>{qty}</span>
        <button type="button" on:click={() => qty = Math.min(selected?.stock ?? 1, qty + 1)}>+</button>
      </div>
      <button class="btn primary" type="button" disabled={!selected || selected.stock === 0} on:click={add}>Añadir al carrito</button>
    </div>
  </div>
</section>

<section class="container page-pad">
  <p class="eyebrow">También puede gustarte</p>
  <div class="gold-rule"></div>
  <ProductGrid products={data.related} />
</section>

<RecentlyViewed currentSlug={data.product.slug} />

<style>
  .pdp { display: grid; grid-template-columns: 1.05fr .95fr; gap: 64px; align-items: start; }
  .panel { display: grid; gap: 26px; position: sticky; top: 98px; }
  .topline { display: flex; justify-content: space-between; color: var(--color-gold); text-transform: uppercase; letter-spacing: .12em; font-size: .78rem; }
  h1 { font-size: clamp(3.2rem, 7vw, 6.8rem); line-height: .86; margin: 0; }
  .tagline { color: var(--color-text-muted); font-size: 1.15rem; line-height: 1.7; max-width: 50ch; }
  .price { font-size: 1.5rem; color: var(--color-gold); }
  .buy { display: grid; grid-template-columns: 92px 1fr; gap: 12px; }
  .qty-stepper { display: flex; border: 1px solid var(--color-border); height: 46px; }
  .qty-stepper button { width: 44px; background: transparent; border: 0; color: var(--color-text-muted); font-size: 1.3rem; }
  .qty-stepper button:hover { color: var(--color-gold); }
  .qty-stepper span { flex: 1; display: grid; place-items: center; border-left: 1px solid var(--color-border); border-right: 1px solid var(--color-border); }
  @media (max-width: 900px) { .pdp { grid-template-columns: 1fr; } .panel { position: static; } }
</style>
