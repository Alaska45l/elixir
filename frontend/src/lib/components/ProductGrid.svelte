<script lang="ts">
  import ProductCard from './ProductCard.svelte';
  import ProductCardSkeleton from './ProductCardSkeleton.svelte';
  import type { Product } from '$lib/api/client';
  export let products: Product[] = [];
  export let loading = false;
</script>

<div class="grid-3">
  {#if loading}
    {#each Array(6) as _}<ProductCardSkeleton />{/each}
  {:else}
    {#each products as product, i (product.id)}
      <ProductCard {product} index={i} />
    {/each}
  {/if}
</div>
