<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch, uploadAdminImage } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { ProductFormValue } from '$lib/types/product-form';
  import { slugify } from '$lib/utils/slug';

  export let id: string | undefined = undefined;
  export let product: ProductFormValue | undefined = undefined;

  const blank: ProductFormValue = {
    slug: '',
    name: '',
    tagline: '',
    description: '',
    scent_family: 'Oriental',
    gender_tag: 'Unisex',
    concentration: 'EDP',
    top_notes: [],
    heart_notes: [],
    base_notes: [],
    featured: false,
    display_order: 0,
    variants: [{ size_ml: 50, price_ars_cents: 8900000, stock: 10, sku: '', weight_grams: 200 }],
    images: [{ url: '', alt_text: '', is_primary: true, sort_order: 0 }]
  };

  function prepare(value: ProductFormValue) {
    const next = structuredClone(value);
    next.variants = next.variants.map((variant) => ({ ...variant, weight_grams: variant.weight_grams ?? 200 }));
    return next;
  }

  let form: ProductFormValue = prepare(product ?? blank);
  let notes = { top: '', heart: '', base: '' };
  let error = '';
  let saving = false;
  let loading = Boolean(id && !product);
  let uploadingImageIndex: number | null = null;

  onMount(async () => {
    if (!id || product) return;
    try {
      const data = await apiFetch<{ product: ProductFormValue }>(`/api/admin/products/${id}`);
      form = prepare(data.product);
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo cargar el producto';
    } finally {
      loading = false;
    }
  });

  function nameInput() {
    if (!form.slug) form.slug = slugify(form.name);
  }
  function addNote(group: 'top_notes' | 'heart_notes' | 'base_notes', value: string) {
    const note = value.trim().replace(/,$/, '');
    if (note && !form[group].includes(note)) form[group] = [...form[group], note];
  }
  function keyNote(event: KeyboardEvent, group: 'top_notes' | 'heart_notes' | 'base_notes') {
    if (event.key !== 'Enter' && event.key !== ',') return;
    event.preventDefault();
    const input = event.currentTarget as HTMLInputElement;
    addNote(group, input.value);
    input.value = '';
  }
  function removeNote(group: 'top_notes' | 'heart_notes' | 'base_notes', note: string) {
    form[group] = form[group].filter((item) => item !== note);
  }
  function addVariant() {
    form.variants = [...form.variants, { size_ml: 100, price_ars_cents: 0, stock: 0, sku: `${form.slug}-100`, weight_grams: 320 }];
  }
  function addImage() {
    form.images = [...form.images, { url: '', alt_text: form.name, is_primary: false, sort_order: form.images.length }];
  }
  function setPrimary(index: number) {
    form.images = form.images.map((img, i) => ({ ...img, is_primary: i === index }));
  }
  async function uploadImage(event: Event, index: number) {
    const input = event.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    error = '';
    uploadingImageIndex = index;
    try {
      const url = await uploadAdminImage(file, 'products');
      form.images = form.images.map((image, i) => i === index ? { ...image, url, alt_text: image.alt_text || form.name } : image);
      toast.push('Imagen subida como WebP');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'No se pudo subir la imagen';
      error = message;
      toast.push(message, 'error');
    } finally {
      uploadingImageIndex = null;
      input.value = '';
    }
  }
  async function save() {
    error = '';
    if (uploadingImageIndex !== null) {
      error = 'Esperá a que termine la carga de imagen.';
      return;
    }
    saving = true;
    if (!form.name.trim() || !form.slug.trim()) {
      error = 'Completá nombre y slug.';
      saving = false;
      return;
    }
    if (!form.variants.length || form.variants.some((variant) => variant.size_ml <= 0 || variant.price_ars_cents <= 0 || variant.weight_grams <= 0)) {
      error = 'Cada variante necesita tamaño, precio y peso válidos.';
      saving = false;
      return;
    }
    if (!form.images.some((image) => image.url.trim())) {
      error = 'Agregá al menos una imagen pública HTTPS.';
      saving = false;
      return;
    }
    try {
      await apiFetch(id ? `/api/admin/products/${id}` : '/api/admin/products', {
        method: id ? 'PUT' : 'POST',
        body: JSON.stringify(form)
      });
      toast.push('Producto guardado');
      location.href = '/admin/productos';
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo guardar el producto';
    } finally {
      saving = false;
    }
  }
  async function deleteProduct() {
    if (!id || !confirm('¿Eliminar este producto?')) return;
    try {
      await apiFetch(`/api/admin/products/${id}`, { method: 'DELETE' });
      toast.push('Producto eliminado');
      location.href = '/admin/productos';
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo eliminar el producto';
    }
  }
</script>

{#if loading}
  <p class="empty">Cargando producto...</p>
{:else}
<form class="form" on:submit|preventDefault={save}>
  <div class="grid-2">
    <label class="field"><span>Nombre comercial</span><input class="input" required bind:value={form.name} on:input={nameInput} /><small>Nombre visible en catálogo, carrito y orden.</small></label>
    <label class="field"><span>Slug de URL</span><input class="input" required bind:value={form.slug} /><small>Texto corto para la URL pública de la fragancia.</small></label>
  </div>
  <label class="field"><span>Frase corta</span><input class="input" bind:value={form.tagline} /><small>Aparece debajo del nombre en tarjetas y detalle.</small></label>
  <label class="field"><span>Descripción completa</span><textarea class="textarea" bind:value={form.description}></textarea></label>
  <div class="grid-2">
    <label class="field"><span>Familia</span><select class="select" bind:value={form.scent_family}><option>Oriental</option><option>Floral</option><option>Amaderado</option><option>Cítrico</option><option>Fresco</option><option>Gourmand</option></select></label>
    <label class="field"><span>Género</span><select class="select" bind:value={form.gender_tag}><option>Unisex</option><option>Masculino</option><option>Femenino</option></select></label>
    <label class="field"><span>Concentración</span><select class="select" bind:value={form.concentration}><option>Extrait de Parfum</option><option>EDP</option><option>EDT</option><option>EDC</option></select></label>
    <label class="field"><span>Orden de aparición</span><input class="input" type="number" bind:value={form.display_order} /></label>
  </div>
  <label class="check"><input type="checkbox" bind:checked={form.featured} /> Fragancia destacada</label>

  <section class="panel">
    <h2>Notas olfativas</h2>
    {#each [['Salida', 'top_notes', 'top'], ['Corazón', 'heart_notes', 'heart'], ['Fondo', 'base_notes', 'base']] as row}
      {@const group = row[1] as 'top_notes' | 'heart_notes' | 'base_notes'}
      {@const noteKey = row[2] as 'top' | 'heart' | 'base'}
      <label class="field"><span>{row[0]}</span><input class="input" bind:value={notes[noteKey]} on:keydown={(event) => keyNote(event, group)} placeholder="Enter o coma para agregar" /></label>
      <div class="tags">{#each form[group] as note}<button type="button" on:click={() => removeNote(group, note)}>{note} ×</button>{/each}</div>
    {/each}
  </section>

  <section class="panel">
    <div class="row-head"><h2>Variantes</h2><button class="btn" type="button" on:click={addVariant}>+</button></div>
    {#if form.variants.length === 0}<p class="empty">Agregá al menos una variante para vender este producto.</p>{/if}
    {#each form.variants as variant, i}
      <div class="variant-row">
        <label class="field compact"><span>ml</span><input class="input" type="number" bind:value={variant.size_ml} placeholder="50" /></label>
        <label class="field compact"><span>Precio ARS</span><input class="input" type="number" value={Math.round(variant.price_ars_cents / 100)} on:input={(event) => variant.price_ars_cents = Number((event.currentTarget as HTMLInputElement).value) * 100} placeholder="89000" /></label>
        <label class="field compact"><span>Stock</span><input class="input" type="number" bind:value={variant.stock} placeholder="10" /></label>
        <label class="field compact"><span>Peso g</span><input class="input" type="number" bind:value={variant.weight_grams} placeholder="200" title="Se usa para cotizar envíos." /></label>
        <label class="field compact"><span>SKU</span><input class="input" bind:value={variant.sku} placeholder={`${form.slug}-${variant.size_ml}`} /></label>
        <button type="button" on:click={() => form.variants = form.variants.filter((_, index) => index !== i)}>×</button>
      </div>
    {/each}
  </section>

  <section class="panel">
    <div class="row-head"><h2>Imágenes</h2><button class="btn" type="button" on:click={addImage}>+</button></div>
    <p class="hint">Subí una imagen o pegá una URL pública HTTPS.</p>
    {#if form.images.length === 0}<p class="empty">Agregá una imagen principal para que el producto se vea en el catálogo.</p>{/if}
    {#each form.images as image, i}
      <div class="image-row">
        <div class="image-preview">{#if image.url}<img src={image.url} alt={image.alt_text} />{/if}</div>
        <input class="input" bind:value={image.url} placeholder="https://..." title="Pegá una URL pública HTTPS." />
        <label class:disabled={uploadingImageIndex !== null} class="upload-control">
          <span>{uploadingImageIndex === i ? 'Subiendo...' : 'Subir'}</span>
          <input type="file" accept="image/jpeg,image/png,image/webp" disabled={uploadingImageIndex !== null} on:change={(event) => uploadImage(event, i)} />
        </label>
        <input class="input" bind:value={image.alt_text} placeholder="Texto alternativo" />
        <label><input type="radio" name="primary" checked={image.is_primary} on:change={() => setPrimary(i)} /> Principal</label>
        <button type="button" on:click={() => form.images = form.images.filter((_, index) => index !== i)}>×</button>
      </div>
    {/each}
  </section>

  <div class="actions">
    {#if error}<p class="error">{error}</p>{/if}
    <button class="btn primary" type="submit" disabled={saving || uploadingImageIndex !== null}>{saving ? 'Guardando...' : 'Guardar producto'}</button>
    {#if id}<button class="btn danger" type="button" on:click={deleteProduct}>Eliminar</button>{/if}
  </div>
</form>
{/if}

<style>
  .form { display: grid; gap: 22px; max-width: 980px; }
  .panel { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 14px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .check { color: var(--color-text-muted); display: flex; gap: 10px; }
  .tags { display: flex; flex-wrap: wrap; gap: 8px; }
  .tags button, .variant-row button, .image-row button { border: 1px solid var(--color-border); color: var(--color-text); background: transparent; min-height: 36px; }
  .row-head { display: flex; justify-content: space-between; align-items: center; }
  small { color: var(--color-text-muted); font-size: .74rem; line-height: 1.4; }
  .hint, .empty { margin: 0; color: var(--color-text-muted); }
  .compact { gap: 5px; }
  .variant-row { display: grid; grid-template-columns: 80px 130px 90px 90px 1fr 44px; gap: 10px; align-items: end; }
  .image-row { display: grid; grid-template-columns: 72px minmax(180px, 1fr) 108px minmax(160px, 1fr) 120px 44px; gap: 10px; align-items: center; }
  .image-preview { width: 72px; height: 72px; border: 1px solid var(--color-border); background: rgba(255,255,255,.03); }
  .image-preview img { width: 100%; height: 100%; object-fit: cover; display: block; }
  .upload-control { min-height: 46px; padding: 0 14px; border: 1px solid var(--color-border); color: var(--color-text); display: inline-flex; align-items: center; justify-content: center; position: relative; overflow: hidden; cursor: pointer; white-space: nowrap; }
  .upload-control input { position: absolute; inset: 0; opacity: 0; cursor: pointer; }
  .upload-control.disabled { opacity: .55; pointer-events: none; }
  .actions { display: flex; gap: 12px; }
  .error { color: var(--color-danger-soft); margin: 0; align-self: center; }
  .danger { border-color: var(--color-danger-soft); color: var(--color-danger-soft); }
  @media (max-width: 800px) { .variant-row, .image-row { grid-template-columns: 1fr; } }
</style>
