<script lang="ts">
  import { onMount } from 'svelte';
  import { fade, fly } from 'svelte/transition';
  import SceneMap from '$lib/components/home/SceneMap.svelte';
  import SecurePackagingScene from '$lib/components/home/scenes/SecurePackagingScene.svelte';
  import OriginGuaranteeScene from '$lib/components/home/scenes/OriginGuaranteeScene.svelte';

  const AUTO_ADVANCE_MS = 6500;

  const scenes = [
    {
      id: 'national-reach',
      label: 'Alcance',
      eyebrow: 'Alcance Nacional',
      title: 'Llegamos a cada destino.',
      copy: 'Coordinamos envíos desde Buenos Aires hacia todo el país con seguimiento claro, tiempos cuidados y una presentación pensada para llegar impecable.',
      component: SceneMap
    },
    {
      id: 'secure-packaging',
      label: 'Empaque',
      eyebrow: 'Empaque Seguro',
      title: 'Protección para piezas delicadas',
      copy: 'Cada fragancia viaja contenida, estabilizada y protegida para conservar la experiencia desde el primer contacto hasta la apertura.',
      component: SecurePackagingScene
    },
    {
      id: 'origin-guarantee',
      label: 'Origen',
      eyebrow: 'Garantía de Origen',
      title: 'Autenticidad verificada',
      copy: 'Trabajamos solo con perfumes originales y controles previos al envío, para que cada compra llegue con respaldo y trazabilidad.',
      component: OriginGuaranteeScene
    }
  ] as const;

  let currentScene = 0;
  let prefersReducedMotion = false;
  let carouselTimer: ReturnType<typeof setInterval> | undefined;

  $: activeScene = scenes[currentScene] ?? scenes[0];

  onMount(() => {
    const motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    const syncMotionPreference = () => {
      prefersReducedMotion = motionQuery.matches;
      if (prefersReducedMotion) {
        clearTimer();
      } else {
        startTimer();
      }
    };

    syncMotionPreference();
    motionQuery.addEventListener('change', syncMotionPreference);

    return () => {
      motionQuery.removeEventListener('change', syncMotionPreference);
      clearTimer();
    };
  });

  function startTimer() {
    clearTimer();
    if (prefersReducedMotion || scenes.length <= 1) return;

    carouselTimer = setInterval(() => {
      currentScene = (currentScene + 1) % scenes.length;
    }, AUTO_ADVANCE_MS);
  }

  function clearTimer() {
    if (!carouselTimer) return;

    clearInterval(carouselTimer);
    carouselTimer = undefined;
  }

  function selectScene(index: number) {
    if (index < 0 || index >= scenes.length || index === currentScene) return;

    clearTimer();
    currentScene = index;
    startTimer();
  }
</script>

<section class="story-band" class:reduced-motion={prefersReducedMotion} aria-labelledby="story-carousel-title">
  <div class="container story-shell">
    <div class="story-copy">
      {#key activeScene.id}
        <div class="copy-frame" in:fly={{ y: 14, duration: 420, opacity: 0 }}>
          <p class="eyebrow">{activeScene.eyebrow}</p>
          <span class="scene-count">0{currentScene + 1} / 03</span>
          <h2 id="story-carousel-title" class="display">{activeScene.title}</h2>
          <p>{activeScene.copy}</p>
        </div>
      {/key}

      <div class="story-controls" aria-label="Escenas de servicio">
        {#each scenes as scene, index}
          <button
            type="button"
            class:active={index === currentScene}
            aria-label={`Ver ${scene.eyebrow}`}
            aria-current={index === currentScene ? 'step' : undefined}
            on:click={() => selectScene(index)}
          >
            <span>{scene.label}</span>
          </button>
        {/each}
      </div>
    </div>

    <div class="story-stage" aria-live="off">
      {#key activeScene.id}
        <div class="scene-frame" in:fade={{ duration: 360 }}>
          <svelte:component this={activeScene.component} />
        </div>
      {/key}
    </div>
  </div>
</section>

<style>
  .story-band {
    margin: clamp(8px, 2vw, 24px) 0 0;
    padding: clamp(64px, 8vw, 104px) 0;
    background: var(--color-surface);
    color: var(--color-text);
    border-top: 1px solid var(--color-border);
    border-bottom: 1px solid var(--color-border);
    overflow: hidden;
  }

  .story-shell {
    display: grid;
    grid-template-columns: minmax(280px, 0.72fr) minmax(0, 1fr);
    gap: clamp(32px, 6vw, 84px);
    align-items: center;
  }

  .story-copy {
    min-height: 360px;
    display: grid;
    align-content: center;
    gap: 32px;
  }

  .copy-frame {
    display: grid;
    gap: 14px;
  }

  .eyebrow {
    color: var(--color-emerald-dark);
    margin: 0;
  }

  .scene-count {
    width: fit-content;
    color: var(--color-text-muted);
    font-size: .76rem;
    letter-spacing: .16em;
  }

  h2 {
    margin: 0;
    max-width: 11ch;
    color: var(--color-text);
    font-size: clamp(2.7rem, 5vw, 5.6rem);
    line-height: .92;
  }

  p:not(.eyebrow) {
    margin: 0;
    max-width: 48ch;
    color: var(--color-text-muted);
    font-size: clamp(1rem, 1.4vw, 1.18rem);
    line-height: 1.75;
  }

  .story-controls {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 10px;
    max-width: 430px;
  }

  .story-controls button {
    position: relative;
    min-width: 0;
    min-height: 44px;
    border: 1px solid var(--color-border);
    border-radius: 6px;
    padding: 0 12px;
    overflow: hidden;
    background: var(--color-bg);
    color: var(--color-text-muted);
    text-align: center;
  }

  .story-controls button::after {
    content: '';
    position: absolute;
    left: 0;
    bottom: 0;
    width: 100%;
    height: 2px;
    background: var(--color-gold);
    transform: scaleX(0);
    transform-origin: left;
  }

  .story-controls button:hover,
  .story-controls button.active {
    border-color: var(--color-gold-dark);
    color: var(--color-text);
    background: color-mix(in srgb, var(--color-gold) 12%, var(--color-bg));
  }

  .story-controls button.active::after {
    animation: story-progress 6500ms linear both;
  }

  .reduced-motion .story-controls button.active::after {
    transform: scaleX(1);
    animation: none;
  }

  .story-controls span {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: .74rem;
    font-weight: 700;
    letter-spacing: .12em;
    text-transform: uppercase;
  }

  .story-stage {
    min-width: 0;
  }

  .scene-frame {
    min-height: clamp(420px, 52vw, 600px);
    display: grid;
    place-items: center;
  }

  @keyframes story-progress {
    from { transform: scaleX(0); }
    to { transform: scaleX(1); }
  }

  @media (max-width: 920px) {
    .story-shell {
      grid-template-columns: 1fr;
    }

    .story-copy {
      min-height: 0;
    }

    h2 {
      max-width: 13ch;
    }

    .scene-frame {
      min-height: clamp(360px, 70vw, 560px);
    }
  }

  @media (max-width: 520px) {
    .story-band {
      padding: 54px 0;
    }

    .story-controls {
      grid-template-columns: 1fr;
      max-width: none;
    }

    .story-controls button {
      min-height: 42px;
    }
  }
</style>
