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
    border-radius: 50%;
    color: var(--color-on-image);
    background: rgba(255,255,255,0.20);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    font-size: 1.25rem;
    text-shadow: 0 1px 4px rgba(0,0,0,0.35);
    transition: color 0.2s ease, background 0.2s ease, transform 0.2s ease;
  }

  button:hover {
    background: rgba(255,255,255,0.35);
    transform: scale(1.08);
  }

  button.active {
    color: var(--color-emerald);
    background: rgba(255,255,255,0.30);
    text-shadow: 0 1px 4px rgba(0,0,0,0.25);
  }
</style>
