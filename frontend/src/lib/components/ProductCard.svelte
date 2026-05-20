<script lang="ts">
  import WishlistButton from './WishlistButton.svelte';
  import type { Product } from '$lib/api/client';
  import { formatARS } from '$lib/utils/currency';
  export let product: Product;
  const primary = product.images[0]?.url ?? '';
  const secondary = product.images[1]?.url ?? primary;
  const limited = product.total_stock > 0 && product.total_stock <= 5;
</script>

<article class="product">
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
  .product { display: grid; gap: 16px; }
  .media { position: relative; aspect-ratio: 4 / 5; overflow: hidden; background: var(--color-surface); }
  img { width: 100%; height: 100%; object-fit: cover; transition: opacity .45s ease, transform 7s ease; }
  .secondary { position: absolute; inset: 0; opacity: 0; }
  .media:hover .primary { opacity: 0; }
  .media:hover .secondary { opacity: 1; transform: scale(1.04); }
  .info { display: flex; justify-content: space-between; gap: 18px; align-items: start; }
  h3 { margin: 0 0 6px; font-size: 1.7rem; }
  p { margin: 0; color: var(--color-text-muted); min-height: 42px; }
  .meta { display: flex; justify-content: space-between; color: var(--color-text-muted); }
  .meta a { color: var(--color-gold); }
  .badge { position: absolute; left: 12px; top: 12px; color: var(--color-gold); border: 1px solid var(--color-border); background: rgba(13,31,21,.8); padding: 7px 10px; font-size: .75rem; }
  .badge.empty { color: #e0a39a; }
</style>
