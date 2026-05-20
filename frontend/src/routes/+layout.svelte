<script lang="ts">
  import '../app.css';
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import CartDrawer from '$lib/components/CartDrawer.svelte';
  import Toast from '$lib/components/Toast.svelte';
  import { page } from '$app/stores';
  let cartOpen = false;
  $: isAdmin = $page.url.pathname.startsWith('/admin');
</script>

{#if !isAdmin}<Header onCart={() => cartOpen = true} />{/if}
<slot />
{#if !isAdmin}<Footer />{/if}
<CartDrawer open={cartOpen} onClose={() => cartOpen = false} />
<Toast />
