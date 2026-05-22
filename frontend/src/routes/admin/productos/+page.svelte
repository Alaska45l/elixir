<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  type Row = { id: string; name: string; slug: string; featured: boolean; active: boolean };
  let products: Row[] = [];
  onMount(async () => products = (await apiFetch<{ items: Row[] }>('/api/admin/products')).items);
  async function toggleActive(product: Row) {
    product.active = !product.active;
    await apiFetch(`/api/admin/products/${product.id}/active`, { method: 'PUT', body: JSON.stringify({ active: product.active }) });
    products = [...products];
  }
</script>
<div class="head"><h1 class="display section-title">Productos</h1><a class="btn primary" href="/admin/productos/nuevo">Nuevo producto</a></div>
<table class="table"><thead><tr><th>Nombre</th><th>Slug</th><th>Destacado</th><th>Activo</th><th></th></tr></thead><tbody>{#each products as p}<tr><td>{p.name}</td><td>{p.slug}</td><td>{p.featured ? 'Sí' : 'No'}</td><td><button class="switch" class:on={p.active} type="button" on:click={() => toggleActive(p)} aria-pressed={p.active} aria-label={p.active ? `Desactivar ${p.name}` : `Activar ${p.name}`}><span></span></button></td><td><a href={`/admin/productos/${p.id}`}>Editar</a></td></tr>{/each}</tbody></table>
<style>.head{display:flex;justify-content:space-between;align-items:center;gap:20px}.switch{width:48px;height:26px;border:1px solid var(--color-border);background:rgba(232,224,208,.06);padding:2px}.switch span{display:block;width:20px;height:20px;background:var(--color-text-muted);transition:transform .2s ease,background .2s ease}.switch.on span{transform:translateX(20px);background:var(--color-gold)}</style>
