<script lang="ts">
  import { apiFetch, defaultHomepage } from '$lib/api/client';
  import type { HomepageSettings } from '$lib/api/client';
  import { onMount } from 'svelte';
  let form: HomepageSettings = defaultHomepage;
  onMount(async () => { const res = await apiFetch<{items: HomepageSettings[]}>('/api/admin/homepage'); form = res.items[0] ?? defaultHomepage; });
  async function save(){ await apiFetch('/api/admin/homepage',{method:'PUT',body:JSON.stringify(form)}); }
</script>
<h1 class="display section-title">Homepage</h1>
<form on:submit|preventDefault={save}>{#each Object.keys(form) as key}<label class="field"><span>{key}</span><textarea class="textarea" bind:value={form[key as keyof HomepageSettings]}></textarea></label>{/each}<button class="btn primary">Guardar</button></form>
{#if form.hero_image_url}<img class="preview" src={form.hero_image_url} alt="Preview" />{/if}
<style>form{display:grid;gap:14px;max-width:760px}.preview{margin-top:24px;width:320px;aspect-ratio:4/3;object-fit:cover}</style>
