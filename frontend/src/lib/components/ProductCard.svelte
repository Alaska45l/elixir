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
  .product { display: grid; gap: 16px; transition: transform .3s ease; }
  .product:hover { transform: translateY(-4px); }
  .media { position: relative; aspect-ratio: 3 / 4; overflow: hidden; background: var(--color-surface); }
  .media::after { content: ''; position: absolute; inset: 0; background: rgba(13,31,21,0); transition: background .45s ease; }
  .product:hover .media::after { background: rgba(13,31,21,.08); }
  img { width: 100%; height: 100%; object-fit: cover; transition: opacity .45s ease, transform .65s ease; }
  .media:hover img.primary { transform: scale(1.04); }
  .secondary { position: absolute; inset: 0; opacity: 0; transition: opacity .45s ease, transform .65s ease; }
  .media:hover .primary { opacity: 0; }
  .media:hover .secondary { opacity: 1; transform: scale(1.04); }
  .info { display: flex; justify-content: space-between; gap: 18px; align-items: start; }
  h3 { margin: 0 0 6px; font-size: clamp(1.4rem, 3vw, 1.9rem); line-height: .92; letter-spacing: -.01em; transition: color .2s; }
  .product:hover .info h3 { color: var(--color-gold); }
  p { margin: 0; color: var(--color-text-muted); min-height: 42px; }
  .meta { display: flex; justify-content: space-between; color: var(--color-text-muted); }
  .meta a { color: var(--color-gold); }
  .badge { position: absolute; left: 12px; top: 12px; color: var(--color-gold); border: 1px solid var(--color-border); background: rgba(13,31,21,.8); padding: 7px 10px; font-size: .75rem; }
  .badge.empty { color: #e0a39a; }
</style>
