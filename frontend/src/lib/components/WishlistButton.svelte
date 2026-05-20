<script lang="ts">
  import { wishlist } from '$lib/stores/wishlist';
  export let slug: string;
</script>

<button class:active={$wishlist.includes(slug)} type="button" on:click={() => wishlist.toggle(slug)} aria-label="Mi lista de deseos">♡</button>

<style>
  button { width: 38px; height: 38px; border: 1px solid var(--color-border); color: var(--color-text); background: rgba(13,31,21,.62); font-size: 1.25rem; }
  button.active { color: var(--color-gold); border-color: var(--color-gold); }
</style>
