<script lang="ts">
  import { demoProducts } from '$lib/api/client';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
</script>

<section class="container page-pad error-page">
  <p class="eyebrow">Catálogo</p>
  <h1 class="display section-title">No pudimos cargar esta fragancia.</h1>
  <div class="gold-rule"></div>
  <a class="btn primary" href="/fragrances">Ver catálogo</a>
  <ProductGrid products={demoProducts.slice(0, 3)} />
</section>

<style>
  .error-page { display: grid; gap: 22px; }
</style>
