<script lang="ts">
  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>Política de devolución | ELIXIR Exclusive</title>
  <meta name="description" content="Política de cambios y devolución de ELIXIR Exclusive." />
</svelte:head>

<section class="container page-pad policy">
  <p class="eyebrow">Postventa</p>
  <h1 class="display section-title">Política de devolución</h1>
  <div class="gold-rule"></div>
  <div class="copy">{@html data.settings.return_policy_html}</div>
</section>

<style>
  .policy { max-width: 860px; }
  .copy { color: var(--color-text-muted); font-size: 1.05rem; line-height: 1.8; }
  .copy :global(a) { color: var(--color-gold); }
</style>
