<script lang="ts">
  import { cartCount } from '$lib/stores/cart';
  export let onCart: () => void = () => undefined;
</script>

<header class="site-header">
  <a class="brand display" href="/">ELIXIR</a>
  <nav>
    <a href="/fragrances">Catálogo</a>
    <a href="/envios">Envíos</a>
    <a href="/contacto">Contacto</a>
  </nav>
  <button class="cart-btn" type="button" on:click={onCart} aria-label="Abrir carrito">Bolsa {$cartCount}</button>
</header>

<style>
  .site-header { position: sticky; top: 0; z-index: 20; height: 72px; padding: 0 max(20px, 5vw); display: flex; align-items: center; justify-content: space-between; background: rgba(13,31,21,.88); border-bottom: 1px solid var(--color-border); backdrop-filter: blur(16px); }
  .brand { color: var(--color-text); font-size: 1.55rem; }
  nav { display: flex; gap: 28px; color: var(--color-text-muted); font-size: .92rem; }
  nav a:hover { color: var(--color-gold); }
  .cart-btn { border: 1px solid var(--color-border); background: transparent; color: var(--color-text); padding: 10px 14px; }
  @media (max-width: 720px) { nav { gap: 12px; font-size: .8rem; } .brand { font-size: 1.25rem; } }
</style>
