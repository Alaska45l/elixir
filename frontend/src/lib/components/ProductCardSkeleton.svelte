<script lang="ts">
  import Skeleton from './Skeleton.svelte';
</script>

<article class="skel-card">
  <Skeleton h="420px" />
  <Skeleton w="70%" h="28px" />
  <Skeleton w="95%" h="18px" />
  <Skeleton w="42%" h="18px" />
</article>

<style>
  .skel-card { display: grid; gap: 14px; }
</style>
