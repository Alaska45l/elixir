<script lang="ts">
  import { page } from '$app/stores';
  import MobileNav from '$lib/components/MobileNav.svelte';
  import { cartCount } from '$lib/stores/cart';
  import { wishlist } from '$lib/stores/wishlist';
  import { onMount } from 'svelte';

  export let onCart: () => void = () => undefined;

  let scrolled = false;
  let mobileNavOpen = false;

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
  <button class="hamburger" type="button" on:click={() => mobileNavOpen = true} aria-label="Abrir menú">
    <svg width="22" height="16" viewBox="0 0 22 16" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M1 1H21M1 8H21M1 15H21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
    </svg>
  </button>

  <a class="brand display" href="/">ELIXIR</a>
  <nav>
    <a href="/fragrances" class:active={isActive('/fragrances')}>Catálogo</a>
    <a href="/envios" class:active={isActive('/envios')}>Envíos</a>
    <a href="/contacto" class:active={isActive('/contacto')}>Contacto</a>
    {#if $wishlist.length > 0}<a href="/lista-de-deseos" class:active={isActive('/lista-de-deseos')}>Mi lista de deseos ({$wishlist.length})</a>{/if}
  </nav>

  <button class="cart-btn" type="button" on:click={onCart} aria-label="Abrir carrito ({$cartCount} productos)">
    <svg width="22" height="22" viewBox="0 0 22 22" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M1 1H4.5L6.72 11.69C6.80 12.07 7.17 12.33 7.56 12.33H17.5C17.89 12.33 18.26 12.07 18.34 11.69L20 4H5.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
      <circle cx="8.5" cy="17.5" r="1.5" fill="currentColor" />
      <circle cx="16.5" cy="17.5" r="1.5" fill="currentColor" />
    </svg>
    {#if $cartCount > 0}
      <span class="cart-count">{$cartCount}</span>
    {/if}
  </button>
</header>

<MobileNav open={mobileNavOpen} onClose={() => mobileNavOpen = false} />

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
  .hamburger {
    display: none;
    background: transparent;
    border: 0;
    color: var(--color-text);
    padding: 4px;
    align-items: center;
    justify-content: center;
  }
  .cart-btn {
    border: 0;
    background: transparent;
    color: var(--color-text);
    padding: 4px;
    position: relative;
    display: flex;
    align-items: center;
  }
  .cart-count {
    position: absolute;
    top: -4px;
    right: -6px;
    min-width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--color-gold);
    color: var(--color-bg);
    font-size: 0.68rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 3px;
  }
  @media (max-width: 720px) {
    nav { display: none; }
    .hamburger { display: flex; }
    .brand { font-size: 1.25rem; }
  }
</style>
