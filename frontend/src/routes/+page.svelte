<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { fade } from 'svelte/transition';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { reveal } from '$lib/utils/reveal';
  import type { PageData } from './$types';
  export let data: PageData;

  const DEFAULT_HERO_ROTATION_INTERVAL_MS = 8000;

  type HeroImage = PageData['heroImages'][number];

  let heroImageIndex = 0;
  let collectionImageIndices: number[] = [];
  let heroCycleInterval: ReturnType<typeof setInterval> | undefined;
  let collectionCycleInterval: ReturnType<typeof setInterval> | undefined;

  $: heroImages = resolveHeroImages(data.homepage, data.heroImages);
  $: activeHeroImage = heroImages[heroImageIndex] ?? heroImages[0];
  $: heroMetaImage = heroImages[0]?.url ?? data.homepage.hero_image_url;
  $: if (heroImages.length > 0 && heroImageIndex >= heroImages.length) {
    heroImageIndex = 0;
  }
  $: if (data.collections) {
    collectionImageIndices = data.collections.map(() => 0);
  }

  onMount(() => {
    if (data.homepage.hero_image_mode === 'product_covers' && heroImages.length > 1) {
      heroCycleInterval = setInterval(() => {
        heroImageIndex = (heroImageIndex + 1) % heroImages.length;
      }, heroRotationInterval());
    }

    collectionCycleInterval = setInterval(() => {
      collectionImageIndices = collectionImageIndices.map((current, index) => {
        const total = data.collections[index]?.images.length ?? 1;
        return total > 1 ? (current + 1) % total : 0;
      });
    }, 8000);
  });

  onDestroy(() => {
    if (heroCycleInterval) clearInterval(heroCycleInterval);
    if (collectionCycleInterval) clearInterval(collectionCycleInterval);
  });

  function resolveHeroImages(homepage: PageData['homepage'], productImages: HeroImage[]): HeroImage[] {
    const fallback = homepage.hero_image_url
      ? [{ url: homepage.hero_image_url, alt: 'ELIXIR Exclusive' }]
      : [];
    if (homepage.hero_image_mode !== 'product_covers') {
      return fallback;
    }
    return productImages.length > 0 ? productImages : fallback;
  }

  function heroRotationInterval() {
    const interval = Number(data.homepage.hero_rotation_interval_ms);
    return Number.isFinite(interval) && interval >= 1000
      ? interval
      : DEFAULT_HERO_ROTATION_INTERVAL_MS;
  }
</script>

<svelte:head>
  <title>ELIXIR Exclusive | Lorem ipsum</title>
  <meta name="description" content="Lorem ipsum dolor sit amet, consectetur adipiscing elit." />
  <meta property="og:title" content="ELIXIR Exclusive" />
  <meta property="og:description" content={data.homepage.hero_subheading} />
  <meta property="og:image" content={heroMetaImage} />
</svelte:head>

<section class="hero">
  {#if activeHeroImage}
    {#key activeHeroImage.url}
      <img
        class="hero-img"
        src={activeHeroImage.url}
        alt={activeHeroImage.alt}
        transition:fade={{ duration: 900 }}
      />
    {/key}
  {/if}
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

{#if data.collections.length}
  <section class="container collections-section" use:reveal>
    <p class="eyebrow">Colecciones</p>
    <div class="gold-rule"></div>
    <h2 class="display collections-title">Explorá por familia</h2>
    <div
      class="collections-grid"
      class:single-collection-grid={data.collections.length === 1}
      class:two-collection-grid={data.collections.length === 2}
      class:four-collection-grid={data.collections.length >= 4}
    >
      {#each data.collections as collection, i}
        <a
          class="collection-card"
          class:hero-card={i === 0}
          class:wide-card={i === 3}
          href={collection.href}
        >
          {#key collection.images[collectionImageIndices[i] ?? 0]}
            <img
              src={collection.images[collectionImageIndices[i] ?? 0]}
              alt={collection.family}
              transition:fade={{ duration: 800 }}
            />
          {/key}
          <span class="collection-copy">
            <strong>{collection.family}</strong>
            <em>Explorar →</em>
          </span>
        </a>
      {/each}
    </div>
  </section>
{/if}

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
    will-change: transform, opacity;
  }
  .hero-overlay {
    position: absolute;
    inset: 0;
    background: linear-gradient(
      to top,
      rgba(25,23,20,0.92) 0%,
      rgba(25,23,20,0.55) 40%,
      rgba(25,23,20,0.15) 100%
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
    color: var(--color-on-image);
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
    color: var(--color-on-image);
    text-shadow: 0 2px 12px rgba(0,0,0,0.3);
  }
  .hero .eyebrow {
    color: var(--color-emerald);
    text-shadow: 0 1px 4px rgba(0,0,0,0.3);
  }
  .rule { width: 64px; height: 1px; background: rgba(255,253,248,0.35); margin: 0 0 28px; }
  .actions { display: flex; gap: 16px; align-items: center; }
  .actions .btn.text { padding: 0 14px; color: var(--color-on-image); }
  .actions .btn.text:hover { color: var(--color-emerald); }
  .editorial {
    display: grid;
    grid-template-columns: .9fr 1.1fr;
    gap: 48px;
    align-items: center;
    border-radius: 16px;
    background: var(--color-surface);
    padding: 48px;
    overflow: hidden;
    box-shadow: 0 4px 32px rgba(0,0,0,0.2);
  }
  .editorial .copy {
    padding-left: 0;
    padding-right: 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  .editorial p:last-child { color: var(--color-text-muted); font-size: 1.2rem; line-height: 1.8; }
  .editorial img {
    display: block;
    width: 100%;
    max-width: 100%;
    height: auto;
    aspect-ratio: 5 / 4;
    object-fit: cover;
    border-radius: 12px;
  }
  .collections-section {
    margin-top: clamp(64px, 8vw, 104px);
    padding-bottom: 28px;
  }
  .collections-title { margin: 0 0 28px; font-size: clamp(2.2rem, 5vw, 4.4rem); line-height: .96; }
  .collections-grid {
    display: grid;
    grid-template-columns: 1.2fr 0.8fr;
    grid-template-rows: 1fr 1fr;
    gap: 16px;
  }
  .collection-card {
    position: relative;
    display: block;
    overflow: hidden;
    min-height: 260px;
    border-radius: 16px;
  }
  .collection-card::after {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(to top, rgba(25,23,20,0.88) 0%, rgba(25,23,20,0.20) 55%, transparent 100%);
    transition: background .3s ease;
    pointer-events: none;
  }
  .hero-card { grid-row: 1 / 3; min-height: 536px; }
  .collection-card img {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform .6s ease, opacity .8s ease;
  }
  .collection-card:hover img { transform: scale(1.06); }
  .collection-card:hover::after { background: linear-gradient(to top, rgba(25,23,20,0.92) 0%, rgba(25,23,20,0.25) 58%, transparent 100%); }
  .collection-copy {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 1;
    padding: 34px 28px 26px;
    color: var(--color-on-image);
  }
  .collection-copy strong {
    display: block;
    font-family: var(--font-display);
    font-size: clamp(1.6rem, 3vw, 2.8rem);
    line-height: .98;
    font-weight: 600;
    letter-spacing: 0;
    text-shadow: 0 1px 6px rgba(0,0,0,0.4);
  }
  .collection-copy em {
    display: inline-flex;
    margin-top: 12px;
    color: var(--color-emerald);
    font-style: normal;
    font-size: .78rem;
    text-transform: uppercase;
    letter-spacing: .12em;
    border-bottom: 1px solid transparent;
    text-shadow: 0 1px 3px rgba(0,0,0,0.3);
  }
  .collection-card:hover .collection-copy em { border-bottom-color: var(--color-emerald); }
  .four-collection-grid {
    grid-template-columns: 1.1fr 0.75fr 0.9fr;
    grid-template-rows: repeat(2, minmax(260px, 1fr));
    grid-template-areas:
      "hero stack-top wide"
      "hero stack-bottom wide";
  }
  .four-collection-grid .hero-card {
    grid-area: hero;
  }
  .four-collection-grid .collection-card:nth-child(2) {
    grid-area: stack-top;
  }
  .four-collection-grid .collection-card:nth-child(3) {
    grid-area: stack-bottom;
  }
  .four-collection-grid .wide-card {
    grid-area: wide;
    min-height: 536px;
  }
  .single-collection-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto;
  }
  .single-collection-grid .hero-card {
    grid-row: auto;
    min-height: 360px;
  }
  .two-collection-grid {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto;
  }
  .two-collection-grid .hero-card {
    grid-row: auto;
  }
  @media (max-width: 600px) {
    .hero { min-height: 560px; height: 80svh; }
    h1 { font-size: clamp(2.8rem, 12vw, 5rem); }
    .actions { flex-direction: column; align-items: flex-start; }
  }
  @media (max-width: 860px) {
    .editorial {
      grid-template-columns: 1fr;
      padding: 32px;
      gap: 0;
    }
    .editorial .copy {
      padding-right: 0;
      padding-bottom: 32px;
    }
    .editorial img {
      width: 100%;
      height: auto;
      aspect-ratio: 16 / 9;
      border-radius: 12px;
    }
    .collections-grid { grid-template-columns: 1fr; grid-template-rows: auto; }
    .collection-card, .hero-card { grid-row: auto; min-height: 280px; }
    .single-collection-grid,
    .two-collection-grid,
    .four-collection-grid {
      grid-template-columns: 1fr;
      grid-template-rows: auto;
      grid-template-areas: none;
    }
    .four-collection-grid .hero-card,
    .four-collection-grid .collection-card:nth-child(2),
    .four-collection-grid .collection-card:nth-child(3),
    .four-collection-grid .wide-card {
      grid-area: auto;
      grid-column: auto;
      grid-row: auto;
      min-height: 280px;
    }
  }

  @media (prefers-reduced-motion: no-preference) {
    .hero-img { animation: ken 22s ease-in-out infinite alternate; }

    @keyframes ken {
      from { transform: scale(1) translateX(0); }
      to { transform: scale(1.04) translateX(0); }
    }
  }
</style>
