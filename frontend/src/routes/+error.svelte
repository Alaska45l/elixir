<script lang="ts">
  import { page } from '$app/stores';
</script>

<section class="container page-pad error-page">
  <p class="eyebrow">Error {$page.status}</p>
  <h1 class="display section-title">{$page.status === 404 ? 'Página no encontrada.' : 'Algo salió mal.'}</h1>
  <p class="muted">{$page.error?.message ?? 'Ocurrió un error inesperado.'}</p>
  <a class="btn" href="/">Volver al inicio</a>
</section>

<style>
  .error-page { min-height: 60vh; display: flex; flex-direction: column; justify-content: center; gap: 20px; }
</style>
