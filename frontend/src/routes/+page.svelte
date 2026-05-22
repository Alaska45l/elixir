<script lang="ts">
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { reveal } from '$lib/utils/reveal';
  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>ELIXIR Exclusive | Perfumería argentina</title>
  <meta name="description" content="Perfumes argentinos de lujo con envíos a todo el país y checkout en ARS." />
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
    <div><h2>Envíos a todo el país</h2><p>Preparación cuidada, seguimiento y despacho desde Buenos Aires.</p></div>
  </div>
</section>

<section class="container collections">
  <a href="/fragrances?family=Oriental"><img src="https://images.unsplash.com/photo-1615634260167-c8cdede054de?auto=format&fit=crop&w=900&q=85" alt="Colección oriental" /><span>Orientales intensos</span></a>
  <a href="/fragrances?family=Fresco"><img src="https://images.unsplash.com/photo-1600612253971-422e7f7faeb6?auto=format&fit=crop&w=900&q=85" alt="Colección fresca" /><span>Frescos de día</span></a>
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
  .editorial { display: grid; grid-template-columns: .9fr 1.1fr; gap: 56px; align-items: center; padding: 56px 0; }
  .editorial .copy { border-left: 1px solid var(--color-gold); padding-left: 32px; }
  .editorial p:last-child { color: var(--color-text-muted); font-size: 1.2rem; line-height: 1.8; }
  .editorial img { width: 100%; aspect-ratio: 5 / 4; object-fit: cover; }
  .shipping-strip { margin: 80px 0; padding: 30px 0; }
  .shipping-strip .container { display: flex; gap: 18px; align-items: center; }
  .shipping-strip span { color: var(--color-gold); font-size: .7rem; opacity: .7; transform: rotate(45deg); display: inline-block; width: 8px; height: 8px; border: 1px solid var(--color-gold); flex-shrink: 0; text-indent: -999px; overflow: hidden; }
  .shipping-strip h2 { margin: 0; font-size: 1.3rem; }
  .shipping-strip p { margin: 4px 0 0; color: var(--color-text-muted); }
  .collections { display: grid; grid-template-columns: 1fr 1fr; gap: 28px; }
  .collections a { position: relative; overflow: hidden; min-height: 360px; }
  .collections img { width: 100%; height: 100%; object-fit: cover; transition: transform .6s ease; }
  .collections a:hover img { transform: scale(1.04); }
  .collections span { position: absolute; left: 0; right: 0; bottom: 0; padding: 40px 24px 22px; background: linear-gradient(to top, color-mix(in srgb, var(--color-bg) 82%, transparent) 0%, transparent 100%); color: var(--color-text); font-family: var(--font-display); font-size: 2.2rem; font-weight: 500; letter-spacing: 0; }
  @media (max-width: 600px) {
    .hero { min-height: 560px; height: 80svh; }
    h1 { font-size: clamp(2.8rem, 12vw, 5rem); }
    .actions { flex-direction: column; align-items: flex-start; }
  }
  @media (max-width: 860px) {
    .editorial, .collections { grid-template-columns: 1fr; }
  }

  @media (prefers-reduced-motion: no-preference) {
    .hero-img { animation: ken 22s ease-in-out infinite alternate; }

    @keyframes ken {
      from { transform: scale(1) translateX(0); }
      to { transform: scale(1.04) translateX(0); }
    }
  }
</style>
