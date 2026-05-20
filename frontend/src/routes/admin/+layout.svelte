<script lang="ts">
  import AdminLayout from '$lib/components/admin/AdminLayout.svelte';
  import { page } from '$app/stores';
  $: login = $page.url.pathname === '/admin/login';
</script>

{#if login}
  <slot />
{:else}
  <AdminLayout><slot /></AdminLayout>
{/if}
