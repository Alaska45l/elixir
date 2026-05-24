<script lang="ts">
  import WishlistButton from './WishlistButton.svelte';
  import type { Product } from '$lib/api/client';
  import { formatARS } from '$lib/utils/currency';
  import { reveal } from '$lib/utils/reveal';
  type ProductCardVariant = 'standard' | 'featured' | 'limited';
  export let product: Product;
  export let index = 0;
  export let variant: ProductCardVariant = 'standard';
  $: primary = product.images[0]?.url ?? '';
  $: secondary = product.images[1]?.url ?? primary;
  $: limited = product.total_stock > 0 && product.total_stock <= 5;
</script>

<article class={`product variant-${variant} reveal-delay-${(index % 3) + 1}`} use:reveal>
  <a class="card-link" href={`/fragrances/${product.slug}`}>
    <img class="primary" src={primary} alt={product.images[0]?.alt_text ?? product.name} loading="lazy" />
    <img class="secondary" src={secondary} alt="" loading="lazy" />

    <div class="card-overlay">
      <div class="card-badges">
        {#if product.total_stock === 0}<span class="badge empty">Agotado</span>{/if}
        {#if limited}<span class="badge">Stock limitado</span>{/if}
      </div>

      <div class="card-info">
        <h3 class="display">{product.name}</h3>
        {#if product.tagline}<p class="tagline">{product.tagline}</p>{/if}
        <span class="price">Desde {formatARS(product.min_price_ars_cents)}</span>
        {#if variant === 'featured'}
          <span class="card-cta btn primary">Añadir al carrito</span>
        {/if}
        <span class="card-detail">Ver detalles</span>
      </div>
    </div>
  </a>
  <WishlistButton slug={product.slug} />
</article>

<style>
  .product {
    position: relative;
    border-radius: 14px;
    overflow: hidden;
    cursor: pointer;
  }

  .product.variant-limited {
    box-shadow: 0 0 20px rgba(184,151,94,0.15);
  }

  .card-link {
    display: block;
    position: relative;
    width: 100%;
    height: 100%;
    min-height: 320px;
    overflow: hidden;
  }

  .card-link img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.7s cubic-bezier(0.25, 0.46, 0.45, 0.94), opacity 0.4s ease;
  }

  .secondary {
    opacity: 0;
  }

  .product:hover img.primary {
    opacity: 0;
    transform: scale(1.06);
  }

  .product:hover img.secondary {
    opacity: 1;
    transform: scale(1.06);
  }

  .card-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    z-index: 2;
    background: linear-gradient(
      to top,
      rgba(13,31,21,0.88) 0%,
      rgba(13,31,21,0.4) 40%,
      rgba(13,31,21,0.05) 70%,
      transparent 100%
    );
    padding: 18px;
  }

  .card-badges {
    display: flex;
    gap: 8px;
    align-self: flex-start;
  }

  .badge {
    padding: 4px 10px;
    font-size: 0.65rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--color-gold);
    background: rgba(13,31,21,0.7);
    border-radius: 4px;
  }

  .badge.empty {
    color: var(--color-danger-soft);
  }

  .card-info {
    display: grid;
    gap: 6px;
    align-self: flex-end;
    width: fit-content;
    max-width: min(100%, 34ch);
    padding: 15px;
    border-radius: 6px;
    background: #0D1F158C;
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    box-shadow: 0 12px 32px rgba(0,0,0,0.18);
  }

  h3 {
    margin: 0;
    font-size: clamp(1.2rem, 2.5vw, 1.7rem);
    line-height: 1.05;
    color: var(--color-text);
    letter-spacing: 0;
    transition: color 0.2s ease;
  }

  .product:hover h3 {
    color: var(--color-gold);
  }

  .tagline {
    margin: 0;
    color: color-mix(in srgb, var(--color-text) 76%, var(--color-text-muted));
    font-size: 0.78rem;
    line-height: 1.4;
    max-width: 28ch;
  }

  .price {
    color: color-mix(in srgb, var(--color-text) 70%, var(--color-text-muted));
    font-size: 0.82rem;
  }

  .card-detail {
    color: var(--color-gold);
    font-size: 0.72rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    margin-top: 4px;
  }

  .card-cta {
    margin-top: 8px;
    width: fit-content;
    min-height: 38px;
    padding: 0 16px;
    font-size: 0.8rem;
  }

  .product :global(.wishlist-btn) {
    position: absolute;
    top: 14px;
    right: 14px;
    z-index: 3;
  }

  .card-cta::before {
    display: none;
  }

  .card-link:hover {
    color: var(--color-text);
  }

  @media (max-width: 860px) {
    .product :global(.wishlist-btn) {
      width: 48px;
      height: 48px;
    }
  }
</style>
