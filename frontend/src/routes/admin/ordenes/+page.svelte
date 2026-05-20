<script lang="ts">
  import { env } from '$env/dynamic/public';
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import type { AdminOrder } from '$lib/types/admin';
  import { formatARS } from '$lib/utils/currency';

  const statuses = ['pending', 'paid', 'failed', 'cancelled', 'shipped', 'delivered'];
  let orders: AdminOrder[] = [];
  let expanded = '';

  async function load() {
    orders = (await apiFetch<{ items: AdminOrder[] }>('/api/admin/orders')).items;
  }
  async function save(order: AdminOrder) {
    await apiFetch(`/api/admin/orders/${order.id}`, {
      method: 'PUT',
      body: JSON.stringify({ status: order.status, tracking_number: order.tracking_number ?? '' })
    });
    await load();
  }
  function exportCSV() {
    window.open(`${env.PUBLIC_API_URL ?? 'http://localhost:8080'}/api/admin/orders/export`, '_blank');
  }
  onMount(load);
</script>

<div class="head">
  <h1 class="display section-title">Órdenes</h1>
  <button class="btn primary" type="button" on:click={exportCSV}>Exportar CSV</button>
</div>

<div class="orders">
  {#each orders as order}
    <article>
      <button class="summary" type="button" on:click={() => expanded = expanded === order.id ? '' : order.id}>
        <span>{order.external_reference}</span>
        <span>{order.customer_name}</span>
        <span>{order.status}</span>
        <strong>{formatARS(order.total_ars_cents)}</strong>
      </button>
      {#if expanded === order.id}
        <div class="detail">
          <div><b>{order.customer_name}</b><p>{order.customer_email} · {order.customer_phone}</p></div>
          <pre>{JSON.stringify(order.shipping_address ?? {}, null, 2)}</pre>
          <table class="table"><tbody>{#each order.items ?? [] as item}<tr><td>{item.product_name}</td><td>{item.size_ml}ml</td><td>{item.quantity}</td><td>{formatARS(item.subtotal_ars_cents)}</td></tr>{/each}</tbody></table>
          <div class="edit-row">
            <select class="select" bind:value={order.status}>{#each statuses as status}<option value={status}>{status}</option>{/each}</select>
            <input class="input" bind:value={order.tracking_number} placeholder="Tracking" />
            <button class="btn primary" type="button" on:click={() => save(order)}>Guardar</button>
          </div>
        </div>
      {/if}
    </article>
  {/each}
</div>

<style>
  .head { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
  .orders { display: grid; gap: 12px; }
  article { border: 1px solid var(--color-border); background: var(--color-surface); }
  .summary { width: 100%; display: grid; grid-template-columns: 1.4fr 1.3fr .8fr auto; gap: 16px; align-items: center; padding: 16px; border: 0; color: var(--color-text); background: transparent; text-align: left; }
  .detail { border-top: 1px solid var(--color-border); padding: 18px; display: grid; gap: 16px; }
  pre { margin: 0; color: var(--color-text-muted); white-space: pre-wrap; }
  .edit-row { display: grid; grid-template-columns: 180px 1fr auto; gap: 12px; }
  @media (max-width: 800px) { .summary, .edit-row { grid-template-columns: 1fr; } }
</style>
