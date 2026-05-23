<script lang="ts">
  import { apiFetch, defaultHomepage } from '$lib/api/client';
  import type { HomepageSettings } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import { onMount } from 'svelte';

  let form: HomepageSettings = { ...defaultHomepage };
  let error = '';
  let saving = false;

  onMount(async () => {
    try {
      const res = await apiFetch<{items: HomepageSettings[]}>('/api/admin/homepage');
      form = res.items[0] ?? { ...defaultHomepage };
    } catch {
      form = { ...defaultHomepage };
    }
  });

  async function save() {
    error = '';
    saving = true;
    try {
      await apiFetch('/api/admin/homepage', { method: 'PUT', body: JSON.stringify(form) });
      toast.push('Homepage guardada');
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo guardar la homepage';
    } finally {
      saving = false;
    }
  }
</script>

<div class="head">
  <h1 class="display section-title">Homepage</h1>
  <a class="btn" href="/" target="_blank" rel="noreferrer">Ver homepage</a>
</div>

<form class="form" on:submit|preventDefault={save}>
  <section>
    <h2>Hero principal</h2>
    <label class="field"><span>Título principal</span><textarea class="textarea" bind:value={form.hero_heading}></textarea></label>
    <label class="field"><span>Subtítulo SEO / redes</span><textarea class="textarea" bind:value={form.hero_subheading}></textarea></label>
    <div class="grid-2">
      <label class="field"><span>Texto del botón</span><input class="input" bind:value={form.hero_cta_label} /></label>
      <label class="field"><span>URL del botón</span><input class="input" bind:value={form.hero_cta_url} /></label>
    </div>
    <label class="field"><span>Imagen del hero</span><input class="input" bind:value={form.hero_image_url} placeholder="https://..." /><small>Usá una imagen pública HTTPS ya subida.</small></label>
    {#if form.hero_image_url}<img class="preview wide" src={form.hero_image_url} alt="Vista previa del hero" />{/if}
  </section>

  <section>
    <h2>Bloque editorial</h2>
    <label class="field"><span>Título editorial</span><input class="input" bind:value={form.editorial_heading} /></label>
    <label class="field"><span>Texto editorial</span><textarea class="textarea" bind:value={form.editorial_body}></textarea></label>
    <label class="field"><span>Imagen editorial</span><input class="input" bind:value={form.editorial_image_url} placeholder="https://..." /></label>
    {#if form.editorial_image_url}<img class="preview" src={form.editorial_image_url} alt="Vista previa editorial" />{/if}
  </section>

  {#if error}<p class="error">{error}</p>{/if}
  <button class="btn primary" disabled={saving}>{saving ? 'Guardando...' : 'Guardar homepage'}</button>
</form>

<style>
  .head { display: flex; justify-content: space-between; gap: 20px; align-items: center; }
  .form { display: grid; gap: 28px; max-width: 900px; }
  section { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 16px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  small { color: var(--color-text-muted); }
  .error { margin: 0; color: var(--color-danger-soft); }
  .preview { width: 320px; aspect-ratio: 4/3; object-fit: cover; border: 1px solid var(--color-border); }
  .preview.wide { width: min(100%, 720px); aspect-ratio: 16/7; }
</style>
