<script lang="ts">
  import ProductCard from './ProductCard.svelte';
  import ProductCardSkeleton from './ProductCardSkeleton.svelte';
  import type { Product } from '$lib/api/client';

  export let products: Product[] = [];
  export let loading = false;
</script>

<div class="product-grid">
  {#if loading}
    {#each Array(6) as _}<ProductCardSkeleton />{/each}
  {:else}
    {#each products as product, i (product.id)}
      <ProductCard {product} index={i} />
    {/each}
  {/if}
</div>

<style>
  .product-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 28px;
  }

  @media (max-width: 860px) {
    .product-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 16px;
    }
  }

  @media (max-width: 480px) {
    .product-grid {
      grid-template-columns: 1fr;
      gap: 20px;
    }
  }
</style>
