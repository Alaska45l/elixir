<script lang="ts">
  import ProductCard from './ProductCard.svelte';
  import ProductCardSkeleton from './ProductCardSkeleton.svelte';
  import type { Product } from '$lib/api/client';

  type CardVariant = 'standard' | 'featured' | 'limited';

  export let products: Product[] = [];
  export let loading = false;

  function isLimited(product: Product): boolean {
    return product.total_stock > 0 && product.total_stock <= 5;
  }

  function variantFor(product: Product): CardVariant {
    if (isLimited(product)) return 'limited';
    if (product.featured) return 'featured';
    return 'standard';
  }
</script>

<div class="mosaic">
  {#if loading}
    {#each Array(8) as _}
      <div class="mosaic-item"><ProductCardSkeleton /></div>
    {/each}
  {:else}
    {#each products as product, i (product.id)}
      {@const limited = isLimited(product)}
      <div class="mosaic-item" class:featured={product.featured} class:limited>
        <ProductCard {product} index={i} variant={variantFor(product)} />
      </div>
    {/each}
  {/if}
</div>

<style>
  .mosaic {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: auto;
    gap: 14px;
    grid-auto-flow: dense;
  }

  .mosaic-item {
    min-width: 0;
  }

  .mosaic-item :global(.product) {
    height: 100%;
  }

  .mosaic-item :global(.card-link) {
    min-height: 300px;
  }

  .mosaic-item:nth-child(1) {
    grid-column: span 2;
    grid-row: span 2;
  }

  .mosaic-item:nth-child(1) :global(.card-link) {
    min-height: 560px;
  }

  .mosaic-item:nth-child(2) :global(.card-link),
  .mosaic-item:nth-child(3) :global(.card-link) {
    min-height: 320px;
  }

  .mosaic-item:nth-child(4) {
    grid-column: span 2;
  }

  .mosaic-item:nth-child(4) :global(.card-link) {
    min-height: 320px;
  }

  .mosaic-item.featured:not(:nth-child(1)):not(:nth-child(4)) {
    grid-column: span 2;
  }

  .mosaic-item.featured:not(:nth-child(1)):not(:nth-child(4)) :global(.card-link) {
    min-height: 440px;
  }

  .mosaic-item.limited:not(:nth-child(1)) :global(.card-link) {
    min-height: 400px;
  }

  @media (max-width: 860px) {
    .mosaic {
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;
    }

    .mosaic-item:nth-child(1) {
      grid-column: span 2;
      grid-row: span 1;
    }

    .mosaic-item:nth-child(1) :global(.card-link) {
      min-height: 380px;
    }

    .mosaic-item:nth-child(4) {
      grid-column: span 2;
    }

    .mosaic-item.featured,
    .mosaic-item.limited:not(:nth-child(1)) {
      grid-column: span 2;
    }
  }

  @media (max-width: 480px) {
    .mosaic {
      grid-template-columns: 1fr;
      gap: 16px;
    }

    .mosaic-item,
    .mosaic-item:nth-child(1),
    .mosaic-item:nth-child(4),
    .mosaic-item.featured,
    .mosaic-item.limited {
      grid-column: span 1;
      grid-row: span 1;
    }
  }
</style>
