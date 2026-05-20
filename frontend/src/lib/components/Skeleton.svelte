<script lang="ts">
  export let w = '100%';
  export let h = '20px';
  export let radius = '0';
</script>

<div class="skel" style="width:{w};height:{h};border-radius:{radius}"></div>

<style>
  .skel {
    background: linear-gradient(90deg, var(--color-surface) 25%, rgba(232,224,208,.06) 50%, var(--color-surface) 75%);
    background-size: 200% 100%;
    animation: shimmer 1.6s ease-in-out infinite;
  }
  @keyframes shimmer { from { background-position: 200% 0; } to { background-position: -200% 0; } }
</style>
