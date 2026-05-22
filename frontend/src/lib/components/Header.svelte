<script lang="ts">
  import { page } from '$app/stores';
  import MobileNav from '$lib/components/MobileNav.svelte';
  import type { SiteSettings } from '$lib/api/client';
  import { cartCount } from '$lib/stores/cart';
  import { wishlist } from '$lib/stores/wishlist';
  import { onMount } from 'svelte';

  export let settings: SiteSettings;
  export let onCart: () => void = () => undefined;

  const contactLinks = [
    { label: 'Política de devolución', href: '/politica-de-devolucion' },
    { label: 'Preguntas frecuentes', href: '/preguntas-frecuentes' },
    { label: 'Envíos', href: '/envios' },
    { label: 'Contacta con nosotros', href: '/contacto' }
  ];

  let scrolled = false;
  let mobileNavOpen = false;
  let activePanel: 'productos' | 'contacto' | 'nosotros' | '' = '';
  let closeTimer: ReturnType<typeof setTimeout> | undefined;

  $: activePath = $page.url.pathname;
  $: panelTitle = activePanel === 'productos' ? 'Productos' : 'Contacto';
  $: panelLinks = activePanel === 'productos' ? settings.navbar_product_categories : activePanel === 'contacto' ? contactLinks : [];
  $: instagramHref = settings.footer_instagram_url || 'https://www.instagram.com/';
  $: tiktokHref = settings.footer_tiktok_url || 'https://www.tiktok.com/';

  function isActive(href: string): boolean {
    const path = href.split('?')[0];
    return path === '/' ? activePath === '/' : activePath.startsWith(path);
  }

  function openPanel(panel: 'productos' | 'contacto' | 'nosotros') {
    clearTimeout(closeTimer);
    activePanel = panel;
  }

  function scheduleClose() {
    clearTimeout(closeTimer);
    closeTimer = setTimeout(() => activePanel = '', 150);
  }

  onMount(() => {
    const onScroll = () => scrolled = window.scrollY > 48;
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  });
</script>

<button class="menu-backdrop" class:open={activePanel !== ''} type="button" aria-label="Cerrar menú" on:click={() => activePanel = ''}></button>

<header class="site-header" class:scrolled class:menu-open={activePanel !== ''}>
  <button class="hamburger" type="button" on:click={() => mobileNavOpen = true} aria-label="Abrir menú" aria-expanded={mobileNavOpen} aria-controls="mobile-nav">
    <svg width="22" height="16" viewBox="0 0 22 16" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <path d="M1 1H21M1 8H21M1 15H21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
    </svg>
  </button>

  <a class="brand display" href="/">ELIXIR</a>
  <nav class="desktop-nav" aria-label="Navegación principal">
    <div class="menu-group" role="none" on:mouseenter={() => openPanel('productos')} on:mouseleave={scheduleClose} on:focusin={() => openPanel('productos')}>
      <a class:active={isActive('/fragrances')} href="/fragrances" aria-haspopup="true" aria-expanded={activePanel === 'productos'}>Productos</a>
    </div>
    <div class="menu-group" role="none" on:mouseenter={() => openPanel('contacto')} on:mouseleave={scheduleClose} on:focusin={() => openPanel('contacto')}>
      <a class:active={isActive('/contacto') || isActive('/envios')} href="/contacto" aria-haspopup="true" aria-expanded={activePanel === 'contacto'}>Contacto</a>
    </div>
    <div class="menu-group" role="none" on:mouseenter={() => openPanel('nosotros')} on:mouseleave={scheduleClose} on:focusin={() => openPanel('nosotros')}>
      <button type="button" aria-haspopup="true" aria-expanded={activePanel === 'nosotros'}>Nosotros</button>
    </div>
    {#if $wishlist.length > 0}<a href="/lista-de-deseos" class:active={isActive('/lista-de-deseos')}>Mi lista ({$wishlist.length})</a>{/if}
  </nav>

  <button class="cart-btn" type="button" on:click={onCart} aria-label="Abrir carrito ({$cartCount} productos)">
    <svg width="22" height="22" viewBox="0 0 22 22" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <path d="M1 1H4.5L6.72 11.69C6.80 12.07 7.17 12.33 7.56 12.33H17.5C17.89 12.33 18.26 12.07 18.34 11.69L20 4H5.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
      <circle cx="8.5" cy="17.5" r="1.5" fill="currentColor" />
      <circle cx="16.5" cy="17.5" r="1.5" fill="currentColor" />
    </svg>
    {#if $cartCount > 0}
      <span class="cart-count">{$cartCount}</span>
    {/if}
  </button>

  <div class="mega-panel" class:open={activePanel !== ''} role="region" aria-label="Menú desplegable" on:mouseenter={() => clearTimeout(closeTimer)} on:mouseleave={scheduleClose}>
    <div class="mega-inner" class:nosotros-panel={activePanel === 'nosotros'}>
      {#if activePanel === 'nosotros'}
        <section class="mega-about" aria-label="Nosotros">
          <p class="panel-label">Nosotros</p>
          <h2>{settings.about_title}</h2>
          <p>{settings.about_description}</p>
          {#if settings.about_location}<span>{settings.about_location}</span>{/if}
        </section>

        <section class="mega-social" aria-label="Redes sociales">
          <p class="panel-label">Social</p>
          <a href={instagramHref} target="_blank" rel="noreferrer">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="4" y="4" width="16" height="16" rx="5" stroke="currentColor" stroke-width="1.6"/><circle cx="12" cy="12" r="3.4" stroke="currentColor" stroke-width="1.6"/><circle cx="17" cy="7" r="1" fill="currentColor"/></svg>
            Instagram
          </a>
          <a href={tiktokHref} target="_blank" rel="noreferrer">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M14 4v10.2a3.8 3.8 0 1 1-3.8-3.8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><path d="M14 4c.6 2.8 2.2 4.5 5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            TikTok
          </a>
        </section>
      {:else}
        <section class="mega-links" aria-label={panelTitle}>
          <p class="panel-label">{panelTitle}</p>
          {#each panelLinks as item}
            <a href={item.href} on:click={() => activePanel = ''}>{item.label}</a>
          {/each}
        </section>
      {/if}
    </div>
  </div>
</header>

<MobileNav id="mobile-nav" open={mobileNavOpen} {settings} contactLinks={contactLinks} onClose={() => mobileNavOpen = false} />

<style>
  .menu-backdrop { position: fixed; inset: 0; z-index: 18; border: 0; padding: 0; background: color-mix(in srgb, var(--color-bg) 22%, transparent); opacity: 0; visibility: hidden; pointer-events: none; backdrop-filter: blur(0); transition: opacity .24s ease, visibility 0s linear .24s, backdrop-filter .24s ease; }
  .menu-backdrop.open { opacity: 1; visibility: visible; pointer-events: auto; backdrop-filter: blur(14px); transition-delay: 0s; }
  .site-header { position: sticky; top: 0; z-index: 32; height: 72px; padding: 0 max(20px, 5vw); display: flex; align-items: center; justify-content: space-between; background: color-mix(in srgb, var(--color-bg) 0%, transparent); border-bottom: 1px solid transparent; transition: background .3s ease, border-color .3s ease, backdrop-filter .3s ease, box-shadow .3s ease; }
  .site-header.scrolled, .site-header.menu-open { background: color-mix(in srgb, var(--color-bg) 92%, transparent); border-bottom-color: var(--color-border); backdrop-filter: blur(16px); box-shadow: 0 1px 0 var(--color-border), 0 8px 24px color-mix(in srgb, var(--color-bg) 40%, transparent); }
  .brand { color: var(--color-text); font-size: 1.55rem; }
  .desktop-nav { display: flex; gap: 26px; align-items: center; color: var(--color-text-muted); font-size: .92rem; }
  .desktop-nav a, .desktop-nav button { position: relative; background: transparent; border: 0; padding: 0; color: inherit; }
  .desktop-nav a::after, .desktop-nav button::after { content: ''; position: absolute; left: 0; bottom: -4px; width: 0; height: 1px; background: var(--color-gold); transition: width 0.25s ease; }
  .desktop-nav a.active { color: var(--color-gold); }
  .desktop-nav a.active::after, .desktop-nav a:hover::after, .desktop-nav button:hover::after, .menu-group:hover > a::after, .menu-group:hover > button::after, .desktop-nav [aria-expanded='true']::after { width: 100%; }
  .desktop-nav a:hover, .desktop-nav button:hover, .menu-group:hover > a, .menu-group:hover > button, .desktop-nav [aria-expanded='true'] { color: var(--color-gold); }
  .menu-group { padding: 26px 0; }
  .mega-panel { position: absolute; top: 100%; left: 50%; width: 100vw; margin-left: -50vw; background: color-mix(in srgb, var(--color-surface) 94%, rgba(232,224,208,.12)); border-top: 1px solid var(--color-border); border-bottom: 1px solid var(--color-border); box-shadow: 0 28px 70px rgba(0,0,0,.32); opacity: 0; visibility: hidden; transform: translateY(-18px); pointer-events: none; transition: opacity .2s ease, transform .26s cubic-bezier(.22, 1, .36, 1), visibility 0s linear .26s; }
  .mega-panel.open { opacity: 1; visibility: visible; transform: translateY(0); pointer-events: auto; transition-delay: 0s; }
  .mega-inner { width: min(100%, 1440px); margin: 0 auto; min-height: 236px; padding: 34px max(24px, 5vw) 38px; display: grid; grid-template-columns: minmax(280px, 420px); gap: clamp(36px, 8vw, 120px); align-items: start; }
  .mega-inner.nosotros-panel { grid-template-columns: minmax(320px, 1fr) minmax(210px, .45fr); }
  .panel-label { margin: 0 0 20px; color: var(--color-gold); font-size: .68rem; text-transform: uppercase; letter-spacing: .16em; font-weight: 700; }
  .mega-about h2 { margin: 0 0 18px; font-family: var(--font-display); font-size: clamp(2rem, 4vw, 3.7rem); line-height: .95; font-weight: 600; }
  .mega-about p:not(.panel-label) { margin: 0 0 26px; color: var(--color-text); max-width: 46ch; line-height: 1.55; font-size: .98rem; }
  .mega-about span { display: block; color: var(--color-text-muted); font-size: .95rem; }
  .mega-links, .mega-social { display: grid; align-content: start; gap: 8px; }
  .mega-links a, .mega-social a { display: flex; align-items: center; gap: 10px; min-height: 34px; color: var(--color-text); font-size: 1.03rem; }
  .mega-links a::before { content: '+'; color: var(--color-gold); font-size: .96rem; }
  .mega-links a::after, .mega-social a::after { display: none; }
  .mega-links a:hover, .mega-social a:hover { color: var(--color-gold); }
  .mega-social svg { color: var(--color-gold); flex: 0 0 auto; }
  .hamburger { display: none; background: transparent; border: 0; color: var(--color-text); padding: 4px; align-items: center; justify-content: center; }
  .cart-btn { border: 0; background: transparent; color: var(--color-text); padding: 4px; position: relative; display: flex; align-items: center; }
  .cart-count { position: absolute; top: -4px; right: -6px; min-width: 16px; height: 16px; border-radius: 50%; background: var(--color-gold); color: var(--color-bg); font-size: 0.68rem; font-weight: 700; display: flex; align-items: center; justify-content: center; padding: 0 3px; }
  @media (max-width: 720px) {
    .menu-backdrop { display: none; }
    .mega-panel { display: none; }
    .desktop-nav { display: none; }
    .hamburger { display: flex; }
    .brand { font-size: 1.25rem; }
  }
</style>
