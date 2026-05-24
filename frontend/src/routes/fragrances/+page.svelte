<script lang="ts">
  import BentoFilterBar from '$lib/components/BentoFilterBar.svelte';
  import BentoProductMosaic from '$lib/components/BentoProductMosaic.svelte';
  import FilterBottomSheet from '$lib/components/FilterBottomSheet.svelte';
  import SearchOverlay from '$lib/components/SearchOverlay.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
  let searchOpen = false;
  let filtersOpen = false;
  let innerWidth = 1024;

  $: params = new URLSearchParams(data.search);
  $: activeFilterCount = getActiveFilterCount(params);

  function pageHref(offset: number) {
    const next = new URLSearchParams(data.search);
    next.set('limit', String(data.limit));
    next.set('offset', String(offset));
    return `/fragrances?${next.toString()}`;
  }

  function getActiveFilterCount(current: URLSearchParams): number {
    let count = current.getAll('family').length + current.getAll('gender').length + current.getAll('concentration').length;
    if (current.get('in_stock') === 'true') count += 1;
    if (Number(current.get('min_price') ?? 0) > 0 || Number(current.get('max_price') ?? 20000000) < 20000000) count += 1;
    return count;
  }
</script>

<svelte:window bind:innerWidth />

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

<section class="container catalog-filters">
  {#if innerWidth > 860}
    <BentoFilterBar {params} />
  {:else}
    <button class="filter-summary" type="button" on:click={() => filtersOpen = true}>
      <span>Filtros</span>
      <strong>({activeFilterCount} {activeFilterCount === 1 ? 'activo' : 'activos'})</strong>
    </button>
  {/if}
</section>

<FilterBottomSheet open={filtersOpen && innerWidth <= 860} {params} onClose={() => filtersOpen = false} />

<section class="container catalog">
  <div>
    <BentoProductMosaic products={data.products} />
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
  .catalog-filters { margin-bottom: 24px; }
  .filter-summary {
    width: 100%;
    min-height: 58px;
    border: 0;
    border-radius: 6px;
    background: var(--color-surface);
    color: var(--color-text);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 0 18px;
    box-shadow: 0 4px 24px rgba(0,0,0,0.15);
  }
  .filter-summary span { color: var(--color-text); font-weight: 700; }
  .filter-summary strong { color: var(--color-gold); font-size: .86rem; font-weight: 700; }
  .catalog { display: grid; gap: 44px; align-items: start; }
  .pagination { display: flex; gap: 10px; margin-top: 44px; flex-wrap: wrap; }
  .pagination .active { border-color: var(--color-gold); color: var(--color-gold); }
  @media (max-width: 860px) {
    .catalog-filters { position: sticky; top: 84px; z-index: 16; margin-bottom: 18px; }
    .pagination .btn { min-width: 48px; }
  }
</style>
