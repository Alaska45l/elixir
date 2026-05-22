<script lang="ts">
  import { formatARS } from '$lib/utils/currency';
  export let params: URLSearchParams;
  export let visible = true;

  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
  const concentrations = ['Extrait de Parfum', 'EDP', 'EDT', 'EDC'];
  let min = Number(params.get('min_price') ?? 0);
  let max = Number(params.get('max_price') ?? 20000000);

  $: selectedFamilies = params.getAll('family');
  $: selectedGenders = params.getAll('gender');
  $: selectedConcentrations = params.getAll('concentration');

  function updateMin(event: Event) {
    min = Number((event.currentTarget as HTMLInputElement).value);
    if (min > max) max = min;
  }
  function updateMax(event: Event) {
    max = Number((event.currentTarget as HTMLInputElement).value);
    if (max < min) min = max;
  }
</script>

{#if visible}
  <aside class="filters">
    <p class="eyebrow">Filtros</p>
    <form method="GET">
      <fieldset>
        <legend>Familia</legend>
        {#each families as family}
          <label class="check"><input name="family" value={family} type="checkbox" checked={selectedFamilies.includes(family)} /> {family}</label>
        {/each}
      </fieldset>
      <fieldset>
        <legend>Género</legend>
        {#each genders as gender}
          <label class="check"><input name="gender" value={gender} type="checkbox" checked={selectedGenders.includes(gender)} /> {gender}</label>
        {/each}
      </fieldset>
      <fieldset>
        <legend>Concentración</legend>
        <div class="chips">
          {#each concentrations as concentration}
            <label class:active={selectedConcentrations.includes(concentration)}>
              <input name="concentration" value={concentration} type="checkbox" checked={selectedConcentrations.includes(concentration)} />
              {concentration}
            </label>
          {/each}
        </div>
      </fieldset>
      <div class="price-range field">
        <span>Precio</span>
        <div class="range-current"><b>{formatARS(min)}</b><b>{formatARS(max)}</b></div>
        <input class="range-input" type="range" name="min_price" min="0" max="20000000" step="500000" value={min} on:input={updateMin} aria-label="Precio mínimo" />
        <input class="range-input" type="range" name="max_price" min="0" max="20000000" step="500000" value={max} on:input={updateMax} aria-label="Precio máximo" />
      </div>
      <label class="check"><input name="in_stock" value="true" type="checkbox" checked={params.get('in_stock') === 'true'} /> En stock</label>
      <div class="actions">
        <button class="btn primary" type="submit">Aplicar</button>
        <a class="btn text" href="/fragrances">Limpiar</a>
      </div>
    </form>
  </aside>
{/if}

<style>
  .filters { position: sticky; top: 96px; }
  form { display: grid; gap: 20px; }
  fieldset { border: 0; padding: 0; margin: 0; display: grid; gap: 10px; }
  legend, .field span { color: var(--color-text-muted); font-size: .85rem; margin-bottom: 2px; }
  .check { color: var(--color-text-muted); display: flex; gap: 10px; align-items: center; min-height: 26px; }
  input[type='checkbox'] { accent-color: var(--color-gold); }
  .chips { display: flex; flex-wrap: wrap; gap: 8px; }
  .chips label { border: 1px solid var(--color-border); color: var(--color-text-muted); padding: 8px 10px; font-size: .78rem; cursor: pointer; }
  .chips label.active, .chips label:hover { border-color: var(--color-gold); color: var(--color-gold); }
  .chips input { position: absolute; opacity: 0; pointer-events: none; }
  .range-input { width: 100%; accent-color: var(--color-gold); }
  .range-current { display: flex; justify-content: space-between; gap: 12px; color: var(--color-gold); font-size: .82rem; font-variant-numeric: tabular-nums; }
  .actions { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
  @media (max-width: 860px) { .filters { position: static; } }
</style>
