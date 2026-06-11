<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { formatARS } from '$lib/utils/currency';

  type Row = {
    id: string;
    name: string;
    slug: string;
    featured: boolean;
    active: boolean;
    min_price_ars_cents: number;
    total_stock: number;
    variant_count: number;
  };
  let products: Row[] = [];
  let loading = true;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      products = (await apiFetch<{ items: Row[] }>('/api/admin/products')).items;
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudieron cargar productos';
    } finally {
      loading = false;
    }
  }

  async function toggleActive(product: Row) {
    const next = !product.active;
    const label = next ? 'activar' : 'desactivar';
    if (!next && !confirm(`¿Desactivar ${product.name}? No se verá en la tienda.`)) return;
    product.active = next;
    products = [...products];
    try {
      await apiFetch(`/api/admin/products/${product.id}/active`, { method: 'PUT', body: JSON.stringify({ active: next }) });
      toast.push(next ? 'Producto activado' : 'Producto desactivado');
    } catch (err) {
      product.active = !next;
      products = [...products];
      error = err instanceof Error ? err.message : `No se pudo ${label} el producto`;
    }
  }

  onMount(load);
</script>

<div class="head">
  <h1 class="display section-title">Productos</h1>
  <a class="btn primary" href="/admin/productos/nuevo">Nuevo producto</a>
</div>

{#if error}<p class="error">{error}</p>{/if}
{#if loading}
  <p class="empty">Cargando productos...</p>
{:else if products.length === 0}
  <section class="empty-box">
    <h2>No hay productos cargados</h2>
    <p>Creá el primer producto con sus variantes, stock e imágenes para publicarlo en el catálogo.</p>
    <a class="btn primary" href="/admin/productos/nuevo">Crear producto</a>
  </section>
{:else}
  <table class="table">
    <thead><tr><th>Producto</th><th>Precio</th><th>Stock</th><th>Variantes</th><th>Estado</th><th></th></tr></thead>
    <tbody>
      {#each products as p}
        <tr>
          <td><strong>{p.name}</strong><br /><span>{p.slug}{p.featured ? ' · Destacado' : ''}</span></td>
          <td>{p.min_price_ars_cents ? formatARS(p.min_price_ars_cents) : 'Sin precio'}</td>
          <td>{p.total_stock}</td>
          <td>{p.variant_count}</td>
          <td><button class="switch" class:on={p.active} type="button" on:click={() => toggleActive(p)} aria-pressed={p.active} aria-label={p.active ? `Desactivar ${p.name}` : `Activar ${p.name}`}><span></span>{p.active ? 'Activo' : 'Inactivo'}</button></td>
          <td><a class="btn" href={`/admin/productos/${p.id}`}>Editar</a></td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .head { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
  .error { color: var(--color-danger-soft); }
  .empty { color: var(--color-text-muted); }
  .empty-box { display: grid; gap: 12px; border-top: 1px solid var(--color-border); padding-top: 22px; max-width: 560px; }
  .empty-box h2 { margin: 0; }
  .empty-box p, td span { color: var(--color-text-muted); }
  .switch { min-width: 104px; height: 34px; display: inline-flex; align-items: center; gap: 8px; border: 1px solid var(--color-border); background: rgba(45,42,36,.04); color: var(--color-text-muted); padding: 2px 10px 2px 2px; }
  .switch span { display: block; width: 28px; height: 28px; background: var(--color-text-muted); transition: background .2s ease; }
  .switch.on { color: var(--color-text); }
  .switch.on span { background: var(--color-emerald); }
</style>
