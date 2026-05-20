<script lang="ts">
  import type { ProductImage } from '$lib/api/client';
  export let images: ProductImage[] = [];
  let active = 0;
</script>

<div class="gallery">
  <img class="main" src={images[active]?.url} alt={images[active]?.alt_text ?? 'Fragancia'} />
  <div class="thumbs">
    {#each images as image, i}
      <button type="button" class:active={i === active} on:click={() => active = i}><img src={image.url} alt={image.alt_text ?? ''} /></button>
    {/each}
  </div>
</div>

<style>
  .main { width: 100%; aspect-ratio: 4 / 5; object-fit: cover; background: var(--color-surface); }
  .thumbs { display: flex; gap: 10px; margin-top: 12px; }
  button { width: 72px; height: 72px; border: 1px solid var(--color-border); background: transparent; padding: 0; }
  button.active { border-color: var(--color-gold); }
  button img { width: 100%; height: 100%; object-fit: cover; }
</style>
