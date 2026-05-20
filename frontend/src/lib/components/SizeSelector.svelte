<script lang="ts">
  import type { Variant } from '$lib/api/client';
  export let variants: Variant[] = [];
  export let selected: Variant | undefined;
  export let onSelect: (variant: Variant) => void = () => undefined;
</script>

<div class="sizes">
  {#each variants as variant}
    <button type="button" disabled={variant.stock === 0} class:active={selected?.id === variant.id} on:click={() => onSelect(variant)}>
      {variant.size_ml}ml {variant.stock === 0 ? 'Agotado' : ''}
    </button>
  {/each}
</div>

<style>
  .sizes { display: flex; flex-wrap: wrap; gap: 10px; }
  button { min-width: 98px; height: 44px; border: 1px solid var(--color-border); background: transparent; color: var(--color-text-muted); font-size: .85rem; letter-spacing: .04em; transition: border-color .2s ease, background .2s ease, color .2s ease; }
  button:hover:not(:disabled):not(.active) { border-color: color-mix(in srgb, var(--color-gold) 50%, transparent); color: var(--color-text); }
  button.active { background: var(--color-gold); border-color: var(--color-gold); color: var(--color-bg); font-weight: 600; }
  button:disabled { opacity: .35; cursor: not-allowed; }
</style>
