<script lang="ts">
  import { env } from '$env/dynamic/public';
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import type { Order } from '$lib/api/client';
  import { formatARS } from '$lib/utils/currency';
  import type { PageData } from './$types';
  export let data: PageData;
  let order: Order | null = data.order;
  onMount(() => {
    if (order?.status !== 'pending') return;
    let count = 0;
    const timer = setInterval(async () => {
      count += 1;
      order = await apiFetch<Order>(`/api/orders/${data.ref}`);
      if (order.status !== 'pending' || count >= 24) clearInterval(timer);
    }, 5000);
    return () => clearInterval(timer);
  });
  $: status = order?.status ?? 'pending';
</script>

<svelte:head><title>Ver orden | ELIXIR Exclusive</title></svelte:head>

<section class="container page-pad confirmation">
  {#if status === 'paid'}
    <div class="icon">✓</div><h1 class="display">¡Gracias por tu compra!</h1><div class="gold-rule"></div><p>Tu orden {order?.external_reference} fue aprobada. Prepararemos el despacho estimado dentro de las próximas 24 a 48 horas hábiles.</p>
  {:else if status === 'failed'}
    <div class="icon error">×</div><h1 class="display">El pago no se completó.</h1><div class="gold-rule"></div><p>Podés volver al carrito e intentar nuevamente.</p><a class="btn primary" href="/carrito">Seguir comprando</a>
  {:else}
    <div class="icon">◷</div><h1 class="display">Tu pago está siendo procesado.</h1><div class="gold-rule"></div><p>Estamos verificando el estado de la orden {data.ref}. Esta página se actualizará automáticamente.</p>
  {/if}
  {#if order}
    <div class="summary"><span>Total</span><strong>{formatARS(order.total_ars_cents)}</strong><span>Estado</span><strong>{order.status}</strong></div>
  {/if}
  <a class="btn" href={`https://wa.me/${env.PUBLIC_WHATSAPP_NUMBER ?? '5491100000000'}?text=${encodeURIComponent(`Hola, consulto por la orden ${data.ref}`)}`} target="_blank" rel="noreferrer">Consultar por WhatsApp</a>
</section>

<style>
  .confirmation { max-width: 760px; text-align: center; display: grid; justify-items: center; gap: 20px; }
  .icon { width: 72px; height: 72px; border: 1px solid var(--color-gold); color: var(--color-gold); display: grid; place-items: center; font-size: 2rem; }
  .icon.error { border-color: #b66; color: #e0a39a; }
  h1 { font-size: clamp(3rem, 7vw, 6rem); line-height: .92; margin: 0; }
  p { color: var(--color-text-muted); line-height: 1.7; }
  .summary { width: min(420px, 100%); border-top: 1px solid var(--color-border); padding-top: 18px; display: grid; grid-template-columns: 1fr auto; gap: 10px; text-align: left; }
</style>
