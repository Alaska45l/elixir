<script lang="ts">
  import ProductForm from '$lib/components/admin/ProductForm.svelte';
  import { page } from '$app/stores';
</script>
<h1 class="display section-title">Editar producto</h1><ProductForm id={$page.params.id} />
