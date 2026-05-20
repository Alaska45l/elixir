<script lang="ts">
  import { formatARS } from '$lib/utils/currency';
  import type { AdminOrder } from '$lib/types/admin';
  export let orders: AdminOrder[] = [];
</script>
<table class="table">
  <thead><tr><th>Orden</th><th>Cliente</th><th>Estado</th><th>Total</th></tr></thead>
  <tbody>
    {#each orders as order}
      <tr><td>{order.external_reference}</td><td>{order.customer_name}<br /><span>{order.customer_email}</span></td><td>{order.status}</td><td>{formatARS(order.total_ars_cents)}</td></tr>
    {/each}
  </tbody>
</table>
<style>span{color:var(--color-text-muted);font-size:.85rem}</style>
