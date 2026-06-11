<script lang="ts">
  export let params: URLSearchParams;

  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
  const priceSorts = [
    { label: 'Mayor a menor', value: 'price_desc' },
    { label: 'Menor a mayor', value: 'price_asc' }
  ];

  $: selectedFamilies = params.getAll('family');
  $: selectedGenders = params.getAll('gender');
  $: selectedSort = params.get('sort') ?? '';

  function submitOnChange(event: Event) {
    (event.currentTarget as HTMLFormElement).requestSubmit();
  }
</script>

<form class="filter-panel" method="GET" on:change={submitOnChange}>
  <div class="filter-grid">
    <div class="filter-box" role="group" aria-labelledby="desktop-family-filter">
      <span id="desktop-family-filter" class="filter-title">Familia</span>
      <div class="chips">
        {#each families as family}
          <label class="chip" class:active={selectedFamilies.includes(family)}>
            <input name="family" value={family} type="checkbox" checked={selectedFamilies.includes(family)} />
            {family}
          </label>
        {/each}
      </div>
    </div>

    <div class="filter-box" role="group" aria-labelledby="desktop-gender-filter">
      <span id="desktop-gender-filter" class="filter-title">Género</span>
      <div class="chips">
        {#each genders as gender}
          <label class="chip" class:active={selectedGenders.includes(gender)}>
            <input name="gender" value={gender} type="checkbox" checked={selectedGenders.includes(gender)} />
            {gender}
          </label>
        {/each}
      </div>
    </div>

    <div class="filter-box sort-box" role="group" aria-labelledby="desktop-sort-filter">
      <span id="desktop-sort-filter" class="filter-title">Ordenar por</span>
      <span class="filter-subtitle">Precio</span>
      <div class="chips">
        {#each priceSorts as option}
          <label class="chip" class:active={selectedSort === option.value}>
            <input name="sort" value={option.value} type="radio" checked={selectedSort === option.value} />
            {option.label}
          </label>
        {/each}
      </div>
    </div>

  </div>

  <div class="filter-actions">
    <button class="btn primary" type="submit">Aplicar</button>
    <a class="btn text" href="/fragrances">Limpiar</a>
  </div>
</form>

<style>
  .filter-panel {
    display: grid;
    gap: 16px;
    padding: 20px 0 0;
    margin-bottom: 32px;
  }

  .filter-grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    align-items: stretch;
  }

  .filter-box {
    min-width: 0;
    border: 0;
    border-radius: 14px;
    background: rgba(45,42,36,0.04);
    padding: 16px;
    display: grid;
    align-content: start;
    gap: 12px;
    margin: 0;
  }

  .sort-box {
    grid-column: span 1;
  }

  .sort-box .chips {
    width: 100%;
  }

  .sort-box .chip {
    flex: 1 1 130px;
    min-width: 0;
    white-space: normal;
    text-align: center;
    line-height: 1.2;
  }

  .filter-title {
    color: var(--color-text-muted);
    font-size: .78rem;
    text-transform: uppercase;
    letter-spacing: .12em;
    margin: 0;
    display: block;
  }

  .filter-subtitle {
    color: var(--color-text);
    font-size: .9rem;
    font-weight: 700;
  }

  .chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .chip {
    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-height: 42px;
    border-radius: 6px;
    padding: 8px 16px;
    border: 1px solid var(--color-border);
    color: var(--color-text-muted);
    background: transparent;
    font-size: .82rem;
    line-height: 1;
    cursor: pointer;
    transition: all .2s ease;
    white-space: nowrap;
  }

  .chip:hover,
  .chip.active {
    background: var(--color-emerald);
    color: #FDF8F0;
    border-color: var(--color-emerald);
  }

  .chip input {
    position: absolute;
    inset: 0;
    opacity: 0;
    pointer-events: none;
  }

  .filter-actions {
    display: flex;
    gap: 12px;
    align-items: center;
    justify-content: flex-end;
  }

  .filter-actions .btn {
    min-height: 42px;
    padding: 0 16px;
    font-size: .82rem;
  }

  @media (max-width: 1180px) {
    .filter-grid {
      grid-template-columns: repeat(2, minmax(0, 1fr));
    }

    .sort-box {
      grid-column: span 2;
    }
  }
</style>
