<script lang="ts">
  import { slugify } from '$lib/utils/slug';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  let name = '';
  let slug = '';
  let price = 89000;
  let stock = 10;
  async function save() {
    await apiFetch('/api/admin/products', {
      method: 'POST',
      body: JSON.stringify({
        name, slug, tagline: 'Fragancia de autor', description: 'Descripción editorial del perfume.',
        scent_family: 'Oriental', gender_tag: 'Unisex', concentration: 'EDP',
        top_notes: ['Bergamota'], heart_notes: ['Rosa'], base_notes: ['Ámbar'],
        featured: true, display_order: 0,
        variants: [{ size_ml: 50, price_ars_cents: price * 100, stock, sku: `${slug}-50` }],
        images: [{ url: 'https://images.unsplash.com/photo-1592945403244-b3fbafd7f539?auto=format&fit=crop&w=900&q=85', alt_text: name, is_primary: true, sort_order: 0 }]
      })
    });
    toast.push('Producto guardado');
  }
</script>
<form class="form" on:submit|preventDefault={save}>
  <label class="field"><span>Nombre</span><input class="input" required bind:value={name} on:input={() => slug = slug || slugify(name)} /></label>
  <label class="field"><span>Slug</span><input class="input" required bind:value={slug} /></label>
  <label class="field"><span>Precio ARS</span><input class="input" type="number" required bind:value={price} /></label>
  <label class="field"><span>Stock</span><input class="input" type="number" required bind:value={stock} /></label>
  <button class="btn primary" type="submit">Guardar producto</button>
</form>
<style>.form{display:grid;gap:18px;max-width:680px}</style>
