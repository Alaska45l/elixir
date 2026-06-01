<script lang="ts">
  import { onMount } from 'svelte';
  import { getProduct } from '$lib/api/client';
  import type { Product } from '$lib/api/client';
  import ProductCard from './ProductCard.svelte';

  export let currentSlug: string;
  let products: Product[] = [];

  onMount(async () => {
    let slugs: string[] = [];
    try {
      const parsed = JSON.parse(localStorage.getItem('elixir_recent') ?? '[]') as unknown;
      slugs = Array.isArray(parsed)
        ? parsed.filter((item): item is string => typeof item === 'string' && item !== currentSlug).slice(0, 6)
        : [];
    } catch {
      slugs = [];
    }
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
