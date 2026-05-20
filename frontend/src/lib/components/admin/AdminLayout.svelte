<script lang="ts">
  import { apiFetch } from '$lib/api/client';
  async function logout() {
    await apiFetch('/api/admin/logout', { method: 'POST', body: '{}' });
    location.href = '/admin/login';
  }
</script>

<div class="admin">
  <aside>
    <a class="display brand" href="/admin">ELIXIR</a>
    <a href="/admin/productos">Productos</a>
    <a href="/admin/ordenes">Órdenes</a>
    <a href="/admin/descuentos">Descuentos</a>
    <a href="/admin/homepage">Homepage</a>
    <a href="/admin/mensajes">Mensajes</a>
    <button type="button" on:click={logout}>Salir</button>
  </aside>
  <main><slot /></main>
</div>

<style>
  .admin { min-height: 100vh; display: grid; grid-template-columns: 240px 1fr; }
  aside { border-right: 1px solid var(--color-border); padding: 28px; display: grid; align-content: start; gap: 16px; background: var(--color-surface); }
  .brand { font-size: 1.7rem; margin-bottom: 20px; }
  a, button { color: var(--color-text-muted); background: transparent; border: 0; text-align: left; padding: 0; }
  a:hover, button:hover { color: var(--color-gold); }
  main { padding: 36px; overflow: auto; }
  @media (max-width: 800px) { .admin { grid-template-columns: 1fr; } aside { position: static; grid-template-columns: repeat(3, 1fr); } }
</style>
