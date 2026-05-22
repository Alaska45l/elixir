<script lang="ts">
  import WishlistButton from './WishlistButton.svelte';
  import type { Product } from '$lib/api/client';
  import { formatARS } from '$lib/utils/currency';
  import { reveal } from '$lib/utils/reveal';
  export let product: Product;
  export let index = 0;
  const primary = product.images[0]?.url ?? '';
  const secondary = product.images[1]?.url ?? primary;
  const limited = product.total_stock > 0 && product.total_stock <= 5;
</script>

<article class={`product reveal-delay-${(index % 3) + 1}`} use:reveal>
  <a class="media" href={`/fragrances/${product.slug}`}>
    <img class="primary" src={primary} alt={product.images[0]?.alt_text ?? product.name} loading="lazy" />
    <img class="secondary" src={secondary} alt="" loading="lazy" />
    {#if product.total_stock === 0}<span class="badge empty">Agotado</span>{/if}
    {#if limited}<span class="badge">Stock limitado</span>{/if}
  </a>
  <div class="info">
    <div>
      <h3 class="display">{product.name}</h3>
      <p>{product.tagline}</p>
    </div>
    <WishlistButton slug={product.slug} />
  </div>
  <div class="meta">
    <span>Desde {formatARS(product.min_price_ars_cents)}</span>
    <a href={`/fragrances/${product.slug}`}>Ver detalles</a>
  </div>
</article>

<style>
  .product {
    display: grid;
    gap: 0;
    cursor: pointer;
  }

  .media {
    position: relative;
    aspect-ratio: 3 / 4;
    overflow: hidden;
    background: color-mix(in srgb, var(--color-bg) 82%, var(--color-surface));
  }

  .media::after {
    content: '';
    position: absolute;
    inset: 0;
    background: color-mix(in srgb, var(--color-bg) 0%, transparent);
    transition: background 0.5s ease;
    pointer-events: none;
  }

  .product:hover .media::after {
    background: color-mix(in srgb, var(--color-bg) 6%, transparent);
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.7s cubic-bezier(0.25, 0.46, 0.45, 0.94), opacity 0.4s ease;
    transform-origin: center bottom;
  }

  .media:hover img.primary {
    transform: scale(1.06);
  }

  .secondary {
    position: absolute;
    inset: 0;
    opacity: 0;
    transition: opacity 0.4s ease;
  }

  .media:hover .primary {
    opacity: 0;
  }

  .media:hover .secondary {
    opacity: 1;
    transform: scale(1.06);
  }

  .badge {
    position: absolute;
    left: 0;
    top: 0;
    color: var(--color-gold);
    background: var(--color-bg);
    padding: 5px 10px;
    font-size: 0.65rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    z-index: 2;
    border-bottom: 1px solid var(--color-border);
    border-right: 1px solid var(--color-border);
  }

  .badge.empty {
    color: var(--color-danger-soft);
  }

  .info {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 12px;
    padding: 14px 0 0;
    border-top: 1px solid var(--color-border);
    margin-top: 12px;
  }

  h3 {
    margin: 0 0 4px;
    font-size: clamp(1.1rem, 2.5vw, 1.45rem);
    line-height: 1;
    letter-spacing: 0;
    transition: color 0.2s;
  }

  .product:hover h3 {
    color: var(--color-gold);
  }

  p {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 0.78rem;
    line-height: 1.5;
    max-width: 28ch;
  }

  .meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0 0;
    gap: 12px;
  }

  .meta span {
    color: var(--color-text-muted);
    font-size: 0.82rem;
  }

  .meta a {
    color: var(--color-gold);
    font-size: 0.78rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    white-space: nowrap;
  }
</style>
