<script lang="ts">
  import { wishlist } from '$lib/stores/wishlist';
  export let slug: string;
</script>

<button class="wishlist-btn" class:active={$wishlist.includes(slug)} type="button" on:click={() => wishlist.toggle(slug)} aria-label="Mi lista de deseos">♡</button>

<style>
  button {
    width: 38px;
    height: 38px;
    border: 0;
    border-radius: 6px;
    color: var(--color-text);
    background: rgba(13,31,21,.68);
    backdrop-filter: blur(10px);
    font-size: 1.25rem;
  }

  button.active {
    color: var(--color-gold);
    background: rgba(13,31,21,.82);
  }
</style>
