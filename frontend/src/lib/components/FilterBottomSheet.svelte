<script lang="ts">
  export let open = false;
  export let params: URLSearchParams;
  export let onClose: () => void = () => undefined;

  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
  const priceSorts = [
    { label: 'Mayor a menor', value: 'price_desc' },
    { label: 'Menor a mayor', value: 'price_asc' }
  ];

  let startY = 0;
  let dragY = 0;
  let dragging = false;

  $: selectedFamilies = params.getAll('family');
  $: selectedGenders = params.getAll('gender');
  $: selectedSort = params.get('sort') ?? '';
  $: inStock = params.get('in_stock') === 'true';

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
          <div class="filter-box" role="group" aria-labelledby="mobile-family-filter">
            <span id="mobile-family-filter" class="filter-title">Familia</span>
            <div class="chips">
              {#each families as family}
                <label class="chip" class:active={selectedFamilies.includes(family)}>
                  <input name="family" value={family} type="checkbox" checked={selectedFamilies.includes(family)} />
                  {family}
                </label>
              {/each}
            </div>
          </div>

          <div class="filter-box" role="group" aria-labelledby="mobile-gender-filter">
            <span id="mobile-gender-filter" class="filter-title">Género</span>
            <div class="chips">
              {#each genders as gender}
                <label class="chip" class:active={selectedGenders.includes(gender)}>
                  <input name="gender" value={gender} type="checkbox" checked={selectedGenders.includes(gender)} />
                  {gender}
                </label>
              {/each}
            </div>
          </div>

          <div class="filter-box sort-box" role="group" aria-labelledby="mobile-sort-filter">
            <span id="mobile-sort-filter" class="filter-title">Ordenar por</span>
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

          <div class="filter-box" role="group" aria-labelledby="mobile-stock-filter">
            <span id="mobile-stock-filter" class="filter-title">Stock</span>
            <label class="chip stock-chip" class:active={inStock}>
              <input name="in_stock" value="true" type="checkbox" checked={inStock} />
              En stock
            </label>
          </div>
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

  .sort-box .chip {
    flex: 1 1 130px;
    min-width: 0;
    white-space: normal;
    text-align: center;
    line-height: 1.2;
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
