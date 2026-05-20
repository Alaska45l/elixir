<script lang="ts">
  import { onMount } from 'svelte';
  let visible = false;
  onMount(() => {
    const onScroll = () => visible = window.scrollY > 400;
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  });
  function top() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
</script>

<button class="scroll-top btn" class:visible type="button" on:click={top} aria-label="Volver arriba">↑</button>

<style>
  .scroll-top { position: fixed; right: 18px; bottom: 82px; z-index: 35; width: 46px; padding: 0; opacity: 0; pointer-events: none; transform: translateY(10px); transition: opacity .2s ease, transform .2s ease; }
  .scroll-top.visible { opacity: 1; pointer-events: auto; transform: none; }
</style>
