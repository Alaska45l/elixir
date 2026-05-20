<script lang="ts">
  import { env } from '$env/dynamic/public';
  import { apiFetch } from '$lib/api/client';
  import type { CartValidation, DiscountValidation, Order } from '$lib/api/client';
  import { cart, cartSubtotal } from '$lib/stores/cart';
  import { toast } from '$lib/stores/toast';
  import { formatARS } from '$lib/utils/currency';
  import type { PageData } from './$types';
  export let data: PageData;
  let name = '';
  let email = '';
  let phone = '';
  let address = '';
  let province = 'CF';
  let postalCode = '';
  let discountCode = '';
  let discount = 0;
  $: shipping = data.zones.find((z) => z.province_codes.includes(province))?.base_cost_cents ?? data.zones.at(-1)?.base_cost_cents ?? 0;
  $: total = $cartSubtotal + shipping - discount;
  async function applyDiscount() {
    const res = await apiFetch<DiscountValidation>('/api/discount/validate', { method: 'POST', body: JSON.stringify({ code: discountCode, subtotal_ars_cents: $cartSubtotal }) });
    discount = res.valid ? res.discount_ars_cents : 0;
    toast.push(res.valid ? 'Código aplicado' : res.message ?? 'Código inválido', res.valid ? 'ok' : 'error');
  }
  async function abandoned() {
    if (email && $cart.length) await apiFetch('/api/contact/abandoned-cart', { method: 'POST', body: JSON.stringify({ email, cart_data: { items: $cart } }) });
  }
  async function checkout() {
    const items = $cart.map((item) => ({ variant_id: item.variantId, quantity: item.quantity, unit_price_ars_cents: item.unitPriceCents }));
    const validation = await apiFetch<CartValidation>('/api/cart/validate', { method: 'POST', body: JSON.stringify({ items }) });
    if (!validation.valid) { toast.push(validation.errors.join('. '), 'error'); return; }
    const order = await apiFetch<Order>('/api/orders', { method: 'POST', body: JSON.stringify({ items, customer_name: name, customer_email: email, customer_phone: phone, shipping_address: { address, province, postalCode }, shipping_cost_ars_cents: shipping, discount_code: discountCode }) });
    const pref = await apiFetch<{ init_point: string }>('/api/checkout/mercadopago/preference', { method: 'POST', body: JSON.stringify({ external_reference: order.external_reference }) });
    location.href = pref.init_point;
  }
  function wa() {
    const summary = $cart.map((i) => `${i.productName} ${i.sizeML}ml x${i.quantity}`).join(', ');
    return `https://wa.me/${env.PUBLIC_WHATSAPP_NUMBER ?? '5491100000000'}?text=${encodeURIComponent(`Hola, quiero consultar por mi carrito: ${summary}`)}`;
  }
  function updateQuantity(event: Event, variantId: string) {
    const input = event.currentTarget as HTMLInputElement;
    cart.setQuantity(variantId, Number(input.value));
  }
</script>

<svelte:head><title>Carrito | ELIXIR Exclusive</title></svelte:head>

<section class="container page-pad">
  <p class="eyebrow">Finalizar compra</p>
  <h1 class="display section-title">Carrito</h1>
  <div class="gold-rule"></div>
  {#if $cart.length === 0}
    <div class="empty">
      <p class="eyebrow">Tu carrito está vacío</p>
      <h2 class="display">Explorá la colección.</h2>
      <a class="btn primary" href="/fragrances">Ver fragancias</a>
    </div>
  {:else}
  <div class="cart-page">
    <div class="lines">
      {#each $cart as item}
        <article>
          <img src={item.image} alt={item.productName} />
          <div><h2>{item.productName}</h2><p>{item.sizeML}ml · {formatARS(item.unitPriceCents)}</p></div>
          <input class="input" type="number" min="1" value={item.quantity} on:input={(event) => updateQuantity(event, item.variantId)} />
          <strong>{formatARS(item.quantity * item.unitPriceCents)}</strong>
          <button type="button" on:click={() => cart.remove(item.variantId)}>Quitar</button>
        </article>
      {/each}
    </div>
    <form class="summary" on:submit|preventDefault={checkout}>
      <label class="field"><span>Nombre</span><input class="input" required bind:value={name} /></label>
      <label class="field"><span>Email</span><input class="input" type="email" required bind:value={email} on:blur={abandoned} /></label>
      <label class="field"><span>Teléfono</span><input class="input" required bind:value={phone} /></label>
      <label class="field"><span>Dirección</span><input class="input" required bind:value={address} /></label>
      <div class="grid-2 compact"><label class="field"><span>Provincia</span><select class="select" bind:value={province}><option value="CF">CABA</option><option value="BA">Buenos Aires</option><option value="AR">Interior</option></select></label><label class="field"><span>Código postal</span><input class="input" required bind:value={postalCode} /></label></div>
      <div class="discount"><input class="input" placeholder="Código de descuento" bind:value={discountCode} /><button class="btn" type="button" on:click={applyDiscount}>Aplicar</button></div>
      <div class="totals"><span>Subtotal</span><b>{formatARS($cartSubtotal)}</b><span>Envío</span><b>{formatARS(shipping)}</b><span>Descuento</span><b>{formatARS(discount)}</b><span>Total</span><strong>{formatARS(total)}</strong></div>
      <button class="btn primary" type="submit" disabled={$cart.length === 0}>Finalizar compra</button>
      <a class="btn" href={wa()} target="_blank" rel="noreferrer">Consultar por WhatsApp</a>
    </form>
  </div>
  {/if}
</section>

<style>
  .cart-page { display: grid; grid-template-columns: 1fr 390px; gap: 52px; align-items: start; }
  .lines { display: grid; gap: 18px; }
  article { display: grid; grid-template-columns: 92px 1fr 74px auto auto; gap: 16px; align-items: center; border-top: 1px solid var(--color-border); padding-top: 18px; }
  article img { width: 92px; height: 110px; object-fit: cover; }
  h2 { margin: 0; font-size: 1rem; }
  p { margin: 6px 0 0; color: var(--color-text-muted); }
  article button { background: transparent; border: 0; color: var(--color-gold); }
  .summary { display: grid; gap: 16px; border: 1px solid var(--color-border); padding: 22px; background: var(--color-surface); }
  .compact { gap: 12px; }
  .discount { display: grid; grid-template-columns: 1fr auto; gap: 10px; }
  .totals { display: grid; grid-template-columns: 1fr auto; gap: 10px; color: var(--color-text-muted); border-top: 1px solid var(--color-border); padding-top: 16px; }
  .totals strong { color: var(--color-gold); font-size: 1.35rem; }
  .empty { min-height: 50vh; display: flex; flex-direction: column; justify-content: center; gap: 18px; }
  .empty h2 { font-size: clamp(2.8rem, 6vw, 5rem); margin: 0; }
  @media (max-width: 920px) { .cart-page, article { grid-template-columns: 1fr; } }
</style>
