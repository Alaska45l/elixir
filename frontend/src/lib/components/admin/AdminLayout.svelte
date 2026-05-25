<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  let ready = false;

  onMount(async () => {
    try {
      await apiFetch('/api/admin/me');
      ready = true;
    } catch {
      location.href = '/admin/login';
    }
  });

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
    <a href="/admin/envios">Envíos</a>
    <a href="/admin/descuentos">Descuentos</a>
    <a href="/admin/homepage">Homepage</a>
    <a href="/admin/configuracion">Configuración</a>
    <a href="/admin/mensajes">Mensajes</a>
    <a href="/admin/cuenta">Cuenta</a>
    <button type="button" on:click={logout}>Salir</button>
  </aside>
  <main>{#if ready}<slot />{:else}<p class="empty">Validando sesión...</p>{/if}</main>
</div>

<style>
  .admin { min-height: 100vh; display: grid; grid-template-columns: 240px 1fr; }
  aside { border-right: 1px solid var(--color-border); padding: 28px; display: grid; align-content: start; gap: 16px; background: var(--color-surface); }
  .brand { font-size: 1.7rem; margin-bottom: 20px; }
  a, button { color: var(--color-text-muted); background: transparent; border: 0; text-align: left; padding: 0; }
  a:hover, button:hover { color: var(--color-emerald-dark); }
  main { padding: 36px; overflow: auto; }
  .empty { color: var(--color-text-muted); }
  @media (max-width: 800px) { .admin { grid-template-columns: 1fr; } aside { position: static; grid-template-columns: repeat(auto-fit, minmax(110px, 1fr)); } }
</style>
