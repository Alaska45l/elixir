<footer class="footer">
  <div class="container footer-grid">
    <div>
      <div class="display logo">ELIXIR Exclusive</div>
      <p class="muted">Perfumería argentina de lujo discreto. Fragancias intensas, envíos nacionales y atención privada.</p>
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
      <a href="/contacto">Preguntas frecuentes</a>
      <a href="/contacto">Contacto</a>
    </nav>
    <nav>
      <h3>Contacto</h3>
      <a href="/contacto">Consultar por WhatsApp</a>
      <a href="/admin">Panel administrativo</a>
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
  @media (max-width: 860px) { .footer-grid { grid-template-columns: 1fr; } }
</style>
