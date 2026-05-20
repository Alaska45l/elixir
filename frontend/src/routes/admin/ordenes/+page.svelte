<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import OrderTable from '$lib/components/admin/OrderTable.svelte';
  import type { AdminOrder } from '$lib/types/admin';
  let orders: AdminOrder[] = [];
  onMount(async () => orders = (await apiFetch<{ items: AdminOrder[] }>('/api/admin/orders')).items);
</script>
<h1 class="display section-title">Órdenes</h1><OrderTable {orders} />
