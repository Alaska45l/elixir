<script lang="ts">
  import { env } from '$env/dynamic/public';
  import { apiFetch, quoteShipping } from '$lib/api/client';
  import type { CartValidation, DiscountValidation, Order, ShippingQuoteOption } from '$lib/api/client';
  import QuantityStepper from '$lib/components/QuantityStepper.svelte';
  import { cart, cartSubtotal } from '$lib/stores/cart';
  import { toast } from '$lib/stores/toast';
  import { formatARS } from '$lib/utils/currency';
  let name = '';
  let email = '';
  let phone = '';
  let address = '';
  let province = 'C';
  let postalCode = '';
  let discountCode = '';
  let discount = 0;
  let shippingOptions: ShippingQuoteOption[] = [];
  let selectedShippingID = '';
  let quoting = false;
  const provinces = [
    ['C', 'Ciudad Autónoma de Buenos Aires'],
    ['B', 'Buenos Aires'],
    ['K', 'Catamarca'],
    ['H', 'Chaco'],
    ['U', 'Chubut'],
    ['X', 'Córdoba'],
    ['W', 'Corrientes'],
    ['E', 'Entre Ríos'],
    ['P', 'Formosa'],
    ['Y', 'Jujuy'],
    ['L', 'La Pampa'],
    ['F', 'La Rioja'],
    ['M', 'Mendoza'],
    ['N', 'Misiones'],
    ['Q', 'Neuquén'],
    ['R', 'Río Negro'],
    ['A', 'Salta'],
    ['J', 'San Juan'],
    ['D', 'San Luis'],
    ['Z', 'Santa Cruz'],
    ['S', 'Santa Fe'],
    ['G', 'Santiago del Estero'],
    ['V', 'Tierra del Fuego'],
    ['T', 'Tucumán']
  ];
  $: selectedShipping = shippingOptions.find((option) => option.id === selectedShippingID);
  $: shipping = selectedShipping?.price_cents ?? 0;
  $: cartWeightGrams = $cart.reduce((sum, item) => sum + (item.weightGrams ?? 200) * item.quantity, 0);
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
    if (!selectedShipping) { toast.push('Seleccioná una opción de envío', 'error'); return; }
    const items = $cart.map((item) => ({ variant_id: item.variantId, quantity: item.quantity, unit_price_ars_cents: item.unitPriceCents }));
    const validation = await apiFetch<CartValidation>('/api/cart/validate', { method: 'POST', body: JSON.stringify({ items }) });
    if (!validation.valid) { toast.push(validation.errors.join('. '), 'error'); return; }
    const order = await apiFetch<Order>('/api/orders', { method: 'POST', body: JSON.stringify({ items, customer_name: name, customer_email: email, customer_phone: phone, shipping_address: { address, province, postalCode, shipping_option: selectedShipping }, shipping_cost_ars_cents: shipping, discount_code: discountCode }) });
    const pref = await apiFetch<{ init_point: string }>('/api/checkout/mercadopago/preference', { method: 'POST', body: JSON.stringify({ external_reference: order.external_reference }) });
    location.href = pref.init_point;
  }
  async function loadShippingOptions() {
    if (!postalCode.trim()) { toast.push('Ingresá el código postal', 'error'); return; }
    quoting = true;
    try {
      shippingOptions = await quoteShipping({
        destination_postal_code: postalCode.trim(),
        province_code: province,
        weight_grams: cartWeightGrams || 200,
        dimensions: { length_cm: 22, width_cm: 14, height_cm: 10 }
      });
      selectedShippingID = shippingOptions[0]?.id ?? '';
      if (!shippingOptions.length) toast.push('No encontramos opciones de envío', 'error');
    } finally {
      quoting = false;
    }
  }
  function wa() {
    const summary = $cart.map((i) => `${i.productName} ${i.sizeML}ml x${i.quantity}`).join(', ');
    return `https://wa.me/${env.PUBLIC_WHATSAPP_NUMBER ?? '5491100000000'}?text=${encodeURIComponent(`Hola, quiero consultar por mi carrito: ${summary}`)}`;
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
        <article class="cart-line">
          <a class="line-image" href={`/fragrances/${item.productSlug}`}>
            <img src={item.image} alt={item.productName} />
          </a>
          <div class="line-details">
            <a class="line-name" href={`/fragrances/${item.productSlug}`}>
              <h2>{item.productName}</h2>
            </a>
            <p class="line-variant">{item.sizeML}ml</p>
            <p class="line-unit-price">{formatARS(item.unitPriceCents)} c/u</p>
          </div>
          <div class="line-quantity">
            <QuantityStepper
              value={item.quantity}
              min={1}
              max={99}
              onchange={(value) => cart.setQuantity(item.variantId, value)}
            />
          </div>
          <div class="line-total">
            <strong>{formatARS(item.quantity * item.unitPriceCents)}</strong>
          </div>
          <button class="line-remove" type="button" on:click={() => cart.remove(item.variantId)} aria-label="Quitar producto">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path d="M3 3L13 13M13 3L3 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
            </svg>
          </button>
        </article>
      {/each}
    </div>
    <form class="summary" on:submit|preventDefault={checkout}>
      <label class="field"><span>Nombre</span><input class="input" required bind:value={name} /></label>
      <label class="field"><span>Email</span><input class="input" type="email" required bind:value={email} on:blur={abandoned} /></label>
      <label class="field"><span>Teléfono</span><input class="input" required bind:value={phone} /></label>
      <label class="field"><span>Dirección</span><input class="input" required bind:value={address} /></label>
      <div class="grid-2 compact"><label class="field"><span>Provincia</span><select class="select" bind:value={province}>{#each provinces as item}<option value={item[0]}>{item[1]}</option>{/each}</select></label><label class="field"><span>Código postal</span><input class="input" required bind:value={postalCode} /></label></div>
      <section class="shipping-options">
        <div class="shipping-head"><span>Envío</span><button class="btn" type="button" on:click={loadShippingOptions} disabled={quoting}>{quoting ? 'Cotizando...' : 'Cotizar'}</button></div>
        {#if shippingOptions.length}
          <div class="option-list">
            {#each shippingOptions as option}
              <label class:active={selectedShippingID === option.id}>
                <input type="radio" name="shipping" bind:group={selectedShippingID} value={option.id} />
                <span><b>{option.service_name}</b><small>{option.carrier_name} · {option.estimated_days_min} a {option.estimated_days_max} días</small></span>
                <strong>{formatARS(option.price_cents)}</strong>
              </label>
            {/each}
          </div>
        {/if}
      </section>
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
  .lines { display: grid; gap: 0; }
  .cart-line {
    display: grid;
    grid-template-columns: 110px 1fr auto auto 40px;
    grid-template-rows: auto;
    gap: 20px;
    align-items: center;
    padding: 24px 0;
    border-bottom: 1px solid var(--color-border);
  }
  .cart-line:first-child { padding-top: 0; }
  .line-image {
    display: block;
    width: 110px;
    height: 130px;
    overflow: hidden;
    border-radius: 10px;
    background: var(--color-surface);
    flex-shrink: 0;
  }
  .line-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }
  .line-image:hover img { transform: scale(1.05); }
  .line-details {
    display: grid;
    gap: 4px;
    align-self: center;
    min-width: 0;
  }
  .line-name { color: var(--color-text); }
  .line-name h2 {
    margin: 0;
    font-size: 1.05rem;
    font-weight: 500;
    line-height: 1.3;
  }
  .line-name:hover h2 { color: var(--color-emerald-dark); }
  .line-variant {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 0.82rem;
  }
  .line-unit-price {
    margin: 0;
    color: var(--color-text-muted);
    font-size: 0.82rem;
  }
  .line-quantity {
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .line-total {
    text-align: right;
    min-width: 100px;
  }
  .line-total strong {
    color: var(--color-text);
    font-size: 1.05rem;
    font-weight: 600;
  }
  .line-remove {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border: 0;
    border-radius: 8px;
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: all 0.2s ease;
  }
  .line-remove:hover {
    color: var(--color-danger-soft);
    background: rgba(224,163,154,0.08);
  }
  .summary { display: grid; gap: 16px; border-radius: 14px; padding: 22px; background: var(--color-surface); box-shadow: 0 4px 24px rgba(0,0,0,0.15); }
  .compact { gap: 12px; }
  .discount { display: grid; grid-template-columns: 1fr auto; gap: 10px; }
  .shipping-options { display: grid; gap: 12px; }
  .shipping-head { display: flex; justify-content: space-between; align-items: center; color: var(--color-text-muted); }
  .option-list { display: grid; gap: 8px; }
  .option-list label { display: grid; grid-template-columns: auto 1fr auto; gap: 12px; align-items: center; border-radius: 6px; background: rgba(45,42,36,0.04); padding: 12px; color: var(--color-text-muted); cursor: pointer; }
  .option-list label.active { background: rgba(107,191,138,0.12); color: var(--color-text); }
  .option-list small { display: block; margin-top: 4px; color: var(--color-text-muted); }
  .totals { display: grid; grid-template-columns: 1fr auto; gap: 10px; color: var(--color-text-muted); border-top: 1px solid var(--color-border); padding-top: 16px; }
  .totals strong { color: var(--color-gold-dark); font-size: 1.35rem; }
  .empty { min-height: 50vh; display: flex; flex-direction: column; justify-content: center; gap: 18px; }
  .empty h2 { font-size: clamp(2.8rem, 6vw, 5rem); margin: 0; }
  @media (max-width: 920px) {
    .cart-page {
      grid-template-columns: 1fr;
      gap: 40px;
    }
    .cart-line {
      position: relative;
      grid-template-columns: 90px 1fr auto;
      grid-template-rows: auto auto;
      gap: 12px 16px;
    }
    .line-image {
      grid-row: 1 / 3;
      width: 90px;
      height: 108px;
    }
    .line-details {
      grid-column: 2 / 4;
      padding-right: 44px;
    }
    .line-quantity {
      grid-column: 2;
      justify-content: flex-start;
    }
    .line-total {
      text-align: right;
      min-width: auto;
    }
    .line-remove {
      position: absolute;
      top: 24px;
      right: 0;
    }
  }
</style>
