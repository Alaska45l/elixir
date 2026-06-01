<script lang="ts">
  import { PUBLIC_SITE_URL } from '$env/static/public';
  import '../app.css';
  import AnnouncementBar from '$lib/components/AnnouncementBar.svelte';
  import Header from '$lib/components/Header.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import CartDrawer from '$lib/components/CartDrawer.svelte';
  import ScrollTop from '$lib/components/ScrollTop.svelte';
  import Toast from '$lib/components/Toast.svelte';
  import { navigating, page } from '$app/stores';
  import type { LayoutData } from './$types';
  export let data: LayoutData;
  let cartOpen = false;
  $: isAdmin = data.isAdmin || $page.url.pathname.startsWith('/admin');
  $: siteUrl = (PUBLIC_SITE_URL || 'http://localhost:5173').replace(/\/$/, '');
</script>

<svelte:head>
  <link rel="canonical" href={`${siteUrl}${$page.url.pathname}`} />
</svelte:head>

{#if !isAdmin}<AnnouncementBar settings={data.settings} />{/if}
{#if !isAdmin}<Header settings={data.settings} onCart={() => cartOpen = true} />{/if}
<div class:fading={$navigating}><slot /></div>
{#if !isAdmin}<Footer settings={data.settings} />{/if}
{#if !isAdmin}<CartDrawer open={cartOpen} onClose={() => cartOpen = false} />{/if}
{#if !isAdmin}<ScrollTop />{/if}
<Toast />

<style>
  div { transition: opacity .22s ease; }
  .fading { opacity: 0; transition: opacity .18s ease; }
</style>
