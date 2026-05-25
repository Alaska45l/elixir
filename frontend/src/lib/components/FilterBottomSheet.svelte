<script lang="ts">
  import { formatARS } from '$lib/utils/currency';

  export let open = false;
  export let params: URLSearchParams;
  export let onClose: () => void = () => undefined;

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
  let startY = 0;
  let dragY = 0;
  let dragging = false;

  $: selectedFamilies = params.getAll('family');
  $: selectedGenders = params.getAll('gender');
  $: selectedConcentrations = params.getAll('concentration');
  $: inStock = params.get('in_stock') === 'true';
  $: if (params) {
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

  function beginDrag(event: PointerEvent) {
    dragging = true;
    startY = event.clientY;
    dragY = 0;
    (event.currentTarget as HTMLElement).setPointerCapture(event.pointerId);
  }

  function moveDrag(event: PointerEvent) {
    if (!dragging) return;
    dragY = Math.max(0, event.clientY - startY);
  }

  function endDrag() {
    if (dragY > 80) onClose();
    dragging = false;
    dragY = 0;
  }

  function closeAfterClick() {
    setTimeout(onClose, 0);
  }
</script>

{#if open}
  <div class="sheet-layer" role="presentation">
    <button class="backdrop" type="button" aria-label="Cerrar filtros" on:click={onClose}></button>
    <div
      class="sheet"
      class:dragging
      style={`--drag-y: ${dragY}px`}
      role="dialog"
      aria-modal="true"
      aria-label="Filtros de catálogo"
    >
      <button
        class="drag-handle"
        type="button"
        aria-label="Arrastrar para cerrar"
        on:pointerdown={beginDrag}
        on:pointermove={moveDrag}
        on:pointerup={endDrag}
        on:pointercancel={endDrag}
      ></button>

      <form method="GET">
        <div class="sheet-content">
          <fieldset class="filter-box">
            <legend>Familia</legend>
            <div class="chips">
              {#each families as family}
                <label class="chip" class:active={selectedFamilies.includes(family)}>
                  <input name="family" value={family} type="checkbox" checked={selectedFamilies.includes(family)} />
                  {family}
                </label>
              {/each}
            </div>
          </fieldset>

          <fieldset class="filter-box">
            <legend>Género</legend>
            <div class="chips">
              {#each genders as gender}
                <label class="chip" class:active={selectedGenders.includes(gender)}>
                  <input name="gender" value={gender} type="checkbox" checked={selectedGenders.includes(gender)} />
                  {gender}
                </label>
              {/each}
            </div>
          </fieldset>

          <fieldset class="filter-box">
            <legend>Concentración</legend>
            <div class="chips">
              {#each concentrations as concentration}
                <label class="chip" class:active={selectedConcentrations.includes(concentration.value)}>
                  <input name="concentration" value={concentration.value} type="checkbox" checked={selectedConcentrations.includes(concentration.value)} />
                  {concentration.label}
                </label>
              {/each}
            </div>
          </fieldset>

          <div class="filter-box price-box">
            <span class="legend">Precio</span>
            <div class="range-stack">
              <input class="range-input" type="range" name="min_price" min="0" max="20000000" step="500000" value={min} on:input={updateMin} aria-label="Precio mínimo" />
              <input class="range-input" type="range" name="max_price" min="0" max="20000000" step="500000" value={max} on:input={updateMax} aria-label="Precio máximo" />
            </div>
            <span class="range-current">{formatARS(min)} — {formatARS(max)}</span>
          </div>

          <fieldset class="filter-box">
            <legend>Stock</legend>
            <label class="chip stock-chip" class:active={inStock}>
              <input name="in_stock" value="true" type="checkbox" checked={inStock} />
              En stock
            </label>
          </fieldset>
        </div>

        <div class="bottom-actions">
          <button class="btn primary" type="submit" on:click={closeAfterClick}>Aplicar</button>
          <a class="btn text" href="/fragrances" on:click={closeAfterClick}>Limpiar</a>
        </div>
      </form>
    </div>
  </div>
{/if}

<style>
  .sheet-layer {
    position: fixed;
    inset: 0;
    z-index: 70;
  }

  .backdrop {
    position: absolute;
    inset: 0;
    border: 0;
    padding: 0;
    background: rgba(0,0,0,0.5);
    backdrop-filter: blur(4px);
  }

  .sheet {
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    height: 85%;
    border-radius: 20px 20px 0 0;
    background: var(--color-surface);
    box-shadow: 0 -24px 70px rgba(0,0,0,.42);
    overflow: hidden;
    transform: translateY(var(--drag-y, 0));
    animation: sheet-in .3s cubic-bezier(0.22, 1, 0.36, 1);
  }

  .sheet.dragging {
    animation: none;
    transition: none;
  }

  .drag-handle {
    width: 54px;
    height: 24px;
    padding: 0;
    border: 0;
    background: transparent;
    display: block;
    margin: 10px auto 4px;
    touch-action: none;
  }

  .drag-handle::before {
    content: '';
    display: block;
    width: 42px;
    height: 4px;
    margin: 10px auto;
    border-radius: 999px;
    background: color-mix(in srgb, var(--color-text-muted) 55%, transparent);
  }

  form {
    height: calc(100% - 38px);
    display: grid;
    grid-template-rows: 1fr auto;
  }

  .sheet-content {
    overflow: auto;
    padding: 18px 18px 96px;
    display: grid;
    gap: 12px;
  }

  .filter-box {
    border-radius: 14px;
    background: rgba(45,42,36,0.04);
    padding: 16px;
    display: grid;
    gap: 12px;
    margin: 0;
  }

  fieldset {
    border: 0;
  }

  legend,
  .legend {
    color: var(--color-text-muted);
    font-size: .78rem;
    text-transform: uppercase;
    letter-spacing: .12em;
    margin: 0;
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
    min-height: 48px;
    border-radius: 6px;
    padding: 8px 16px;
    border: 1px solid var(--color-border);
    color: var(--color-text-muted);
    background: transparent;
    font-size: .82rem;
    line-height: 1;
    cursor: pointer;
    transition: all .2s ease;
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

  .stock-chip {
    width: fit-content;
  }

  .range-stack {
    display: grid;
    gap: 8px;
  }

  .range-input {
    width: 100%;
    accent-color: var(--color-emerald);
  }

  .range-current {
    display: block;
    color: var(--color-emerald-dark);
    font-size: .82rem;
    line-height: 1.4;
    font-variant-numeric: tabular-nums;
  }

  .bottom-actions {
    position: sticky;
    bottom: 0;
    display: flex;
    gap: 12px;
    align-items: center;
    justify-content: space-between;
    padding: 14px 18px 18px;
    background: color-mix(in srgb, var(--color-surface) 94%, rgba(0,0,0,.28));
    box-shadow: 0 -14px 32px rgba(0,0,0,0.18);
  }

  .bottom-actions .btn {
    min-height: 48px;
    flex: 1;
  }

  @keyframes sheet-in {
    from { transform: translateY(100%); }
    to { transform: translateY(var(--drag-y, 0)); }
  }
</style>
