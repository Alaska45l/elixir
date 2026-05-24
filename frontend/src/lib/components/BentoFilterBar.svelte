<script lang="ts">
  import { formatARS } from '$lib/utils/currency';

  export let params: URLSearchParams;

  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
  const concentrations = [
    { label: 'Extrait', value: 'Extrait de Parfum' },
    { label: 'EDP', value: 'EDP' },
    { label: 'EDT', value: 'EDT' },
    { label: 'EDC', value: 'EDC' }
  ];

  let min = 0;
  let max = 20000000;
  let selectedFamily = '';
  let selectedGender = '';

  $: selectedConcentrations = params.getAll('concentration');
  $: inStock = params.get('in_stock') === 'true';
  $: if (params) {
    selectedFamily = params.get('family') ?? '';
    selectedGender = params.get('gender') ?? '';
    min = Number(params.get('min_price') ?? 0);
    max = Number(params.get('max_price') ?? 20000000);
  }

  function updateMin(event: Event) {
    min = Number((event.currentTarget as HTMLInputElement).value);
    if (min > max) max = min;
  }

  function updateMax(event: Event) {
    max = Number((event.currentTarget as HTMLInputElement).value);
    if (max < min) min = max;
  }

  function prepareSubmit(event: SubmitEvent) {
    const form = event.currentTarget as HTMLFormElement;
    const selects = Array.from(form.querySelectorAll('select'));
    selects.forEach((select) => {
      if (!select.value) select.disabled = true;
    });
    window.setTimeout(() => selects.forEach((select) => select.disabled = false), 0);
  }
</script>

<form class="filter-row" method="GET" on:submit={prepareSubmit}>
  <label class="filter-field">
    <span>Familia</span>
    <select class="select" name="family" bind:value={selectedFamily}>
      <option value="">Todas</option>
      {#each families as family}<option value={family}>{family}</option>{/each}
    </select>
  </label>

  <label class="filter-field">
    <span>Género</span>
    <select class="select" name="gender" bind:value={selectedGender}>
      <option value="">Todos</option>
      {#each genders as gender}<option value={gender}>{gender}</option>{/each}
    </select>
  </label>

  <div class="filter-field concentration-chips">
    <span>Concentración</span>
    <div class="chips">
      {#each concentrations as concentration}
        <label class="chip" class:active={selectedConcentrations.includes(concentration.value)}>
          <input name="concentration" value={concentration.value} type="checkbox" checked={selectedConcentrations.includes(concentration.value)} />
          {concentration.label}
        </label>
      {/each}
    </div>
  </div>

  <div class="filter-field price-field">
    <span>Precio</span>
    <div class="price-row">
      <span class="price-label">{formatARS(min)}</span>
      <input type="range" name="min_price" min="0" max="20000000" step="500000" value={min} on:input={updateMin} aria-label="Precio mínimo" />
      <input type="range" name="max_price" min="0" max="20000000" step="500000" value={max} on:input={updateMax} aria-label="Precio máximo" />
      <span class="price-label">{formatARS(max)}</span>
    </div>
  </div>

  <label class="chip stock-chip" class:active={inStock}>
    <input name="in_stock" value="true" type="checkbox" checked={inStock} />
    En stock
  </label>

  <div class="filter-actions">
    <button class="btn primary" type="submit">Aplicar</button>
    <a class="btn text" href="/fragrances">Limpiar</a>
  </div>
</form>

<style>
  .filter-row {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-end;
    gap: 20px;
    padding: 20px 0;
    border-bottom: 1px solid var(--color-border);
    margin-bottom: 32px;
  }

  .filter-field {
    display: grid;
    gap: 6px;
    min-width: 150px;
  }

  .filter-field > span {
    color: var(--color-text-muted);
    font-size: .72rem;
    text-transform: uppercase;
    letter-spacing: .12em;
  }

  .filter-field .select {
    min-height: 38px;
    padding-block: 8px;
  }

  .chips {
    display: flex;
    gap: 6px;
  }

  .chip {
    position: relative;
    display: inline-flex;
    align-items: center;
    border-radius: 6px;
    padding: 7px 14px;
    border: 1px solid var(--color-border);
    color: var(--color-text-muted);
    background: transparent;
    font-size: .8rem;
    cursor: pointer;
    transition: all .2s ease;
    white-space: nowrap;
  }

  .chip.active {
    background: var(--color-gold);
    color: var(--color-bg);
    border-color: var(--color-gold);
  }

  .chip input {
    position: absolute;
    inset: 0;
    opacity: 0;
    pointer-events: none;
  }

  .price-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .price-row input[type="range"] {
    width: 80px;
    accent-color: var(--color-gold);
  }

  .price-label {
    color: var(--color-text-muted);
    font-size: .78rem;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }

  .stock-chip {
    min-height: 38px;
    align-self: flex-end;
  }

  .filter-actions {
    display: flex;
    gap: 12px;
    align-items: center;
    align-self: flex-end;
  }

  .filter-actions .btn {
    min-height: 38px;
    padding: 0 16px;
    font-size: .82rem;
  }

  @media (max-width: 1040px) {
    .price-field {
      flex: 1 1 100%;
    }
  }
</style>
