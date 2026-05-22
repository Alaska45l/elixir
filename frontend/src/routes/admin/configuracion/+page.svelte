<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch, defaultSettings } from '$lib/api/client';
  import type { SiteSettings } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';

  let form: SiteSettings = structuredClone(defaultSettings);

  onMount(async () => {
    try {
      form = await apiFetch<SiteSettings>('/api/admin/settings');
    } catch {
      form = structuredClone(defaultSettings);
    }
  });

  function addFAQ() {
    form.faq_items = [...form.faq_items, { question: '', answer: '' }];
  }

  function addCategory() {
    form.navbar_product_categories = [...form.navbar_product_categories, { label: '', href: '/fragrances' }];
  }

  async function save() {
    await apiFetch('/api/admin/settings', { method: 'PUT', body: JSON.stringify(form) });
    toast.push('Configuración guardada');
  }
</script>

<h1 class="display section-title">Configuración</h1>

<form class="form" on:submit|preventDefault={save}>
  <section>
    <h2>Barra de anuncio</h2>
    <label class="check"><input type="checkbox" bind:checked={form.announcement_bar_active} /> Mostrar barra</label>
    <label class="field"><span>Texto</span><input class="input" bind:value={form.announcement_bar_text} /></label>
  </section>

  <section>
    <h2>Footer y redes sociales</h2>
    <div class="grid-2">
      <label class="field"><span>Instagram</span><input class="input" bind:value={form.footer_instagram_url} placeholder="https://instagram.com/..." /></label>
      <label class="field"><span>TikTok</span><input class="input" bind:value={form.footer_tiktok_url} placeholder="https://tiktok.com/@..." /></label>
      <label class="field"><span>WhatsApp</span><input class="input" bind:value={form.footer_whatsapp_url} placeholder="https://wa.me/..." /></label>
    </div>
  </section>

  <section>
    <h2>Nosotros</h2>
    <label class="field"><span>Título</span><input class="input" bind:value={form.about_title} /></label>
    <label class="field"><span>Descripción</span><textarea class="textarea" bind:value={form.about_description}></textarea></label>
    <div class="grid-2">
      <label class="field"><span>Ubicación</span><input class="input" bind:value={form.about_location} /></label>
      <label class="field"><span>Teléfono</span><input class="input" bind:value={form.about_phone} /></label>
    </div>
  </section>

  <section>
    <div class="row-head"><h2>Preguntas frecuentes</h2><button class="btn" type="button" on:click={addFAQ}>Agregar</button></div>
    {#each form.faq_items as item, i}
      <div class="repeat-row">
        <label class="field"><span>Pregunta</span><input class="input" bind:value={item.question} /></label>
        <label class="field"><span>Respuesta</span><textarea class="textarea" bind:value={item.answer}></textarea></label>
        <button type="button" on:click={() => form.faq_items = form.faq_items.filter((_, index) => index !== i)}>Quitar</button>
      </div>
    {/each}
  </section>

  <section>
    <h2>Política de devolución</h2>
    <label class="field"><span>Contenido HTML</span><textarea class="textarea tall" bind:value={form.return_policy_html}></textarea></label>
  </section>

  <section>
    <div class="row-head"><h2>Categorías de navegación</h2><button class="btn" type="button" on:click={addCategory}>Agregar</button></div>
    {#each form.navbar_product_categories as item, i}
      <div class="nav-row">
        <label class="field"><span>Etiqueta</span><input class="input" bind:value={item.label} /></label>
        <label class="field"><span>URL</span><input class="input" bind:value={item.href} /></label>
        <button type="button" on:click={() => form.navbar_product_categories = form.navbar_product_categories.filter((_, index) => index !== i)}>Quitar</button>
      </div>
    {/each}
  </section>

  <button class="btn primary" type="submit">Guardar configuración</button>
</form>

<style>
  .form { display: grid; gap: 30px; max-width: 980px; }
  section { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 16px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .check { display: flex; gap: 10px; color: var(--color-text-muted); }
  .row-head { display: flex; justify-content: space-between; align-items: center; gap: 16px; }
  .repeat-row, .nav-row { display: grid; grid-template-columns: 1fr 1.4fr auto; gap: 12px; align-items: end; }
  .repeat-row button, .nav-row button { border: 1px solid var(--color-border); background: transparent; color: var(--color-text); min-height: 46px; padding: 0 14px; }
  .tall { min-height: 220px; }
  @media (max-width: 820px) { .repeat-row, .nav-row { grid-template-columns: 1fr; } }
</style>
