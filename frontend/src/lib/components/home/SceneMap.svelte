<script lang="ts">
  import { draw, fade } from 'svelte/transition';
  import { origin, provincePaths, routes, VIEWBOX_HEIGHT, VIEWBOX_WIDTH } from '$lib/data/argentina-map';
</script>

<svg class="scene-map" viewBox={`0 0 ${VIEWBOX_WIDTH} ${VIEWBOX_HEIGHT}`} role="img" aria-labelledby="map-title map-desc">
  <title id="map-title">Alcance Nacional</title>
  <desc id="map-desc">Mapa de Argentina con rutas desde Buenos Aires hacia distintas provincias.</desc>
  <defs>
    <filter id="map-soft-shadow" x="-30%" y="-30%" width="160%" height="160%">
      <feDropShadow dx="0" dy="16" stdDeviation="14" flood-color="var(--color-text)" flood-opacity=".14" />
    </filter>
    <marker id="route-arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto" markerUnits="strokeWidth">
      <path d="M 0 0 L 6 3 L 0 6 z" fill="var(--color-gold-dark)" />
    </marker>
  </defs>

  <rect class="scene-bg" x="0" y="0" width={VIEWBOX_WIDTH} height={VIEWBOX_HEIGHT} rx="14" />
  <path class="grid-line" d="M 48 104 H 472 M 48 204 H 472 M 48 304 H 472 M 48 404 H 472" />
  <path class="grid-line vertical" d="M 126 52 V 508 M 232 52 V 508 M 338 52 V 508 M 444 52 V 508" />

  <g class="map-layer" filter="url(#map-soft-shadow)">
    {#each provincePaths as province, index (province.id)}
      <path
        class="province"
        d={province.path}
        in:draw={{ duration: 850, delay: index * 28 }}
      />
    {/each}
  </g>

  <g class="route-layer">
    <circle class="origin-ring" cx={origin.x} cy={origin.y} r="12" in:draw={{ duration: 600, delay: 280 }} />
    <circle class="origin-dot" cx={origin.x} cy={origin.y} r="4.8" in:fade={{ duration: 360, delay: 460 }} />

    {#each routes as route, index}
      <path
        class="route"
        d={route.path}
        marker-end="url(#route-arrow)"
      />
      <circle
        class="destination-dot"
        cx={route.target.x}
        cy={route.target.y}
        r="4.2"
        in:fade={{ duration: 260, delay: 980 + index * 120 }}
      />
      <text
        class="destination-label"
        x={route.target.x + (route.target.x > origin.x ? 12 : -12)}
        y={route.target.y - 10}
        text-anchor={route.target.x > origin.x ? 'start' : 'end'}
        in:fade={{ duration: 300, delay: 1060 + index * 120 }}
      >
        {route.name}
      </text>
    {/each}
  </g>

  <g class="caption" in:fade={{ duration: 420, delay: 1050 }}>
    <text x="354" y="438">Buenos Aires</text>
    <path d="M 340 448 H 468" />
    <text class="small" x="354" y="474">Despacho nacional</text>
  </g>
</svg>

<style>
  .scene-map {
    width: min(100%, 640px);
    height: auto;
    display: block;
  }

  .scene-bg {
    fill: color-mix(in srgb, var(--color-bg) 74%, var(--color-surface));
    stroke: var(--color-border);
  }

  .grid-line {
    fill: none;
    stroke: color-mix(in srgb, var(--color-border) 62%, transparent);
    stroke-width: 1;
  }

  .vertical {
    opacity: .74;
  }

  .province {
    fill: color-mix(in srgb, var(--color-surface) 78%, var(--color-emerald) 22%);
    stroke: color-mix(in srgb, var(--color-text-muted) 50%, transparent);
    stroke-width: .85;
    vector-effect: non-scaling-stroke;
  }

  .route {
    fill: none;
    stroke: var(--color-gold-dark);
    stroke-width: 2.2;
    stroke-linecap: round;
    stroke-dasharray: 8 8;
    animation: march-ants 30s linear infinite;
  }

  @keyframes march-ants {
    from { stroke-dashoffset: 400; }
    to { stroke-dashoffset: 0; }
  }

  .origin-ring {
    fill: color-mix(in srgb, var(--color-emerald) 18%, transparent);
    stroke: var(--color-emerald-dark);
    stroke-width: 1.8;
  }

  .origin-dot {
    fill: var(--color-emerald-dark);
    stroke: var(--color-bg);
    stroke-width: 1;
  }

  .destination-dot {
    fill: var(--color-gold-dark);
    stroke: var(--color-bg);
    stroke-width: 1;
  }

  .destination-label,
  .caption text {
    fill: var(--color-text);
    font-size: 13px;
    font-weight: 700;
    letter-spacing: .1em;
    text-transform: uppercase;
  }

  .caption path {
    stroke: var(--color-gold-dark);
    stroke-width: 1;
  }

  .caption .small {
    fill: var(--color-text-muted);
    font-size: 10px;
  }
</style>
