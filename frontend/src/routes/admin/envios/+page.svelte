<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { AdminShippingZone } from '$lib/types/admin';
  import { formatARS } from '$lib/utils/currency';

  type ZoneForm = AdminShippingZone & { province_text: string };

  const blank = (): ZoneForm => ({
    id: '',
    zone_name: '',
    province_codes: ['AR'],
    province_text: 'AR',
    base_cost_cents: 0,
    per_kg_cents: 0,
    estimated_days_min: 1,
    estimated_days_max: 5,
    active: true
  });

  let zones: ZoneForm[] = [];
  let form: ZoneForm = blank();
  let loading = true;
  let saving = false;
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      const data = await apiFetch<{ items: AdminShippingZone[] }>('/api/admin/shipping/zones');
      zones = data.items.map(toForm);
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudieron cargar zonas';
    } finally {
      loading = false;
    }
  }

  function toForm(zone: AdminShippingZone): ZoneForm {
    return { ...zone, province_text: (zone.province_codes ?? []).join(', ') };
  }

  function edit(zone: ZoneForm) {
    form = structuredClone(zone);
  }

  function cancelEdit() {
    form = blank();
  }

  function payload(zone: ZoneForm) {
    return {
      zone_name: zone.zone_name,
      province_codes: zone.province_text.split(',').map((item) => item.trim()).filter(Boolean),
      base_cost_cents: Number(zone.base_cost_cents),
      per_kg_cents: Number(zone.per_kg_cents),
      estimated_days_min: Number(zone.estimated_days_min),
      estimated_days_max: Number(zone.estimated_days_max),
      active: zone.active
    };
  }

  async function save() {
    error = '';
    saving = true;
    try {
      const target = form.id ? `/api/admin/shipping/zones/${form.id}` : '/api/admin/shipping/zones';
      await apiFetch(target, { method: form.id ? 'PUT' : 'POST', body: JSON.stringify(payload(form)) });
      toast.push('Zona de envío guardada');
      form = blank();
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo guardar la zona';
    } finally {
      saving = false;
    }
  }

  async function deactivate(zone: ZoneForm) {
    if (!confirm(`¿Desactivar la zona ${zone.zone_name}?`)) return;
    try {
      await apiFetch(`/api/admin/shipping/zones/${zone.id}`, { method: 'DELETE' });
      toast.push('Zona desactivada');
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo desactivar la zona';
    }
  }

  onMount(load);
</script>

<div class="head">
  <h1 class="display section-title">Envíos</h1>
</div>

<div class="layout">
  <form class="panel" on:submit|preventDefault={save}>
    <h2>{form.id ? 'Editar zona' : 'Nueva zona'}</h2>
    <label class="field"><span>Nombre</span><input class="input" required bind:value={form.zone_name} placeholder="CABA, GBA, Interior" /></label>
    <label class="field"><span>Provincias</span><input class="input" bind:value={form.province_text} placeholder="CF, BA o AR para todo el país" /><small>Separá códigos por coma. Usá AR para cubrir provincias no específicas.</small></label>
    <div class="grid-2 tight">
      <label class="field"><span>Costo base</span><input class="input" type="number" min="0" value={Math.round(form.base_cost_cents / 100)} on:input={(event) => form.base_cost_cents = Number((event.currentTarget as HTMLInputElement).value) * 100} /></label>
      <label class="field"><span>Extra por kg</span><input class="input" type="number" min="0" value={Math.round(form.per_kg_cents / 100)} on:input={(event) => form.per_kg_cents = Number((event.currentTarget as HTMLInputElement).value) * 100} /></label>
    </div>
    <div class="grid-2 tight">
      <label class="field"><span>Días mínimos</span><input class="input" type="number" min="0" bind:value={form.estimated_days_min} /></label>
      <label class="field"><span>Días máximos</span><input class="input" type="number" min="0" bind:value={form.estimated_days_max} /></label>
    </div>
    <label class="check"><input type="checkbox" bind:checked={form.active} /> Zona activa</label>
    {#if error}<p class="error">{error}</p>{/if}
    <div class="actions">
      <button class="btn primary" type="submit" disabled={saving}>{saving ? 'Guardando...' : 'Guardar zona'}</button>
      {#if form.id}<button class="btn" type="button" on:click={cancelEdit}>Cancelar</button>{/if}
    </div>
  </form>

  <section class="panel">
    <h2>Zonas configuradas</h2>
    {#if loading}
      <p class="empty">Cargando zonas...</p>
    {:else if zones.length === 0}
      <p class="empty">No hay zonas de envío. Creá una zona para que el checkout pueda cotizar.</p>
    {:else}
      <div class="zones">
        {#each zones as zone}
          <article class:inactive={!zone.active}>
            <div>
              <strong>{zone.zone_name}</strong>
              <span>{zone.province_codes.join(', ') || 'AR'} · {zone.estimated_days_min} a {zone.estimated_days_max} días</span>
            </div>
            <span>{formatARS(zone.base_cost_cents)}{zone.per_kg_cents ? ` + ${formatARS(zone.per_kg_cents)}/kg` : ''}</span>
            <span class:ok={zone.active}>{zone.active ? 'Activa' : 'Inactiva'}</span>
            <button class="btn" type="button" on:click={() => edit(zone)}>Editar</button>
            {#if zone.active}<button class="btn danger" type="button" on:click={() => deactivate(zone)}>Desactivar</button>{/if}
          </article>
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .layout { display: grid; grid-template-columns: minmax(300px, 420px) 1fr; gap: 32px; align-items: start; }
  .panel { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 16px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .tight { gap: 12px; }
  small, .empty { color: var(--color-text-muted); }
  .check { display: flex; gap: 10px; color: var(--color-text-muted); }
  .actions { display: flex; gap: 12px; flex-wrap: wrap; }
  .error { margin: 0; color: var(--color-danger-soft); }
  .zones { display: grid; gap: 12px; }
  article { display: grid; grid-template-columns: 1fr auto 90px auto auto; gap: 12px; align-items: center; border: 1px solid var(--color-border); background: var(--color-surface); padding: 14px; }
  article.inactive { opacity: .65; }
  article div { display: grid; gap: 5px; }
  article span { color: var(--color-text-muted); }
  .ok { color: var(--color-gold); }
  .danger { border-color: #9f5d55; color: #e0a39a; }
  @media (max-width: 980px) { .layout, article { grid-template-columns: 1fr; } }
</style>
