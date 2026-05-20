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
  <div class="hero-copy">
    <span class="watermark" aria-hidden="true">ELIXIR</span>
    <div class="hero-content">
      <p class="eyebrow">ELIXIR Exclusive</p>
      <h1 class="display">{data.homepage.hero_heading}</h1>
      <div class="rule"></div>
      <p>{data.homepage.hero_subheading}</p>
      <div class="actions">
        <a class="btn primary" href={data.homepage.hero_cta_url}>{data.homepage.hero_cta_label}</a>
        <a class="btn text" href="/contacto">Consultar por WhatsApp</a>
      </div>
    </div>
  </div>
  <div class="hero-image"><img src={data.homepage.hero_image_url} alt="Perfume ELIXIR Exclusive" /></div>
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
  .hero { min-height: calc(100svh - 72px); display: grid; grid-template-columns: minmax(0, .88fr) minmax(0, 1.12fr); }
  .hero-copy { padding: clamp(40px, 8vw, 120px) clamp(24px, 6vw, 96px); display: flex; flex-direction: column; justify-content: center; gap: 0; position: relative; overflow: hidden; }
  .hero-content { position: relative; z-index: 1; }
  .watermark { position: absolute; top: 50%; left: -10%; transform: translateY(-50%); font-family: var(--font-display); font-size: clamp(8rem, 22vw, 24rem); font-weight: 600; color: transparent; -webkit-text-stroke: 1px color-mix(in srgb, var(--color-gold) 6%, transparent); pointer-events: none; user-select: none; white-space: nowrap; letter-spacing: -.04em; z-index: 0; }
  h1 { font-size: clamp(4rem, 8vw, 8.8rem); line-height: .85; margin: 18px 0; max-width: 760px; }
  .rule { width: 82px; height: 1px; background: var(--color-gold); margin: 12px 0 28px; }
  .hero-copy p:not(.eyebrow) { color: var(--color-text-muted); max-width: 560px; font-size: 1.12rem; line-height: 1.7; }
  .actions { display: flex; gap: 18px; margin-top: 28px; align-items: center; }
  .hero-image { position: relative; overflow: hidden; min-height: 520px; }
  .hero-image img { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; will-change: transform; }
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
  .collections span { position: absolute; left: 0; right: 0; bottom: 0; padding: 40px 24px 22px; background: linear-gradient(to top, color-mix(in srgb, var(--color-bg) 82%, transparent) 0%, transparent 100%); color: var(--color-text); font-family: var(--font-display); font-size: 2.2rem; font-weight: 500; letter-spacing: -.01em; }
  @media (max-width: 860px) {
    .hero { grid-template-columns: 1fr; position: relative; min-height: calc(100svh - 72px); }
    .hero-image { position: absolute; inset: 0; min-height: 0; }
    .hero-image::after { content: ''; position: absolute; inset: 0; background: linear-gradient(135deg, color-mix(in srgb, var(--color-bg) 94%, transparent) 50%, color-mix(in srgb, var(--color-bg) 62%, transparent) 100%); }
    .hero-copy { position: relative; z-index: 1; padding: 60px 24px; min-height: calc(100svh - 72px); }
    .editorial, .collections { grid-template-columns: 1fr; }
  }

  @media (prefers-reduced-motion: no-preference) {
    .hero-content > * { animation: heroIn .9s ease both; }
    .hero-content > *:nth-child(1) { animation-delay: .1s; }
    .hero-content > *:nth-child(2) { animation-delay: .24s; }
    .hero-content > *:nth-child(3) { animation-delay: .38s; }
    .hero-content > *:nth-child(4) { animation-delay: .5s; }
    .hero-content > *:nth-child(5) { animation-delay: .62s; }
    .hero-image img { animation: ken 22s ease-in-out infinite alternate; }

    @keyframes ken {
      from { transform: scale(1) translateX(0); }
      to { transform: scale(1.07) translateX(-2%); }
    }

    @keyframes heroIn {
      from { opacity: 0; transform: translateY(22px); }
      to { opacity: 1; transform: none; }
    }
  }
</style>
