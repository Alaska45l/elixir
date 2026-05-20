<script lang="ts">
  import { cartCount } from '$lib/stores/cart';
  import { wishlist } from '$lib/stores/wishlist';
  import { onMount } from 'svelte';
  export let onCart: () => void = () => undefined;
  let scrolled = false;
  onMount(() => {
    const onScroll = () => scrolled = window.scrollY > 48;
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  });
</script>

<header class="site-header" class:scrolled>
  <a class="brand display" href="/">ELIXIR</a>
  <nav>
    <a href="/fragrances">Catálogo</a>
    <a href="/envios">Envíos</a>
    <a href="/contacto">Contacto</a>
    {#if $wishlist.length > 0}<a href="/lista-de-deseos">Mi lista de deseos ({$wishlist.length})</a>{/if}
  </nav>
  <button class="cart-btn" type="button" on:click={onCart} aria-label="Abrir carrito">Bolsa {$cartCount}</button>
</header>

<style>
  .site-header { position: sticky; top: 0; z-index: 20; height: 72px; padding: 0 max(20px, 5vw); display: flex; align-items: center; justify-content: space-between; background: rgba(13,31,21,0); border-bottom: 1px solid transparent; transition: background .3s ease, border-color .3s ease, backdrop-filter .3s ease; }
  .site-header.scrolled { background: rgba(13,31,21,.88); border-bottom-color: var(--color-border); backdrop-filter: blur(16px); }
  .brand { color: var(--color-text); font-size: 1.55rem; }
  nav { display: flex; gap: 28px; color: var(--color-text-muted); font-size: .92rem; }
  nav a:hover { color: var(--color-gold); }
  .cart-btn { border: 1px solid var(--color-border); background: transparent; color: var(--color-text); padding: 10px 14px; }
  @media (max-width: 720px) { nav { gap: 12px; font-size: .8rem; } .brand { font-size: 1.25rem; } }
</style>
