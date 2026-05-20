<script lang="ts">
  import '../app.css';
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import CartDrawer from '$lib/components/CartDrawer.svelte';
  import ScrollTop from '$lib/components/ScrollTop.svelte';
  import Toast from '$lib/components/Toast.svelte';
  import { navigating, page } from '$app/stores';
  let cartOpen = false;
  $: isAdmin = $page.url.pathname.startsWith('/admin');
</script>

{#if !isAdmin}<Header onCart={() => cartOpen = true} />{/if}
<div class:fading={$navigating}><slot /></div>
{#if !isAdmin}<Footer />{/if}
<CartDrawer open={cartOpen} onClose={() => cartOpen = false} />
{#if !isAdmin}<ScrollTop />{/if}
<Toast />

<style>
  div { transition: opacity .22s ease; }
  .fading { opacity: 0; transition: opacity .18s ease; }
</style>
