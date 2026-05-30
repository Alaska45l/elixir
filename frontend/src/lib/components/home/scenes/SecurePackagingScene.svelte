<script lang="ts">
  import { draw, fade, fly } from 'svelte/transition';
</script>

<svg class="packaging-scene" viewBox="0 0 560 520" role="img" aria-labelledby="packaging-title packaging-desc">
  <title id="packaging-title">Empaque Seguro</title>
  <desc id="packaging-desc">Una botella de perfume desciende y entra perfectamente en una caja.</desc>
  <defs>
    <clipPath id="box-interior">
      <rect x="160" y="0" width="240" height="460" />
    </clipPath>
  </defs>

  <rect class="scene-bg" x="0" y="0" width="560" height="520" rx="18" />
  <path class="measurement-line" d="M 114 92 H 446 M 132 142 H 428 M 114 420 H 446" in:draw={{ duration: 900, delay: 120 }} />

  <g class="bottle-wrapper" clip-path="url(#box-interior)">
    <g class="bottle" transform="scale(0.55) translate(335, 365)" in:fly={{ y: -250, duration: 1100, delay: 200, opacity: 0 }}>
      <rect class="cap" x="250" y="100" width="60" height="36" rx="4" />
      <rect class="neck" x="265" y="136" width="30" height="14" />
      <rect class="bottle-body" x="210" y="150" width="140" height="170" rx="12" />
      <rect class="label" x="235" y="210" width="90" height="60" rx="4" />
      <text class="brand-text" x="280" y="246" text-anchor="middle">ELIXIR</text>
    </g>
  </g>

  <g class="box-front" in:fade={{ duration: 600, delay: 100 }}>
    <path class="box-body" d="M 160 250 H 400 V 460 H 160 Z" />
    <path class="box-line" d="M 160 270 H 400" />
  </g>

  <g class="flaps">
    <line class="flap left-flap" x1="160" y1="250" x2="160" y2="130" />
    <line class="flap right-flap" x1="400" y1="250" x2="400" y2="130" />
  </g>

  <g class="seal" in:fade={{ duration: 400, delay: 1400 }}>
    <circle cx="280" cy="250" r="24" />
    <path class="seal-check" d="M 270 250 L 278 258 L 292 242" in:draw={{ duration: 400, delay: 1600 }} />
  </g>
</svg>

<style>
  .packaging-scene { width: min(100%, 660px); height: auto; display: block; }
  .scene-bg { fill: color-mix(in srgb, var(--color-bg) 74%, var(--color-surface)); stroke: var(--color-border); }
  .measurement-line { fill: none; stroke: color-mix(in srgb, var(--color-border) 70%, transparent); stroke-linecap: round; }

  .cap { fill: var(--color-text); }
  .neck { fill: var(--color-gold-dark); }
  .bottle-body { fill: color-mix(in srgb, var(--color-bg) 82%, var(--color-emerald) 18%); stroke: var(--color-emerald-dark); stroke-width: 1.5; }
  .label { fill: var(--color-surface); stroke: var(--color-gold); stroke-width: 1.5; }
  .brand-text { fill: var(--color-text); font-size: 13px; font-weight: 700; letter-spacing: 0.2em; }

  .box-body { fill: color-mix(in srgb, var(--color-surface) 80%, var(--color-gold) 20%); stroke: var(--color-gold-dark); stroke-width: 1.5; }
  .box-line { fill: none; stroke: rgba(212,168,75,0.42); stroke-width: 1.5; }

  .flap { stroke: var(--color-gold-dark); stroke-width: 4; stroke-linecap: round; transform-box: fill-box; }
  .left-flap { transform-origin: bottom; animation: close-left 0.8s cubic-bezier(0.22, 1, 0.36, 1) 1.2s both; }
  .right-flap { transform-origin: bottom; animation: close-right 0.8s cubic-bezier(0.22, 1, 0.36, 1) 1.3s both; }

  .seal circle { fill: var(--color-emerald); stroke: var(--color-bg); stroke-width: 2; }
  .seal-check { fill: none; stroke: var(--color-bg); stroke-width: 3; stroke-linecap: round; stroke-linejoin: round; }

  @keyframes close-left { from { transform: rotate(0deg); } to { transform: rotate(90deg); } }
  @keyframes close-right { from { transform: rotate(0deg); } to { transform: rotate(-90deg); } }
</style>
