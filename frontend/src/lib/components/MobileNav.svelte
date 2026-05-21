<script lang="ts">
  import { browser } from '$app/environment';
  import { onDestroy, onMount } from 'svelte';

  export let open = false;
  export let onClose: () => void = () => undefined;

  onMount(() => {
    if (open) document.body.style.overflow = 'hidden';
  });

  onDestroy(() => {
    if (browser) document.body.style.overflow = '';
  });

  $: if (browser) {
    document.body.style.overflow = open ? 'hidden' : '';
  }
</script>

{#if open}
  <button class="shade" type="button" on:click={onClose} aria-label="Cerrar menú"></button>
  <nav class="mobile-nav" aria-label="Menú de navegación">
    <div class="nav-header">
      <a class="brand display" href="/" on:click={onClose}>ELIXIR</a>
      <button type="button" on:click={onClose} aria-label="Cerrar">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M2 2L18 18M18 2L2 18" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
        </svg>
      </button>
    </div>
    <ul>
      <li><a href="/fragrances" on:click={onClose}>Catálogo</a></li>
      <li><a href="/envios" on:click={onClose}>Envíos</a></li>
      <li><a href="/contacto" on:click={onClose}>Contacto</a></li>
      <li><a href="/lista-de-deseos" on:click={onClose}>Mi lista de deseos</a></li>
    </ul>
    <div class="nav-footer">
      <p class="eyebrow">ELIXIR Exclusive</p>
      <p class="muted">Perfumería argentina de lujo discreto.</p>
    </div>
  </nav>
{/if}

<style>
  .shade {
    position: fixed;
    inset: 0;
    background: var(--color-overlay);
    z-index: 50;
    border: 0;
    padding: 0;
  }

  .mobile-nav {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: min(320px, 85vw);
    background: var(--color-surface);
    border-right: 1px solid var(--color-border);
    z-index: 51;
    padding: 28px;
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .nav-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 48px;
  }

  .nav-header .brand {
    font-size: 1.55rem;
    color: var(--color-text);
  }

  .nav-header button {
    background: transparent;
    border: 0;
    color: var(--color-text-muted);
    padding: 4px;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    gap: 0;
    flex: 1;
  }

  ul li {
    border-bottom: 1px solid var(--color-border);
  }

  ul li a {
    display: block;
    padding: 18px 0;
    font-family: var(--font-display);
    font-size: 1.9rem;
    font-weight: 600;
    color: var(--color-text);
    letter-spacing: -0.02em;
    line-height: 1;
  }

  ul li a:hover {
    color: var(--color-gold);
  }

  .nav-footer {
    margin-top: 40px;
  }

  .nav-footer .muted {
    font-size: 0.82rem;
    line-height: 1.6;
    margin: 6px 0 0;
  }

  @media (prefers-reduced-motion: no-preference) {
    .shade {
      animation: fadeIn 0.22s ease both;
    }

    .mobile-nav {
      animation: slideInLeft 0.3s cubic-bezier(0.22, 1, 0.36, 1) both;
    }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    @keyframes slideInLeft {
      from { transform: translateX(-100%); }
      to { transform: none; }
    }
  }
</style>
