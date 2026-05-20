<script lang="ts">
  import { page } from '$app/stores';
  import { cartCount } from '$lib/stores/cart';
  import { wishlist } from '$lib/stores/wishlist';
  import { onMount } from 'svelte';
  export let onCart: () => void = () => undefined;
  let scrolled = false;
  $: activePath = $page.url.pathname;
  function isActive(href: string): boolean {
    return href === '/' ? activePath === '/' : activePath.startsWith(href);
  }
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
    <a href="/fragrances" class:active={isActive('/fragrances')}>Catálogo</a>
    <a href="/envios" class:active={isActive('/envios')}>Envíos</a>
    <a href="/contacto" class:active={isActive('/contacto')}>Contacto</a>
    {#if $wishlist.length > 0}<a href="/lista-de-deseos" class:active={isActive('/lista-de-deseos')}>Mi lista de deseos ({$wishlist.length})</a>{/if}
  </nav>
  <button class="cart-btn" type="button" on:click={onCart} aria-label="Abrir carrito">Bolsa {$cartCount}</button>
</header>

<style>
  .site-header { position: sticky; top: 0; z-index: 20; height: 72px; padding: 0 max(20px, 5vw); display: flex; align-items: center; justify-content: space-between; background: color-mix(in srgb, var(--color-bg) 0%, transparent); border-bottom: 1px solid transparent; transition: background .3s ease, border-color .3s ease, backdrop-filter .3s ease, box-shadow .3s ease; }
  .site-header.scrolled { background: color-mix(in srgb, var(--color-bg) 88%, transparent); border-bottom-color: var(--color-border); backdrop-filter: blur(16px); box-shadow: 0 1px 0 var(--color-border), 0 8px 24px color-mix(in srgb, var(--color-bg) 40%, transparent); }
  .brand { color: var(--color-text); font-size: 1.55rem; }
  nav { display: flex; gap: 28px; color: var(--color-text-muted); font-size: .92rem; }
  nav a { position: relative; }
  nav a::after { content: ''; position: absolute; left: 0; bottom: -4px; width: 0; height: 1px; background: var(--color-gold); transition: width 0.25s ease; }
  nav a.active { color: var(--color-gold); }
  nav a.active::after, nav a:hover::after { width: 100%; }
  nav a:hover { color: var(--color-gold); }
  .cart-btn { border: 1px solid var(--color-border); background: transparent; color: var(--color-text); padding: 10px 14px; }
  @media (max-width: 720px) { nav { gap: 12px; font-size: .8rem; } .brand { font-size: 1.25rem; } }
</style>
