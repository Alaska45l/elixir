<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { onDestroy } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import type { Product } from '$lib/api/client';
  export let open = false;
  const dispatch = createEventDispatcher<{ close: void }>();
  let q = '';
  let debounceTimer: ReturnType<typeof setTimeout>;
  let results: Pick<Product, 'id' | 'slug' | 'name' | 'min_price_ars_cents'>[] = [];
  async function search() {
    if (q.length < 2) { results = []; return; }
    const data = await apiFetch<{ items: Pick<Product, 'id' | 'slug' | 'name' | 'min_price_ars_cents'>[] }>(`/api/products/search?q=${encodeURIComponent(q)}`);
    results = data.items;
  }
  function debouncedSearch() {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(search, 300);
  }
  function close() {
    dispatch('close');
  }
  function keydown(event: KeyboardEvent) {
    if (event.key === 'Escape' && open) close();
  }
  onDestroy(() => clearTimeout(debounceTimer));
</script>

<svelte:window on:keydown={keydown} />

{#if open}
  <div class="overlay">
    <input class="input" bind:value={q} on:input={debouncedSearch} placeholder="Buscar fragancias" />
    <div class="results">
      {#each results as item}
        <a href={`/fragrances/${item.slug}`} on:click={close}>{item.name}</a>
      {/each}
    </div>
  </div>
{/if}

<style>
  .overlay { position: fixed; inset: 72px 0 auto; z-index: 30; background: var(--color-bg); border-bottom: 1px solid var(--color-border); padding: 24px max(20px, 8vw); }
  .results { display: grid; gap: 12px; margin-top: 16px; color: var(--color-gold); }
</style>
