<script lang="ts">
  import FilterSidebar from '$lib/components/FilterSidebar.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import SearchOverlay from '$lib/components/SearchOverlay.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
  let searchOpen = false;
  $: params = new URLSearchParams(data.search);
</script>

<svelte:head>
  <title>Catálogo | ELIXIR Exclusive</title>
  <meta name="description" content="Catálogo de fragancias ELIXIR Exclusive: orientales, florales, amaderadas, cítricas, frescas y gourmand." />
</svelte:head>

<section class="container page-pad catalog-head">
  <p class="eyebrow">Catálogo</p>
  <h1 class="display section-title">Fragancias para presencia precisa</h1>
  <button class="btn" type="button" on:click={() => searchOpen = true}>Buscar fragancias</button>
</section>

<SearchOverlay open={searchOpen} />

<section class="container catalog">
  <FilterSidebar {params} />
  <div>
    <ProductGrid products={data.products} />
    <nav class="pagination">
      <a class="btn" href="/fragrances?limit=24&offset=0">1</a>
      <a class="btn" href="/fragrances?limit=24&offset=24">2</a>
    </nav>
  </div>
</section>

<style>
  .catalog-head { padding-bottom: 28px; }
  .catalog-head .btn { margin-top: -12px; }
  .catalog { display: grid; grid-template-columns: 240px 1fr; gap: 44px; align-items: start; }
  .pagination { display: flex; gap: 10px; margin-top: 44px; }
  @media (max-width: 860px) { .catalog { grid-template-columns: 1fr; } }
</style>
