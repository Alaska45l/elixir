<script lang="ts">
  import type { SiteSettings } from '$lib/api/client';
  export let settings: SiteSettings;
</script>

{#if settings.announcement_bar_active && settings.announcement_bar_text}
  <div class="bar" role="banner">
    <p>{settings.announcement_bar_text}</p>
  </div>
{/if}

<style>
  .bar {
    background: var(--color-gold);
    color: var(--color-bg);
    text-align: center;
    padding: 9px 16px;
    font-size: 0.72rem;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    font-weight: 600;
  }

  .bar p {
    margin: 0;
  }
</style>
