<script lang="ts">
  import BentoFilterBar from './BentoFilterBar.svelte';

  export let params: URLSearchParams;
  export let visible = true;
</script>

{#if visible}
  <BentoFilterBar {params} />
{/if}
