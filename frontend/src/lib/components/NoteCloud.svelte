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
  section div span { display: inline-flex; align-items: center; gap: 6px; border-radius: 999px; background: rgba(45,42,36,0.04); color: var(--color-text); padding: 5px 12px; font-size: .78rem; letter-spacing: .04em; transition: background .2s ease, color .2s ease; }
  section div span::before { content: ''; display: inline-block; width: 4px; height: 4px; border-radius: 50%; background: var(--color-emerald); opacity: .6; flex-shrink: 0; }
  section div span:hover { background: rgba(107,191,138,0.12); color: var(--color-emerald-dark); }
</style>
