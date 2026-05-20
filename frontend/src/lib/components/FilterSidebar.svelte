<script lang="ts">
  export let params: URLSearchParams;
  const families = ['Oriental', 'Floral', 'Amaderado', 'Cítrico', 'Fresco', 'Gourmand'];
  const genders = ['Unisex', 'Masculino', 'Femenino'];
</script>

<aside class="filters">
  <p class="eyebrow">Filtros</p>
  <form method="GET">
    <label class="field"><span>Familia</span><select class="select" name="family"><option value="">Todas</option>{#each families as f}<option selected={params.get('family') === f}>{f}</option>{/each}</select></label>
    <label class="field"><span>Género</span><select class="select" name="gender"><option value="">Todos</option>{#each genders as g}<option selected={params.get('gender') === g}>{g}</option>{/each}</select></label>
    <label class="field"><span>Concentración</span><select class="select" name="concentration"><option value="">Todas</option><option>Extrait de Parfum</option><option>EDP</option><option>EDT</option><option>EDC</option></select></label>
    <label class="check"><input name="in_stock" value="true" type="checkbox" checked={params.get('in_stock') === 'true'} /> En stock</label>
    <button class="btn primary" type="submit">Aplicar</button>
  </form>
</aside>

<style>
  .filters { position: sticky; top: 96px; }
  form { display: grid; gap: 18px; }
  .check { color: var(--color-text-muted); display: flex; gap: 10px; align-items: center; }
  @media (max-width: 860px) { .filters { position: static; } }
</style>
