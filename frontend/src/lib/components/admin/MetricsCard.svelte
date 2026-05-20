<script lang="ts">
  import { reveal } from '$lib/utils/reveal';
  export let label: string;
  export let value: string;
</script>
<article class="metric" use:reveal>
  <span>{label}</span>
  <strong>{value}</strong>
</article>
<style>
  .metric { border: 1px solid var(--color-border); padding: 22px; background: var(--color-surface); }
  span { color: var(--color-text-muted); font-size: .8rem; text-transform: uppercase; letter-spacing: .12em; }
  strong { display: block; margin-top: 14px; font-size: 1.8rem; color: var(--color-text); }
</style>
