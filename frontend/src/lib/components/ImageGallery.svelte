<script lang="ts">
  import type { ProductImage } from '$lib/api/client';
  export let images: ProductImage[] = [];
  let active = 0;
</script>

<div class="gallery">
  <img class="main" src={images[active]?.url} alt={images[active]?.alt_text ?? 'Fragancia'} width="900" height="1125" />
  <div class="thumbs">
    {#each images as image, i}
      <button type="button" class:active={i === active} on:click={() => active = i}><img src={image.url} alt={image.alt_text ?? ''} loading="lazy" width="136" height="136" /></button>
    {/each}
  </div>
</div>

<style>
  .gallery { display: grid; gap: 12px; }
  .main { width: 100%; aspect-ratio: 4 / 5; object-fit: cover; object-position: center; border-radius: 14px; background: color-mix(in srgb, var(--color-bg) 82%, var(--color-surface)); }
  .thumbs { display: flex; gap: 8px; flex-wrap: wrap; }
  button { width: 68px; height: 68px; border: 1px solid transparent; border-radius: 6px; background: var(--color-surface); padding: 0; transition: border-color .2s ease, opacity .2s ease; overflow: hidden; position: relative; }
  button::after { content: ''; position: absolute; inset: 0; background: color-mix(in srgb, var(--color-emerald) 0%, transparent); transition: background .2s ease; }
  button:hover::after { background: color-mix(in srgb, var(--color-emerald) 8%, transparent); }
  button.active { border-color: var(--color-emerald); }
  button.active::after { background: color-mix(in srgb, var(--color-emerald) 6%, transparent); }
  button img { width: 100%; height: 100%; object-fit: cover; display: block; }
</style>
