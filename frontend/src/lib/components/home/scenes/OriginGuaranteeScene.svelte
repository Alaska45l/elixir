<script lang="ts">
  import { draw, fade, fly } from 'svelte/transition';
</script>

<svg class="origin-scene" viewBox="0 0 560 560" role="img" aria-labelledby="origin-title origin-desc">
  <title id="origin-title">Garantía de Origen</title>
  <desc id="origin-desc">Sello de garantía minimalista.</desc>

  <rect class="scene-bg" x="0" y="0" width="560" height="560" rx="18" />
  
  <g class="badge">
    <circle class="ring-outer" cx="280" cy="240" r="140" in:draw={{ duration: 1000, delay: 100 }} />
    <circle class="ring-inner" cx="280" cy="240" r="125" in:draw={{ duration: 1000, delay: 250 }} />
    <path class="sharp-check" d="M 230 250 L 265 285 L 340 195" in:draw={{ duration: 600, delay: 800 }} />
  </g>

  <g class="typography" in:fly={{ y: 20, duration: 800, delay: 1200, opacity: 0 }}>
    <text class="percent" x="280" y="445" text-anchor="middle">100%</text>
    <text class="title" x="280" y="492" text-anchor="middle">ORIGINAL</text>
    <text class="subtitle" x="280" y="523" text-anchor="middle">Autenticidad Verificada</text>
    <path class="separator" d="M 250 542 H 310" />
  </g>
</svg>

<style>
  .origin-scene { width: min(100%, 640px); height: auto; display: block; }
  .scene-bg { fill: color-mix(in srgb, var(--color-bg) 74%, var(--color-surface)); stroke: var(--color-border); }
  
  .ring-outer { fill: none; stroke: var(--color-gold); stroke-width: 2; opacity: 0.6; }
  .ring-inner { fill: color-mix(in srgb, var(--color-surface) 60%, transparent); stroke: var(--color-emerald); stroke-width: 4; }
  .sharp-check { fill: none; stroke: var(--color-emerald-dark); stroke-width: 12; stroke-linecap: round; stroke-linejoin: round; }
  
  .percent { fill: var(--color-text); font-family: var(--font-display); font-size: 68px; font-weight: 700; letter-spacing: 0; }
  .title { fill: var(--color-gold-dark); font-family: var(--font-display); font-size: 26px; font-weight: 700; letter-spacing: 0.24em; }
  .subtitle { fill: var(--color-text-muted); font-size: 14px; font-weight: 600; letter-spacing: 0.2em; text-transform: uppercase; }
  .separator { fill: none; stroke: var(--color-gold-dark); stroke-width: 2; opacity: 0.5; }
</style>
