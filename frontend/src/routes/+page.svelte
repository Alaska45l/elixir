<script lang="ts">
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { reveal } from '$lib/utils/reveal';
  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>ELIXIR Exclusive | Lorem ipsum</title>
  <meta name="description" content="Lorem ipsum dolor sit amet, consectetur adipiscing elit." />
  <meta property="og:title" content="ELIXIR Exclusive" />
  <meta property="og:description" content={data.homepage.hero_subheading} />
  <meta property="og:image" content={data.homepage.hero_image_url} />
</svelte:head>

<section class="hero">
  <img class="hero-img" src={data.homepage.hero_image_url} alt="ELIXIR Exclusive" />
  <div class="hero-overlay"></div>
  <div class="hero-text container">
    <span class="watermark" aria-hidden="true">ELIXIR</span>
    <div class="hero-content">
      <p class="eyebrow">ELIXIR Exclusive</p>
      <h1 class="display">{data.homepage.hero_heading}</h1>
      <div class="rule"></div>
      <div class="actions">
        <a class="btn primary" href={data.homepage.hero_cta_url}>{data.homepage.hero_cta_label}</a>
        <a class="btn text" href="/contacto">Consultar por WhatsApp</a>
      </div>
    </div>
  </div>
</section>

<section class="container page-pad">
  <p class="eyebrow">Fragancias destacadas</p>
  <div class="gold-rule"></div>
  <ProductGrid products={data.featured} />
</section>

<section class="container editorial" use:reveal>
  <div class="copy">
    <p class="eyebrow">{data.homepage.editorial_heading}</p>
    <div class="gold-rule"></div>
    <p>{data.homepage.editorial_body}</p>
  </div>
  <img src={data.homepage.editorial_image_url} alt="Editorial ELIXIR" />
</section>

<section class="shipping-strip hairline">
  <div class="container" use:reveal>
    <span>◇</span>
    <div><h2>Lorem ipsum</h2><p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p></div>
  </div>
</section>

<section class="container collections-section" use:reveal>
  <p class="eyebrow">Colecciones</p>
  <div class="gold-rule"></div>
  <h2 class="display collections-title">Lorem ipsum dolor</h2>
  <div class="collections-grid">
    <a class="collection-card hero-card" href="/fragrances?family=Oriental">
      <img src="https://images.unsplash.com/photo-1563170351-be82bc888aa4?auto=format&fit=crop&w=900&q=85" alt="Lorem ipsum" />
      <span class="collection-copy"><strong>Lorem ipsum</strong><em>Explorar →</em></span>
    </a>
    <a class="collection-card" href="/fragrances?family=Amaderado">
      <img src="https://images.unsplash.com/photo-1624811742200-69166e7b7bcc?auto=format&fit=crop&w=900&q=85" alt="Dolor sit amet" />
      <span class="collection-copy"><strong>Dolor sit amet</strong><em>Explorar →</em></span>
    </a>
    <a class="collection-card" href="/fragrances?family=Gourmand">
      <img src="https://images.unsplash.com/photo-1586495777744-4413f21062fa?auto=format&fit=crop&w=900&q=85" alt="Consectetur adipiscing" />
      <span class="collection-copy"><strong>Consectetur adipiscing</strong><em>Explorar →</em></span>
    </a>
  </div>
</section>

<style>
  .hero {
    position: relative;
    min-height: clamp(560px, 78svh, 820px);
    height: min(86svh, 860px);
    overflow: hidden;
  }
  .hero-img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center 30%;
    will-change: transform;
  }
  .hero-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to top,
      rgba(13,31,21,0.92) 0%,
      rgba(13,31,21,0.5) 40%,
      rgba(13,31,21,0.1) 100%
    );
  }
  .hero-text {
    position: absolute;
    inset: 0;
    left: 0;
    right: 0;
    display: flex;
    align-items: center;
    padding-top: clamp(20px, 6vh, 70px);
  }
  .hero-content {
    max-width: clamp(640px, 48vw, 900px);
    position: relative;
  }
  .watermark { display: none; }
  h1 {
    font-size: clamp(3.8rem, 9vw, 9rem);
    line-height: 0.88;
    margin: 10px 0 24px;
  }
  .rule { width: 64px; height: 1px; background: var(--color-gold); margin: 0 0 28px; }
  .actions { display: flex; gap: 16px; align-items: center; }
  .actions .btn.text { padding: 0 14px; }
  .editorial {
    display: grid;
    grid-template-columns: .9fr 1.1fr;
    gap: 56px;
    align-items: center;
    border-radius: 16px;
    background: var(--color-surface);
    padding: 48px;
    overflow: hidden;
    box-shadow: 0 4px 32px rgba(0,0,0,0.2);
  }
  .editorial .copy { padding-left: 0; }
  .editorial p:last-child { color: var(--color-text-muted); font-size: 1.2rem; line-height: 1.8; }
  .editorial img { width: calc(100% + 48px); aspect-ratio: 5 / 4; object-fit: cover; margin: -48px -48px -48px 0; }
  .shipping-strip { margin: 80px 0; padding: 30px 0; }
  .shipping-strip .container { display: flex; gap: 18px; align-items: center; }
  .shipping-strip span { color: var(--color-gold); font-size: .7rem; opacity: .7; transform: rotate(45deg); display: inline-block; width: 8px; height: 8px; border: 1px solid var(--color-gold); flex-shrink: 0; text-indent: -999px; overflow: hidden; }
  .shipping-strip h2 { margin: 0; font-size: 1.3rem; }
  .shipping-strip p { margin: 4px 0 0; color: var(--color-text-muted); }
  .collections-section { padding-bottom: 28px; }
  .collections-title { margin: 0 0 28px; font-size: clamp(2.2rem, 5vw, 4.4rem); line-height: .96; }
  .collections-grid {
    display: grid;
    grid-template-columns: 1.2fr 0.8fr;
    grid-template-rows: 1fr 1fr;
    gap: 16px;
  }
  .collection-card {
    position: relative;
    overflow: hidden;
    min-height: 260px;
    border-radius: 16px;
  }
  .collection-card::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(to top, rgba(13,31,21,0.85) 0%, rgba(13,31,21,0.1) 60%, transparent 100%);
    transition: background .3s ease;
    pointer-events: none;
  }
  .hero-card { grid-row: 1 / 3; min-height: 536px; }
  .collection-card img { width: 100%; height: 100%; object-fit: cover; transition: transform .6s ease; }
  .collection-card:hover img { transform: scale(1.06); }
  .collection-card:hover::after { background: linear-gradient(to top, rgba(13,31,21,0.92) 0%, rgba(13,31,21,0.18) 62%, transparent 100%); }
  .collection-copy {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1;
    padding: 34px 28px 26px;
    color: var(--color-text);
  }
  .collection-copy strong {
    display: block;
    font-family: var(--font-display);
    font-size: clamp(1.6rem, 3vw, 2.8rem);
    line-height: .98;
    font-weight: 600;
    letter-spacing: 0;
  }
  .collection-copy em {
    display: inline-flex;
    margin-top: 12px;
    color: var(--color-gold);
    font-style: normal;
    font-size: .78rem;
    text-transform: uppercase;
    letter-spacing: .12em;
    border-bottom: 1px solid transparent;
  }
  .collection-card:hover .collection-copy em { border-bottom-color: var(--color-gold); }
  @media (max-width: 600px) {
    .hero { min-height: 560px; height: 80svh; }
    h1 { font-size: clamp(2.8rem, 12vw, 5rem); }
    .actions { flex-direction: column; align-items: flex-start; }
  }
  @media (max-width: 860px) {
    .editorial { grid-template-columns: 1fr; }
    .collections-grid { grid-template-columns: 1fr; grid-template-rows: auto; }
    .collection-card, .hero-card { grid-row: auto; min-height: 280px; }
    .editorial { padding: 32px; }
    .editorial img { margin: 0 -48px -48px; width: calc(100% + 96px); }
  }

  @media (prefers-reduced-motion: no-preference) {
    .hero-img { animation: ken 22s ease-in-out infinite alternate; }

    @keyframes ken {
      from { transform: scale(1) translateX(0); }
      to { transform: scale(1.04) translateX(0); }
    }
  }
</style>
