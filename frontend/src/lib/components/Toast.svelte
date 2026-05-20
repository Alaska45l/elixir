<script lang="ts">
  import { toast } from '$lib/stores/toast';
</script>

<div class="toasts">
  {#each $toast as item (item.id)}
    <div class:error={item.tone === 'error'}>{item.text}</div>
  {/each}
</div>

<style>
  .toasts { position: fixed; right: 18px; bottom: 18px; display: grid; gap: 10px; z-index: 60; }
  .toasts div { background: var(--color-surface); border: 1px solid var(--color-gold); color: var(--color-text); padding: 12px 14px; min-width: 240px; }
  .toasts div.error { border-color: #b66; }
</style>
