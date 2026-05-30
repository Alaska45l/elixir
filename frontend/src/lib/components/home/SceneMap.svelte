<script lang="ts">
  import { draw, fade } from 'svelte/transition';
  import { geoPath, geoTransverseMercator } from 'd3-geo';
  import argentinaGeojsonRaw from '$lib/data/argentina.geojson?raw';
  import type { Feature, FeatureCollection, MultiPolygon, Polygon, Position } from 'geojson';

  type Coordinates = [number, number];
  type ProvinceFeature = Feature<Polygon | MultiPolygon, { name?: string }>;
  type ProvinceCollection = FeatureCollection<Polygon | MultiPolygon, { name?: string }>;

  const VIEWBOX_WIDTH = 520;
  const VIEWBOX_HEIGHT = 560;
  const argentina = normalizeCollection(JSON.parse(argentinaGeojsonRaw) as ProvinceCollection);
  const buenosAires: Coordinates = [-58.3816, -34.6037];
  const destinations = [
    { name: 'Salta', coordinates: [-65.4122, -24.7821] as Coordinates },
    { name: 'Misiones', coordinates: [-54.5736, -25.5972] as Coordinates },
    { name: 'Córdoba', coordinates: [-64.1888, -31.4201] as Coordinates },
    { name: 'Mendoza', coordinates: [-68.8458, -32.8895] as Coordinates },
    { name: 'Neuquén', coordinates: [-68.0591, -38.9516] as Coordinates }
  ];

  const projection = geoTransverseMercator()
    .rotate([64, 0])
    .center([0, -38])
    .fitExtent([[56, 28], [464, 526]], argentina);
  const pathGenerator = geoPath(projection);
  const provincePaths = argentina.features.map((feature, index) => ({
    id: feature.properties?.name ?? `province-${index}`,
    path: pathGenerator(feature) ?? ''
  }));
  const origin = projectPoint(buenosAires);
  const routes = destinations.map((destination, index) => {
    const target = projectPoint(destination.coordinates);
    return {
      ...destination,
      target,
      path: routePath(origin, target, index)
    };
  });

  function normalizeCollection(collection: ProvinceCollection): ProvinceCollection {
    return {
      ...collection,
      features: collection.features.map(normalizeFeature)
    };
  }

  function normalizeFeature(feature: ProvinceFeature): ProvinceFeature {
    const geometry = feature.geometry;
    if (geometry.type === 'Polygon') {
      return {
        ...feature,
        geometry: {
          ...geometry,
          coordinates: normalizePolygon(geometry.coordinates)
        }
      };
    }

    return {
      ...feature,
      geometry: {
        ...geometry,
        coordinates: geometry.coordinates.map(normalizePolygon)
      }
    };
  }

  function normalizePolygon(polygon: Polygon['coordinates']): Polygon['coordinates'] {
    return polygon.map((ring, index) => {
      const shouldReverse = index === 0 ? signedRingArea(ring) > 0 : signedRingArea(ring) < 0;
      return shouldReverse ? [...ring].reverse() : ring;
    });
  }

  function signedRingArea(ring: Position[]) {
    let area = 0;
    for (let index = 0, previous = ring.length - 1; index < ring.length; previous = index++) {
      area += ring[previous][0] * ring[index][1] - ring[index][0] * ring[previous][1];
    }
    return area / 2;
  }

  function projectPoint(coordinates: Coordinates) {
    const point = projection(coordinates);
    return {
      x: point?.[0] ?? VIEWBOX_WIDTH / 2,
      y: point?.[1] ?? VIEWBOX_HEIGHT / 2
    };
  }

  function routePath(start: { x: number; y: number }, end: { x: number; y: number }, index: number) {
    const dx = end.x - start.x;
    const dy = end.y - start.y;
    const isSouth = dy > 0;
    const isNorthEast = dy < 0 && dx > 0; // Targets Misiones
    
    // Push midX aggressively left (inland) if going to Misiones
    const midX = start.x + dx * (isNorthEast ? 0.05 : (isSouth ? 0.2 : 0.4));
    const midY = start.y + dy * (isSouth ? 0.6 : 0.3);
    
    return `M ${start.x} ${start.y} L ${midX} ${midY} L ${end.x} ${end.y}`;
  }
</script>

<svg class="scene-map" viewBox={`0 0 ${VIEWBOX_WIDTH} ${VIEWBOX_HEIGHT}`} role="img" aria-labelledby="map-title map-desc">
  <title id="map-title">Alcance Nacional</title>
  <desc id="map-desc">Mapa de Argentina renderizado desde GeoJSON con rutas desde Buenos Aires hacia distintas provincias.</desc>
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
