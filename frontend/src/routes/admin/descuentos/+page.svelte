<script lang="ts">
  import { apiFetch } from '$lib/api/client';
  import { onMount } from 'svelte';
  let code = ''; let discount_type = 'percent'; let discount_value = 10; let rows: { id: string; code: string; discount_type: string; discount_value: number; active: boolean }[] = [];
  async function load(){ rows = (await apiFetch<{items: typeof rows}>('/api/admin/discounts')).items; }
  async function create(){ await apiFetch('/api/admin/discounts',{method:'POST',body:JSON.stringify({code,discount_type,discount_value,min_order_cents:0,max_uses:null})}); code=''; await load(); }
  onMount(load);
</script>
<h1 class="display section-title">Descuentos</h1>
<form class="inline" on:submit|preventDefault={create}><input class="input" required placeholder="Código" bind:value={code}/><select class="select" bind:value={discount_type}><option value="percent">Percent</option><option value="fixed">Fixed</option></select><input class="input" type="number" bind:value={discount_value}/><button class="btn primary">Crear</button></form>
<table class="table"><tbody>{#each rows as row}<tr><td>{row.code}</td><td>{row.discount_type}</td><td>{row.discount_value}</td><td>{row.active ? 'Activo' : 'Inactivo'}</td></tr>{/each}</tbody></table>
<style>.inline{display:grid;grid-template-columns:1fr 160px 120px auto;gap:12px;margin-bottom:24px}</style>
