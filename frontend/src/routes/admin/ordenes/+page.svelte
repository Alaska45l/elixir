<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { AdminOrder } from '$lib/types/admin';
  import { formatARS } from '$lib/utils/currency';

  const statuses = [
    ['pending', 'Pendiente'],
    ['paid', 'Pagada'],
    ['failed', 'Fallida'],
    ['cancelled', 'Cancelada'],
    ['shipped', 'Despachada'],
    ['delivered', 'Entregada']
  ];
  let orders: AdminOrder[] = [];
  let expanded = '';
  let loading = true;
  let error = '';
  let saving = '';

  async function load() {
    loading = true;
    error = '';
    try {
      orders = (await apiFetch<{ items: AdminOrder[] }>('/api/admin/orders')).items;
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudieron cargar órdenes';
    } finally {
      loading = false;
    }
  }
  async function save(order: AdminOrder) {
    saving = order.id;
    error = '';
    try {
      await apiFetch(`/api/admin/orders/${order.id}`, {
        method: 'PUT',
        body: JSON.stringify({
          status: order.status,
          tracking_number: order.tracking_number ?? '',
          shipping_carrier: order.shipping_carrier ?? '',
          shipped_at: order.shipped_at || null,
          internal_notes: order.internal_notes ?? ''
        })
      });
      toast.push('Orden guardada');
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo guardar la orden';
    } finally {
      saving = '';
    }
  }
  function exportCSV() {
    window.open('/api/admin/orders/export', '_blank');
  }
  function label(status: string) {
    return statuses.find(([value]) => value === status)?.[1] ?? status;
  }
  function address(order: AdminOrder) {
    const data = order.shipping_address ?? {};
    return [
      data.address,
      data.province,
      data.postalCode,
      data.shipping_option
    ].filter(Boolean).join(' · ');
  }
  function dateInput(value?: string) {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    return new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
  }
  function setShippedAt(order: AdminOrder, value: string) {
    order.shipped_at = value ? new Date(value).toISOString() : undefined;
  }
  onMount(load);
</script>

<div class="head">
  <h1 class="display section-title">Órdenes</h1>
  <button class="btn primary" type="button" on:click={exportCSV}>Exportar CSV</button>
</div>

{#if error}<p class="error">{error}</p>{/if}
{#if loading}
  <p class="empty">Cargando órdenes...</p>
{:else if orders.length === 0}
  <section class="empty-box"><h2>No hay órdenes todavía</h2><p>Cuando alguien compre, la orden aparecerá acá para revisar pago, envío y estado.</p></section>
{:else}
  <div class="orders">
    {#each orders as order}
      <article>
        <button class="summary" type="button" on:click={() => expanded = expanded === order.id ? '' : order.id}>
          <span><b>{order.external_reference}</b><small>{new Date(order.created_at).toLocaleString('es-AR')}</small></span>
          <span>{order.customer_name}<small>{order.customer_email}</small></span>
          <span class={`badge ${order.status}`}>{label(order.status)}</span>
          <strong>{formatARS(order.total_ars_cents)}</strong>
        </button>
        {#if expanded === order.id}
          <div class="detail">
            <div class="info-grid">
              <section>
                <h2>Cliente</h2>
                <p><b>{order.customer_name}</b><br />{order.customer_email}{order.customer_phone ? ` · ${order.customer_phone}` : ''}</p>
              </section>
              <section>
                <h2>Pago</h2>
                {#if order.payment?.mp_status}
                  <p><b>MercadoPago: {order.payment.mp_status}</b><br />{order.payment.mp_status_detail || 'Sin detalle adicional'}</p>
                {:else}
                  <p>Sin pago registrado todavía.</p>
                {/if}
              </section>
              <section>
                <h2>Envío</h2>
                <p>{address(order) || 'Sin dirección cargada'}</p>
              </section>
            </div>

            <table class="table">
              <thead><tr><th>Producto</th><th>Tamaño</th><th>Cantidad</th><th>Total</th></tr></thead>
              <tbody>{#each order.items ?? [] as item}<tr><td>{item.product_name}</td><td>{item.size_ml}ml</td><td>{item.quantity}</td><td>{formatARS(item.subtotal_ars_cents)}</td></tr>{/each}</tbody>
            </table>

            <div class="totals">
              <span>Subtotal <b>{formatARS(order.subtotal_ars_cents)}</b></span>
              <span>Envío <b>{formatARS(order.shipping_cost_ars_cents)}</b></span>
              <span>Descuento <b>{formatARS(order.discount_ars_cents)}</b></span>
              <span>Total <b>{formatARS(order.total_ars_cents)}</b></span>
            </div>

            <div class="edit-row">
              <label class="field"><span>Estado</span><select class="select" bind:value={order.status}>{#each statuses as status}<option value={status[0]}>{status[1]}</option>{/each}</select></label>
              <label class="field"><span>Correo / transporte</span><input class="input" bind:value={order.shipping_carrier} placeholder="Correo Argentino, Andreani..." /></label>
              <label class="field"><span>Número de seguimiento</span><input class="input" bind:value={order.tracking_number} placeholder="Tracking" /></label>
              <label class="field"><span>Fecha de despacho</span><input class="input" type="datetime-local" value={dateInput(order.shipped_at)} on:input={(event) => setShippedAt(order, (event.currentTarget as HTMLInputElement).value)} /></label>
            </div>
            <label class="field"><span>Notas internas</span><textarea class="textarea" bind:value={order.internal_notes} placeholder="Solo para administración"></textarea></label>
            <div class="actions">
              <button class="btn primary" type="button" disabled={saving === order.id} on:click={() => save(order)}>{saving === order.id ? 'Guardando...' : 'Guardar orden'}</button>
            </div>
          </div>
        {/if}
      </article>
    {/each}
  </div>
{/if}

<style>
  .head { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
  .orders { display: grid; gap: 12px; }
  article { border: 1px solid var(--color-border); background: var(--color-surface); }
  .summary { width: 100%; display: grid; grid-template-columns: 1.4fr 1.3fr .8fr auto; gap: 16px; align-items: center; padding: 16px; border: 0; color: var(--color-text); background: transparent; text-align: left; }
  .detail { border-top: 1px solid var(--color-border); padding: 18px; display: grid; gap: 16px; }
  .summary span { display: grid; gap: 4px; }
  small, p, .empty { color: var(--color-text-muted); }
  p { margin: 0; }
  .error { color: var(--color-danger-soft); }
  .empty-box { border-top: 1px solid var(--color-border); padding-top: 22px; }
  .info-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
  .info-grid section { border: 1px solid var(--color-border); padding: 14px; }
  h2 { margin: 0 0 8px; font-size: .78rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .badge { justify-self: start; border: 1px solid var(--color-border); padding: 5px 10px; color: var(--color-text-muted); }
  .badge.paid, .badge.shipped, .badge.delivered { color: var(--color-gold); border-color: color-mix(in srgb, var(--color-gold) 55%, transparent); }
  .badge.failed, .badge.cancelled { color: var(--color-danger-soft); border-color: #9f5d55; }
  .totals { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; color: var(--color-text-muted); }
  .totals span { border-top: 1px solid var(--color-border); padding-top: 10px; display: flex; justify-content: space-between; gap: 8px; }
  .edit-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
  .actions { display: flex; justify-content: flex-end; }
  @media (max-width: 980px) { .summary, .edit-row, .info-grid, .totals { grid-template-columns: 1fr; } }
</style>
