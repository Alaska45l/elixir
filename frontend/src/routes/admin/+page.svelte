<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import MetricsCard from '$lib/components/admin/MetricsCard.svelte';
  import OrderTable from '$lib/components/admin/OrderTable.svelte';
  import type { AdminOrder } from '$lib/types/admin';
  import { formatARS } from '$lib/utils/currency';
  type Metrics = { total_products: number; total_orders: number; paid_revenue_cents: number; pending_orders: number; low_stock_count: number; recent_orders: AdminOrder[] };
  let metrics: Metrics = { total_products: 0, total_orders: 0, paid_revenue_cents: 0, pending_orders: 0, low_stock_count: 0, recent_orders: [] };
  let lowStock: { id: string; name: string; size_ml: number; stock: number }[] = [];
  let error = '';
  onMount(async () => {
    try {
      metrics = await apiFetch<Metrics>('/api/admin/metrics');
      lowStock = (await apiFetch<{ items: typeof lowStock }>('/api/admin/low-stock')).items;
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo cargar el panel';
    }
  });
</script>
<h1 class="display section-title">Panel administrativo</h1>
{#if error}<p class="error">{error}</p>{/if}
<div class="metrics">
  <MetricsCard label="Productos" value={String(metrics.total_products)} />
  <MetricsCard label="Ingresos pagos" value={formatARS(metrics.paid_revenue_cents)} />
  <MetricsCard label="Pendientes" value={String(metrics.pending_orders)} />
  <MetricsCard label="Órdenes" value={String(metrics.total_orders)} />
  <MetricsCard label="Stock bajo" value={String(metrics.low_stock_count)} />
</div>
<h2>Órdenes recientes</h2><OrderTable orders={metrics.recent_orders} />
<h2>Stock limitado</h2>
{#if lowStock.length === 0}
  <p class="empty">No hay variantes por debajo del aviso de stock.</p>
{:else}
  <ul>{#each lowStock as item}<li>{item.name} {item.size_ml}ml · {item.stock}</li>{/each}</ul>
{/if}
<style>.metrics{display:grid;grid-template-columns:repeat(5,1fr);gap:16px;margin-bottom:36px}h2{margin-top:38px}.empty{color:var(--color-text-muted)}.error{color:var(--color-danger-soft)}@media(max-width:1100px){.metrics{grid-template-columns:repeat(2,1fr)}}@media(max-width:700px){.metrics{grid-template-columns:1fr}}</style>
