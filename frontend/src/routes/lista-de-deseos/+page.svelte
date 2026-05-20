<script lang="ts">
  import { onMount } from 'svelte';
  import { getProduct } from '$lib/api/client';
  import type { Product } from '$lib/api/client';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { wishlist } from '$lib/stores/wishlist';

  let products: Product[] = [];
  let loading = true;

  onMount(async () => {
    const loaded = await Promise.all($wishlist.map((slug) => getProduct(slug)));
    products = loaded.filter((item): item is Product => item !== null);
    loading = false;
  });
</script>

<svelte:head><title>Mi lista de deseos | ELIXIR Exclusive</title></svelte:head>

<section class="container page-pad">
  <p class="eyebrow">Mi lista de deseos</p>
  <h1 class="display section-title">Fragancias guardadas</h1>
  <div class="gold-rule"></div>
  {#if $wishlist.length === 0 && !loading}
    <div class="empty">
      <p class="muted">Todavía no guardaste fragancias.</p>
      <a class="btn primary" href="/fragrances">Ver catálogo</a>
    </div>
  {:else}
    <ProductGrid {products} {loading} />
  {/if}
</section>

<style>
  .empty { min-height: 40vh; display: grid; align-content: center; gap: 16px; }
</style>
