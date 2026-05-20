<script lang="ts">
  import ProductCard from './ProductCard.svelte';
  import type { Product } from '$lib/api/client';
  export let products: Product[] = [];
</script>

<div class="grid-3">
  {#each products as product (product.id)}
    <ProductCard {product} />
  {/each}
</div>
