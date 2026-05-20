<script lang="ts">
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { ImageFormValue, ProductFormValue, VariantFormValue } from '$lib/types/product-form';
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
    variants: [{ size_ml: 50, price_ars_cents: 8900000, stock: 10, sku: '' }],
    images: [{ url: '', alt_text: '', is_primary: true, sort_order: 0 }]
  };

  let form: ProductFormValue = product ? structuredClone(product) : structuredClone(blank);
  let notes = { top: '', heart: '', base: '' };

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
    form.variants = [...form.variants, { size_ml: 100, price_ars_cents: 0, stock: 0, sku: `${form.slug}-100` }];
  }
  function addImage() {
    form.images = [...form.images, { url: '', alt_text: form.name, is_primary: false, sort_order: form.images.length }];
  }
  function setPrimary(index: number) {
    form.images = form.images.map((img, i) => ({ ...img, is_primary: i === index }));
  }
  async function save() {
    await apiFetch(id ? `/api/admin/products/${id}` : '/api/admin/products', {
      method: id ? 'PUT' : 'POST',
      body: JSON.stringify(form)
    });
    toast.push('Producto guardado');
    location.href = '/admin/productos';
  }
  async function deleteProduct() {
    if (!id || !confirm('¿Eliminar este producto?')) return;
    await apiFetch(`/api/admin/products/${id}`, { method: 'DELETE' });
    toast.push('Producto eliminado');
    location.href = '/admin/productos';
  }
</script>

<form class="form" on:submit|preventDefault={save}>
  <div class="grid-2">
    <label class="field"><span>Nombre</span><input class="input" required bind:value={form.name} on:input={nameInput} /></label>
    <label class="field"><span>Slug</span><input class="input" required bind:value={form.slug} /></label>
  </div>
  <label class="field"><span>Tagline</span><input class="input" bind:value={form.tagline} /></label>
  <label class="field"><span>Descripción</span><textarea class="textarea" bind:value={form.description}></textarea></label>
  <div class="grid-2">
    <label class="field"><span>Familia</span><select class="select" bind:value={form.scent_family}><option>Oriental</option><option>Floral</option><option>Amaderado</option><option>Cítrico</option><option>Fresco</option><option>Gourmand</option></select></label>
    <label class="field"><span>Género</span><select class="select" bind:value={form.gender_tag}><option>Unisex</option><option>Masculino</option><option>Femenino</option></select></label>
    <label class="field"><span>Concentración</span><select class="select" bind:value={form.concentration}><option>Extrait de Parfum</option><option>EDP</option><option>EDT</option><option>EDC</option></select></label>
    <label class="field"><span>Orden</span><input class="input" type="number" bind:value={form.display_order} /></label>
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
    {#each form.variants as variant, i}
      <div class="variant-row">
        <input class="input" type="number" bind:value={variant.size_ml} placeholder="ml" />
        <input class="input" type="number" value={Math.round(variant.price_ars_cents / 100)} on:input={(event) => variant.price_ars_cents = Number((event.currentTarget as HTMLInputElement).value) * 100} placeholder="ARS" />
        <input class="input" type="number" bind:value={variant.stock} placeholder="Stock" />
        <input class="input" bind:value={variant.sku} placeholder="SKU" />
        <button type="button" on:click={() => form.variants = form.variants.filter((_, index) => index !== i)}>×</button>
      </div>
    {/each}
  </section>

  <section class="panel">
    <div class="row-head"><h2>Imágenes</h2><button class="btn" type="button" on:click={addImage}>+</button></div>
    {#each form.images as image, i}
      <div class="image-row">
        {#if image.url}<img src={image.url} alt={image.alt_text} />{/if}
        <input class="input" bind:value={image.url} placeholder="URL" />
        <input class="input" bind:value={image.alt_text} placeholder="Alt text" />
        <label><input type="radio" name="primary" checked={image.is_primary} on:change={() => setPrimary(i)} /> Principal</label>
        <button type="button" on:click={() => form.images = form.images.filter((_, index) => index !== i)}>×</button>
      </div>
    {/each}
  </section>

  <div class="actions">
    <button class="btn primary" type="submit">Guardar producto</button>
    {#if id}<button class="btn danger" type="button" on:click={deleteProduct}>Eliminar</button>{/if}
  </div>
</form>

<style>
  .form { display: grid; gap: 22px; max-width: 980px; }
  .panel { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 14px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .check { color: var(--color-text-muted); display: flex; gap: 10px; }
  .tags { display: flex; flex-wrap: wrap; gap: 8px; }
  .tags button, .variant-row button, .image-row button { border: 1px solid var(--color-border); color: var(--color-text); background: transparent; min-height: 36px; }
  .row-head { display: flex; justify-content: space-between; align-items: center; }
  .variant-row { display: grid; grid-template-columns: 90px 140px 100px 1fr 44px; gap: 10px; }
  .image-row { display: grid; grid-template-columns: 72px 1fr 1fr 120px 44px; gap: 10px; align-items: center; }
  .image-row img { width: 72px; height: 72px; object-fit: cover; }
  .actions { display: flex; gap: 12px; }
  .danger { border-color: #9f5d55; color: #e0a39a; }
  @media (max-width: 800px) { .variant-row, .image-row { grid-template-columns: 1fr; } }
</style>
