<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  type Row = { id: string; name: string; slug: string; featured: boolean; active: boolean };
  let products: Row[] = [];
  onMount(async () => products = (await apiFetch<{ items: Row[] }>('/api/admin/products')).items);
</script>
<div class="head"><h1 class="display section-title">Productos</h1><a class="btn primary" href="/admin/productos/nuevo">Nuevo producto</a></div>
<table class="table"><thead><tr><th>Nombre</th><th>Slug</th><th>Destacado</th><th>Activo</th><th></th></tr></thead><tbody>{#each products as p}<tr><td>{p.name}</td><td>{p.slug}</td><td>{p.featured ? 'Sí' : 'No'}</td><td>{p.active ? 'Sí' : 'No'}</td><td><a href={`/admin/productos/${p.id}`}>Editar</a></td></tr>{/each}</tbody></table>
<style>.head{display:flex;justify-content:space-between;align-items:center;gap:20px}</style>
