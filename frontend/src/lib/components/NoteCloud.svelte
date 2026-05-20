<script lang="ts">
  export let top: string[] = [];
  export let heart: string[] = [];
  export let base: string[] = [];
  const rows = [
    ['Notas de salida', top],
    ['Notas de corazón', heart],
    ['Notas de fondo', base]
  ] as const;
</script>

<div class="notes">
  {#each rows as row}
    <section>
      <h3>{row[0]}</h3>
      <div>{#each row[1] as note}<span>{note}</span>{/each}</div>
    </section>
  {/each}
</div>

<style>
  .notes { display: grid; gap: 18px; }
  h3 { margin: 0 0 10px; color: var(--color-text-muted); font-size: .85rem; text-transform: uppercase; letter-spacing: .12em; }
  section div { display: flex; flex-wrap: wrap; gap: 8px; }
  span { border: 1px solid var(--color-border); color: var(--color-text); padding: 7px 10px; }
</style>
