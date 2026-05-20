<script lang="ts">
  import { formatARS } from '$lib/utils/currency';
  export let params: URLSearchParams;
  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
  let min = Number(params.get('min_price') ?? 0);
  let max = Number(params.get('max_price') ?? 20000000);
  function updateMin(event: Event) {
    min = Number((event.currentTarget as HTMLInputElement).value);
    if (min > max) max = min;
  }
  function updateMax(event: Event) {
    max = Number((event.currentTarget as HTMLInputElement).value);
    if (max < min) min = max;
  }
</script>

<aside class="filters">
  <p class="eyebrow">Filtros</p>
  <form method="GET">
    <label class="field"><span>Familia</span><select class="select" name="family"><option value="">Todas</option>{#each families as f}<option selected={params.get('family') === f}>{f}</option>{/each}</select></label>
    <label class="field"><span>Género</span><select class="select" name="gender"><option value="">Todos</option>{#each genders as g}<option selected={params.get('gender') === g}>{g}</option>{/each}</select></label>
    <label class="field"><span>Concentración</span><select class="select" name="concentration"><option value="">Todas</option><option>Extrait de Parfum</option><option>EDP</option><option>EDT</option><option>EDC</option></select></label>
    <div class="price-range field">
      <span>Precio</span>
      <div class="range-labels"><span>{formatARS(min)}</span><span>{formatARS(max)}</span></div>
      <p class="range-current">ARS {formatARS(min)} - {formatARS(max)}</p>
      <input class="range-input" type="range" name="min_price" min="0" max="20000000" step="500000" value={min} on:input={updateMin} />
      <input class="range-input" type="range" name="max_price" min="0" max="20000000" step="500000" value={max} on:input={updateMax} />
    </div>
    <label class="check"><input name="in_stock" value="true" type="checkbox" checked={params.get('in_stock') === 'true'} /> En stock</label>
    <button class="btn primary" type="submit">Aplicar</button>
  </form>
</aside>

<style>
  .filters { position: sticky; top: 96px; }
  form { display: grid; gap: 18px; }
  .check { color: var(--color-text-muted); display: flex; gap: 10px; align-items: center; }
  .range-input { width: 100%; accent-color: var(--color-gold); }
  .range-labels { display: flex; justify-content: space-between; color: var(--color-text-muted); font-size: .78rem; }
  .range-current { margin: 0; color: var(--color-gold); font-size: .82rem; font-variant-numeric: tabular-nums; }
  @media (max-width: 860px) { .filters { position: static; } }
</style>
