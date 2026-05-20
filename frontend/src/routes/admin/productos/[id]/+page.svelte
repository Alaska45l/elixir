<script lang="ts">
  import ProductForm from '$lib/components/admin/ProductForm.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
</script>
<h1 class="display section-title">Editar producto</h1><ProductForm id={data.id} product={data.product} />
