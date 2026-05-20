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
  button { min-width: 98px; height: 44px; border: 1px solid var(--color-border); background: transparent; color: var(--color-text); }
  button.active { border-color: var(--color-gold); color: var(--color-gold); }
  button:disabled { opacity: .45; cursor: not-allowed; }
</style>
