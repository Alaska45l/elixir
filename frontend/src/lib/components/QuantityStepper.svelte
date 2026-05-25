<script lang="ts">
  export let value = 1;
  export let min = 1;
  export let max = 99;
  export let onchange: (v: number) => void = () => undefined;

  function decrement() {
    if (value > min) {
      value -= 1;
      onchange(value);
    }
  }

  function increment() {
    if (value < max) {
      value += 1;
      onchange(value);
    }
  }
</script>

<div class="stepper" role="group" aria-label="Cantidad">
  <button
    type="button"
    class="step-btn"
    disabled={value <= min}
    on:click={decrement}
    aria-label="Reducir cantidad"
  >-</button>
  <div class="step-val" aria-live="polite" aria-atomic="true">
    {#key value}
      <span class="num">{value}</span>
    {/key}
  </div>
  <button
    type="button"
    class="step-btn"
    disabled={value >= max}
    on:click={increment}
    aria-label="Aumentar cantidad"
  >+</button>
</div>

<style>
  .stepper {
    display: flex;
    height: 46px;
    width: 120px;
    background: rgba(45,42,36,0.04);
    border-radius: 6px;
    overflow: hidden;
  }

  .step-btn {
    flex: 0 0 38px;
    background: transparent;
    border: 0;
    color: var(--color-text-muted);
    font-size: 1.3rem;
    display: grid;
    place-items: center;
    cursor: pointer;
    transition: color 0.15s ease, background 0.15s ease;
  }

  .step-btn:hover:not(:disabled) {
    color: var(--color-emerald-dark);
    background: var(--color-surface-hover);
  }

  .step-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .step-val {
    flex: 1;
    display: grid;
    place-items: center;
    background: rgba(45,42,36,0.06);
    overflow: hidden;
    position: relative;
  }

  .num {
    display: block;
    font-size: 1rem;
    font-variant-numeric: tabular-nums;
  }

  @media (prefers-reduced-motion: no-preference) {
    .num {
      animation: numIn 0.18s ease both;
    }

    @keyframes numIn {
      from { opacity: 0; transform: translateY(6px); }
      to { opacity: 1; transform: none; }
    }
  }
</style>
