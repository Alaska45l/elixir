<script lang="ts">
  import { browser } from '$app/environment';
  import type { NavItem, SiteSettings } from '$lib/api/client';
  import { onDestroy } from 'svelte';

  export let id = 'mobile-nav';
  export let open = false;
  export let settings: SiteSettings;
  export let productLinks: NavItem[] = [];
  export let contactLinks: NavItem[] = [];
  export let onClose: () => void = () => undefined;

  let expanded = 'productos';
  $: productNavLinks = productLinks.length > 0 ? productLinks : settings.navbar_product_categories;
  $: instagramHref = settings.footer_instagram_url || 'https://www.instagram.com/';
  $: tiktokHref = settings.footer_tiktok_url || 'https://www.tiktok.com/';

  function toggle(group: string) {
    expanded = expanded === group ? '' : group;
  }

  onDestroy(() => {
    if (browser) document.body.style.overflow = '';
  });

  $: if (browser) {
    document.body.style.overflow = open ? 'hidden' : '';
  }
</script>

<button class="shade" class:open type="button" on:click={onClose} aria-label="Cerrar menú" tabindex={open ? 0 : -1}></button>
<nav {id} class="mobile-nav" class:open aria-label="Menú de navegación" aria-hidden={!open}>
  <div class="nav-header">
    <a class="brand display" href="/" on:click={onClose}>ELIXIR</a>
    <button type="button" on:click={onClose} aria-label="Cerrar">
      <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
        <path d="M2 2L18 18M18 2L2 18" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
      </svg>
    </button>
  </div>
  <ul>
    <li>
      <button class="group-trigger" type="button" aria-expanded={expanded === 'productos'} aria-controls="mobile-productos" on:click={() => toggle('productos')}>Productos</button>
      <div id="mobile-productos" class="children" class:expanded={expanded === 'productos'}>
        {#each productNavLinks as item}<a href={item.href} on:click={onClose}>{item.label}</a>{/each}
      </div>
    </li>
    <li>
      <button class="group-trigger" type="button" aria-expanded={expanded === 'contacto'} aria-controls="mobile-contacto" on:click={() => toggle('contacto')}>Contacto</button>
      <div id="mobile-contacto" class="children" class:expanded={expanded === 'contacto'}>
        {#each contactLinks as item}<a href={item.href} on:click={onClose}>{item.label}</a>{/each}
      </div>
    </li>
    <li>
      <button class="group-trigger" type="button" aria-expanded={expanded === 'nosotros'} aria-controls="mobile-nosotros" on:click={() => toggle('nosotros')}>Nosotros</button>
      <div id="mobile-nosotros" class="children about" class:expanded={expanded === 'nosotros'}>
        <strong>{settings.about_title}</strong>
        <p>{settings.about_description}</p>
        {#if settings.about_location}<span>{settings.about_location}</span>{/if}
        <span class="social-label">Social</span>
        <div class="mobile-social">
          <a class="social-link" href={instagramHref} target="_blank" rel="noreferrer">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="4" y="4" width="16" height="16" rx="5" stroke="currentColor" stroke-width="1.6"/><circle cx="12" cy="12" r="3.4" stroke="currentColor" stroke-width="1.6"/><circle cx="17" cy="7" r="1" fill="currentColor"/></svg>
            Instagram
          </a>
          <a class="social-link" href={tiktokHref} target="_blank" rel="noreferrer">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M14 4v10.2a3.8 3.8 0 1 1-3.8-3.8" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/><path d="M14 4c.6 2.8 2.2 4.5 5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>
            TikTok
          </a>
        </div>
      </div>
    </li>
    <li><a class="solo" href="/lista-de-deseos" on:click={onClose}>Mi lista de deseos</a></li>
  </ul>
</nav>

<style>
  .shade {
    position: fixed;
    inset: 0;
    background: var(--color-overlay);
    backdrop-filter: blur(10px);
    z-index: 50;
    border: 0;
    padding: 0;
    opacity: 0;
    visibility: hidden;
    transition: opacity .22s ease, visibility 0s linear .22s;
  }

  .shade.open {
    opacity: 1;
    visibility: visible;
    transition-delay: 0s;
  }

  .mobile-nav {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: min(430px, 100vw);
    background: var(--color-surface);
    border-right: 1px solid var(--color-border);
    z-index: 51;
    padding: 30px max(24px, 6vw);
    display: flex;
    flex-direction: column;
    gap: 0;
    transform: translateX(-100%);
    visibility: hidden;
    transition: transform .3s cubic-bezier(.22, 1, .36, 1), visibility 0s linear .3s;
  }

  .mobile-nav.open {
    transform: translateX(0);
    visibility: visible;
    transition-delay: 0s;
  }

  .nav-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 44px;
  }

  .nav-header .brand {
    font-size: 1.7rem;
    color: var(--color-text);
  }

  .nav-header button,
  .group-trigger {
    background: transparent;
    border: 0;
    color: var(--color-text);
    padding: 0;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    gap: 0;
    flex: 1;
  }

  li {
    border-bottom: 1px solid var(--color-border);
  }

  .group-trigger,
  .solo {
    width: 100%;
    display: flex;
    justify-content: space-between;
    min-height: 70px;
    padding: 20px 0;
    font-family: var(--font-display);
    font-size: clamp(2.1rem, 9vw, 3rem);
    font-weight: 600;
    color: var(--color-text);
    line-height: 1;
    text-align: left;
  }

  .group-trigger::after {
    content: '+';
    color: var(--color-gold);
    font-family: var(--font-ui);
    font-size: 1.35rem;
  }

  .group-trigger[aria-expanded='true']::after {
    content: '-';
  }

  .children {
    display: grid;
    gap: 4px;
    max-height: 0;
    overflow: hidden;
    opacity: 0;
    transition: max-height .24s ease, opacity .2s ease;
  }

  .children.expanded {
    max-height: 520px;
    opacity: 1;
    padding-bottom: 14px;
  }

  .children a,
  .children div,
  .children span,
  .children p,
  .children strong {
    color: var(--color-text-muted);
    min-height: 48px;
    padding: 10px 0;
    font-size: 1.12rem;
    line-height: 1.5;
  }

  .children a {
    display: flex;
    align-items: center;
  }

  .children div {
    min-height: 0;
    padding: 0;
  }

  .children p {
    min-height: 0;
    font-size: 1rem;
  }

  .children strong {
    color: var(--color-text);
    padding-top: 0;
    font-size: 1.18rem;
    min-height: 0;
  }

  .social-link {
    gap: 12px;
  }

  .social-label {
    min-height: 0;
    padding: 20px 0 4px;
    color: var(--color-gold);
    font-size: .76rem;
    font-weight: 700;
    letter-spacing: .16em;
    text-transform: uppercase;
  }

  .mobile-social {
    display: grid;
    gap: 2px;
  }

  .social-link svg {
    color: var(--color-gold);
    flex: 0 0 auto;
  }

  .children a:hover,
  .solo:hover {
    color: var(--color-gold);
  }
</style>
