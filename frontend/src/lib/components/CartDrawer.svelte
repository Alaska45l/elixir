<script lang="ts">
  import { page } from '$app/stores';
  import QuantityStepper from '$lib/components/QuantityStepper.svelte';
  import { cart, cartSubtotal } from '$lib/stores/cart';
  import { formatARS } from '$lib/utils/currency';
  export let open = false;
  export let onClose: () => void = () => undefined;
  let lastPath = '';
  $: if (lastPath && $page.url.pathname !== lastPath && open) onClose();
  $: lastPath = $page.url.pathname;
</script>

{#if open}
  <button class="shade" type="button" on:click={onClose} aria-label="Cerrar carrito"></button>
  <aside class="drawer">
    <button class="close" type="button" on:click={onClose}>Cerrar</button>
    <h2 class="display">Carrito</h2>
    <div class="items">
      {#if $cart.length === 0}
        <div class="empty">
          <p class="eyebrow">Tu carrito está vacío</p>
          <a class="btn primary" href="/fragrances">Ver fragancias</a>
        </div>
      {:else}
        {#each $cart as item}
        <article>
          <img src={item.image} alt={item.productName} />
          <div>
            <strong>{item.productName}</strong>
            <span>{item.sizeML}ml · {formatARS(item.unitPriceCents)}</span>
            <div class="qty">
              <QuantityStepper value={item.quantity} min={1} max={99} onchange={(value) => cart.setQuantity(item.variantId, value)} />
              <button type="button" on:click={() => cart.remove(item.variantId)}>Quitar</button>
            </div>
            <span class="item-total">{formatARS(item.quantity * item.unitPriceCents)}</span>
          </div>
        </article>
        {/each}
      {/if}
    </div>
    <div class="summary">
      <span>Total parcial</span><strong>{formatARS($cartSubtotal)}</strong>
    </div>
    <a class="btn primary" href="/carrito">Finalizar compra</a>
  </aside>
{/if}

<style>
  .shade { position: fixed; inset: 0; background: var(--color-overlay); z-index: 40; border: 0; padding: 0; }
  .drawer { position: fixed; right: 0; top: 0; bottom: 0; width: min(420px, 100vw); background: var(--color-bg); border-left: 1px solid var(--color-border); z-index: 50; padding: 28px; display: grid; grid-template-rows: auto auto 1fr auto auto; gap: 20px; }
  .close { justify-self: end; background: transparent; border: 0; color: var(--color-text-muted); }
  h2 { margin: 0; font-size: 2.6rem; }
  .items { overflow: auto; display: grid; align-content: start; gap: 16px; }
  article { display: grid; grid-template-columns: 78px 1fr; gap: 14px; }
  img { width: 78px; height: 96px; object-fit: cover; }
  article span { display: block; color: var(--color-text-muted); margin-top: 4px; }
  .qty { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; margin-top: 10px; }
  .qty button { background: rgba(232,224,208,0.04); color: var(--color-gold); border: 0; border-radius: 6px; min-width: 72px; height: 46px; }
  .item-total { display: block; margin-top: 8px; color: var(--color-gold); font-size: .9rem; }
  .summary { border-top: 1px solid var(--color-border); padding-top: 16px; display: flex; justify-content: space-between; }
  .empty { display: grid; gap: 16px; align-content: center; min-height: 220px; }

  @media (prefers-reduced-motion: no-preference) {
    .shade { animation: fadeIn .22s ease both; }
    .drawer { animation: slideIn .32s cubic-bezier(.22, 1, .36, 1) both; }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    @keyframes slideIn {
      from { transform: translateX(100%); opacity: 0; }
      to { transform: none; opacity: 1; }
    }
  }
</style>
