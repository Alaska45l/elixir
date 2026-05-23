<script lang="ts">
  import { apiFetch } from '$lib/api/client';
  import { onMount } from 'svelte';
  import { toast } from '$lib/stores/toast';
  import { formatARS } from '$lib/utils/currency';

  type Discount = {
    id: string;
    code: string;
    discount_type: 'percent' | 'fixed';
    discount_value: number;
    min_order_cents: number;
    max_uses: number | null;
    uses: number;
    active: boolean;
    expires_at: string | null;
  };

  type DiscountForm = {
    id: string;
    code: string;
    discount_type: 'percent' | 'fixed';
    discount_value: number;
    min_order_cents: number;
    max_uses: number | null;
    active: boolean;
    expires_at: string;
  };

  const blank = (): DiscountForm => ({
    id: '',
    code: '',
    discount_type: 'percent',
    discount_value: 10,
    min_order_cents: 0,
    max_uses: null,
    active: true,
    expires_at: ''
  });

  let rows: Discount[] = [];
  let form = blank();
  let loading = true;
  let saving = false;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      rows = (await apiFetch<{ items: Discount[] }>('/api/admin/discounts')).items;
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudieron cargar descuentos';
    } finally {
      loading = false;
    }
  }

  function edit(row: Discount) {
    form = {
      id: row.id,
      code: row.code,
      discount_type: row.discount_type,
      discount_value: row.discount_value,
      min_order_cents: row.min_order_cents,
      max_uses: row.max_uses,
      active: row.active,
      expires_at: row.expires_at ? new Date(row.expires_at).toISOString().slice(0, 16) : ''
    };
  }

  function payload() {
    return {
      code: form.code,
      discount_type: form.discount_type,
      discount_value: Number(form.discount_value),
      min_order_cents: Number(form.min_order_cents),
      max_uses: form.max_uses && form.max_uses > 0 ? Number(form.max_uses) : null,
      active: form.active,
      expires_at: form.expires_at ? new Date(form.expires_at).toISOString() : null
    };
  }

  async function save() {
    saving = true;
    error = '';
    try {
      await apiFetch(form.id ? `/api/admin/discounts/${form.id}` : '/api/admin/discounts', {
        method: form.id ? 'PUT' : 'POST',
        body: JSON.stringify(payload())
      });
      toast.push(form.id ? 'Descuento actualizado' : 'Descuento creado');
      form = blank();
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo guardar el descuento';
    } finally {
      saving = false;
    }
  }

  async function remove(row: Discount) {
    if (!confirm(`¿Eliminar el cupón ${row.code}?`)) return;
    try {
      await apiFetch(`/api/admin/discounts/${row.id}`, { method: 'DELETE' });
      toast.push('Descuento eliminado');
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo eliminar el descuento';
    }
  }

  function valueLabel(row: Discount) {
    return row.discount_type === 'percent' ? `${row.discount_value}%` : formatARS(row.discount_value);
  }

  onMount(load);
</script>

<h1 class="display section-title">Descuentos</h1>

<form class="panel" on:submit|preventDefault={save}>
  <h2>{form.id ? 'Editar cupón' : 'Nuevo cupón'}</h2>
  <div class="grid">
    <label class="field"><span>Código</span><input class="input" required bind:value={form.code} placeholder="BIENVENIDO10" /></label>
    <label class="field"><span>Tipo</span><select class="select" bind:value={form.discount_type}><option value="percent">Porcentaje</option><option value="fixed">Monto fijo</option></select></label>
    <label class="field"><span>{form.discount_type === 'percent' ? 'Porcentaje' : 'Monto en pesos'}</span><input class="input" type="number" min="1" value={form.discount_type === 'fixed' ? Math.round(form.discount_value / 100) : form.discount_value} on:input={(event) => form.discount_value = Number((event.currentTarget as HTMLInputElement).value) * (form.discount_type === 'fixed' ? 100 : 1)} /></label>
    <label class="field"><span>Mínimo de compra</span><input class="input" type="number" min="0" value={Math.round(form.min_order_cents / 100)} on:input={(event) => form.min_order_cents = Number((event.currentTarget as HTMLInputElement).value) * 100} /></label>
    <label class="field"><span>Máximo de usos</span><input class="input" type="number" min="0" bind:value={form.max_uses} placeholder="Sin límite" /></label>
    <label class="field"><span>Vence</span><input class="input" type="datetime-local" bind:value={form.expires_at} /></label>
  </div>
  <label class="check"><input type="checkbox" bind:checked={form.active} /> Cupón activo</label>
  {#if error}<p class="error">{error}</p>{/if}
  <div class="actions">
    <button class="btn primary" type="submit" disabled={saving}>{saving ? 'Guardando...' : 'Guardar descuento'}</button>
    {#if form.id}<button class="btn" type="button" on:click={() => form = blank()}>Cancelar</button>{/if}
  </div>
</form>

{#if loading}
  <p class="empty">Cargando descuentos...</p>
{:else if rows.length === 0}
  <p class="empty">No hay descuentos cargados. Creá un cupón para usarlo en el checkout.</p>
{:else}
  <table class="table">
    <thead><tr><th>Código</th><th>Descuento</th><th>Mínimo</th><th>Usos</th><th>Estado</th><th></th></tr></thead>
    <tbody>
      {#each rows as row}
        <tr>
          <td><strong>{row.code}</strong><br /><span>{row.expires_at ? `Vence ${new Date(row.expires_at).toLocaleString('es-AR')}` : 'Sin vencimiento'}</span></td>
          <td>{valueLabel(row)}</td>
          <td>{formatARS(row.min_order_cents)}</td>
          <td>{row.uses}{row.max_uses ? ` / ${row.max_uses}` : ''}</td>
          <td>{row.active ? 'Activo' : 'Inactivo'}</td>
          <td class="row-actions"><button class="btn" type="button" on:click={() => edit(row)}>Editar</button><button class="btn danger" type="button" on:click={() => remove(row)}>Eliminar</button></td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .panel { border-top: 1px solid var(--color-border); padding-top: 22px; margin-bottom: 28px; display: grid; gap: 16px; max-width: 1040px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
  .check { display: flex; gap: 10px; color: var(--color-text-muted); }
  .actions, .row-actions { display: flex; gap: 10px; flex-wrap: wrap; }
  .error { margin: 0; color: var(--color-danger-soft); }
  .empty, td span { color: var(--color-text-muted); }
  .danger { border-color: #9f5d55; color: #e0a39a; }
  @media (max-width: 900px) { .grid { grid-template-columns: 1fr; } }
</style>
