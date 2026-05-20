<script lang="ts">
  import FilterSidebar from '$lib/components/FilterSidebar.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import SearchOverlay from '$lib/components/SearchOverlay.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
  let searchOpen = false;
  $: params = new URLSearchParams(data.search);
  function pageHref(offset: number) {
    const next = new URLSearchParams(data.search);
    next.set('limit', String(data.limit));
    next.set('offset', String(offset));
    return `/fragrances?${next.toString()}`;
  }
</script>

<svelte:head>
  <title>Catálogo | ELIXIR Exclusive</title>
  <meta name="description" content="Catálogo de fragancias ELIXIR Exclusive: orientales, florales, amaderadas, cítricas, frescas y gourmand." />
</svelte:head>

<section class="container page-pad catalog-head">
  <p class="eyebrow">Catálogo</p>
  <h1 class="display section-title">Fragancias para presencia precisa</h1>
  <div class="gold-rule"></div>
  <button class="btn" type="button" on:click={() => searchOpen = true}>Buscar fragancias</button>
</section>

<SearchOverlay open={searchOpen} on:close={() => searchOpen = false} />

<section class="container catalog">
  <FilterSidebar {params} />
  <div>
    <ProductGrid products={data.products} />
    {#if data.total > data.limit}
      <nav class="pagination">
        {#each Array(Math.ceil(data.total / data.limit)) as _, i}
          {@const offset = i * data.limit}
          <a class="btn" class:active={offset === data.offset} href={pageHref(offset)}>{i + 1}</a>
        {/each}
      </nav>
    {/if}
  </div>
</section>

<style>
  .catalog-head { padding-bottom: 28px; }
  .catalog-head .btn { margin-top: -12px; }
  .catalog { display: grid; grid-template-columns: 240px 1fr; gap: 44px; align-items: start; }
  .pagination { display: flex; gap: 10px; margin-top: 44px; }
  .pagination .active { border-color: var(--color-gold); color: var(--color-gold); }
  @media (max-width: 860px) { .catalog { grid-template-columns: 1fr; } }
</style>
