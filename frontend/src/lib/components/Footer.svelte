<script lang="ts">
  import type { SiteSettings } from '$lib/api/client';
  export let settings: SiteSettings;
</script>

<footer class="footer">
  <div class="container footer-grid">
    <div>
      <div class="display logo">{settings.about_title}</div>
      <p class="muted">{settings.about_description}</p>
      <div class="socials" aria-label="Redes sociales">
        {#if settings.footer_instagram_url}<a href={settings.footer_instagram_url} target="_blank" rel="noreferrer" aria-label="Instagram">IG</a>{/if}
        {#if settings.footer_tiktok_url}<a href={settings.footer_tiktok_url} target="_blank" rel="noreferrer" aria-label="TikTok">TT</a>{/if}
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
      <span>© {new Date().getFullYear()} ELIXIR Exclusive</span>
    </nav>
  </div>
</footer>

<style>
  .footer { padding: 54px 0 38px; margin-top: 80px; border-top: 1px solid rgba(184,151,94,.22); }
  .footer-grid { display: grid; grid-template-columns: 1.4fr 1fr 1fr 1fr; gap: 42px; }
  .logo { font-size: 2rem; }
  nav { display: grid; align-content: start; gap: 10px; color: var(--color-text-muted); }
  h3 { color: var(--color-gold); margin: 0 0 8px; font-size: .75rem; text-transform: uppercase; letter-spacing: .14em; }
  span { font-size: .85rem; margin-top: 8px; }
  .socials { display: flex; gap: 10px; margin-top: 18px; min-height: 36px; }
  .socials a { width: 36px; height: 36px; border: 1px solid var(--color-border); display: grid; place-items: center; color: var(--color-gold); font-size: .72rem; letter-spacing: .1em; }
  @media (max-width: 860px) { .footer-grid { grid-template-columns: 1fr; } }
</style>
