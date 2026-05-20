<script lang="ts">
  import { onMount } from 'svelte';
  import { getProduct } from '$lib/api/client';
  import type { Product } from '$lib/api/client';
  import ProductCard from './ProductCard.svelte';

  export let currentSlug: string;
  let products: Product[] = [];

  onMount(async () => {
    const slugs = (JSON.parse(localStorage.getItem('elixir_recent') ?? '[]') as string[])
      .filter((slug) => slug !== currentSlug)
      .slice(0, 6);
    const loaded = await Promise.all(slugs.map((slug) => getProduct(slug)));
    products = loaded.filter((item): item is Product => item !== null);
  });
</script>

{#if products.length >= 2}
  <section class="container page-pad">
    <p class="eyebrow">Vistas recientemente</p>
    <div class="gold-rule"></div>
    <div class="recent-strip">
      {#each products as product, i}
        <div class="recent-card"><ProductCard {product} index={i} /></div>
      {/each}
    </div>
  </section>
{/if}

<style>
  .recent-strip { display: grid; grid-auto-flow: column; grid-auto-columns: minmax(240px, 320px); gap: 22px; overflow-x: auto; padding-bottom: 10px; }
</style>
