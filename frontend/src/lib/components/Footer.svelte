<script lang="ts">
  import type { SiteSettings } from '$lib/api/client';
  export let settings: SiteSettings;
</script>

<footer class="footer">
  <div class="container footer-grid">
    <div class="footer-brand">
      <img class="footer-logo" src="/elixir-logo.webp" alt="ELIXIR Exclusive" />
      <div class="footer-copy">
        <h2 class="footer-title display">ELIXIR Exclusive</h2>
        <p class="muted">{settings.about_description}</p>
      </div>
    </div>
    <nav>
      <h3>Catálogo</h3>
      <a href="/fragrances">Fragancias</a>
      <a href="/lista-de-deseos">Mi lista de deseos</a>
      <a href="/fragrances?featured=true">Destacadas</a>
    </nav>
    <nav>
      <h3>Información</h3>
      <a href="/envios">Envíos</a>
      <a href="/preguntas-frecuentes">Preguntas frecuentes</a>
      <a href="/politica-de-devolucion">Política de devolución</a>
      <a href="/contacto">Contacto</a>
    </nav>
    <nav>
      <h3>Contacto</h3>
      {#if settings.footer_whatsapp_url}<a href={settings.footer_whatsapp_url} target="_blank" rel="noreferrer">Consultar por WhatsApp</a>{/if}
      {#if settings.about_phone}<a href={`tel:${settings.about_phone}`}>{settings.about_phone}</a>{/if}
      {#if settings.about_location}<span>{settings.about_location}</span>{/if}
      <a class="social-link" href="https://instagram.com/elixirexclusive" target="_blank" rel="noopener noreferrer">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <rect x="4.5" y="4.5" width="15" height="15" rx="4.5" stroke="currentColor" stroke-width="1.5" />
          <circle cx="12" cy="12" r="3.5" stroke="currentColor" stroke-width="1.5" />
          <circle cx="17" cy="7" r="1" fill="currentColor" />
        </svg>
        Instagram
      </a>
      <a class="social-link" href="https://tiktok.com/@elixirexclusive" target="_blank" rel="noopener noreferrer">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path d="M14 4v10.25a4.25 4.25 0 1 1-4.25-4.25" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M14 4c.7 3.1 2.55 4.95 5.4 5.55" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
        TikTok
      </a>
      <span>© {new Date().getFullYear()} ELIXIR Exclusive</span>
    </nav>
  </div>
</footer>

<style>
  .footer { padding: 54px 0 38px; margin-top: 80px; border-top: 1px solid rgba(184,151,94,.22); }
  .footer-grid { display: grid; grid-template-columns: 1.4fr 1fr 1fr 1fr; gap: 42px; }
  .footer-brand {
    display: grid;
    grid-template-columns: 92px minmax(0, 32ch);
    gap: 20px;
    align-items: center;
  }
  .footer-logo {
    width: 92px;
    height: auto;
    object-fit: contain;
    object-position: left center;
    border-radius: 6px;
  }
  .footer-copy { min-width: 0; }
  .footer-title {
    margin: 0 0 12px;
    color: var(--color-text);
    font-size: 1.55rem;
    line-height: 1;
  }
  nav { display: grid; align-content: start; gap: 10px; color: var(--color-text-muted); }
  h3 { color: var(--color-gold); margin: 0 0 8px; font-size: .75rem; text-transform: uppercase; letter-spacing: .14em; }
  span { font-size: .85rem; margin-top: 8px; }
  .social-link { display: flex; gap: 8px; align-items: center; }
  .social-link svg { color: var(--color-gold); flex: 0 0 auto; }
  @media (max-width: 860px) { .footer-grid { grid-template-columns: 1fr; } }
</style>
