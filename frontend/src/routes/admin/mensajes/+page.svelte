<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  type Msg = { id: string; name: string; email: string; subject: string; message: string; read: boolean };
  let messages: Msg[] = [];
  async function load(){ messages = (await apiFetch<{items: Msg[]}>('/api/admin/contact')).items; }
  async function mark(id: string){ await apiFetch(`/api/admin/contact/${id}/read`,{method:'PUT',body:'{}'}); await load(); }
  onMount(load);
</script>
<h1 class="display section-title">Mensajes</h1>
<div class="messages">{#each messages as m}<article class:read={m.read}><strong>{m.name}</strong><span>{m.email} · {m.subject}</span><p>{m.message}</p><button class="btn" on:click={() => mark(m.id)}>Marcar como leído</button></article>{/each}</div>
<style>.messages{display:grid;gap:16px}article{border:1px solid var(--color-border);padding:18px;background:var(--color-surface)}article.read{opacity:.62}span{display:block;color:var(--color-text-muted);margin-top:4px}p{line-height:1.6}</style>
