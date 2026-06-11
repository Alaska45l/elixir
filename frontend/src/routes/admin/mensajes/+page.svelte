<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  type Msg = { id: string; name: string; email: string; subject: string; message: string; read: boolean };
  let messages: Msg[] = [];
  let loading = true;
  let savingId = '';
  let error = '';

  async function load() {
    loading = true;
    error = '';
    try {
      messages = (await apiFetch<{items: Msg[]}>('/api/admin/contact')).items;
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudieron cargar mensajes';
    } finally {
      loading = false;
    }
  }

  async function setRead(message: Msg, read: boolean) {
    savingId = message.id;
    error = '';
    try {
      await apiFetch(`/api/admin/contact/${message.id}/read`, {
        method: 'PUT',
        body: JSON.stringify({ read })
      });
      toast.push(read ? 'Mensaje marcado como leído' : 'Mensaje marcado como pendiente');
      await load();
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo actualizar el mensaje';
    } finally {
      savingId = '';
    }
  }

  onMount(load);
</script>
<h1 class="display section-title">Mensajes</h1>
{#if error}<p class="error">{error}</p>{/if}
{#if loading}
  <p class="empty">Cargando mensajes...</p>
{:else if messages.length === 0}
  <p class="empty">No hay mensajes de contacto.</p>
{:else}
  <div class="messages">
    {#each messages as m}
      <article class:read={m.read}>
        <strong>{m.name}</strong>
        <span>{m.email}{m.subject ? ` · ${m.subject}` : ''}</span>
        <p>{m.message}</p>
        <button class="btn" type="button" disabled={savingId === m.id} on:click={() => setRead(m, !m.read)}>
          {savingId === m.id ? 'Guardando...' : m.read ? 'Marcar como pendiente' : 'Marcar como leído'}
        </button>
      </article>
    {/each}
  </div>
{/if}
<style>.messages{display:grid;gap:16px}article{border:1px solid var(--color-border);padding:18px;background:var(--color-surface)}article.read{opacity:.62}span{display:block;color:var(--color-text-muted);margin-top:4px}p{line-height:1.6}.empty{color:var(--color-text-muted)}.error{color:var(--color-danger-soft)}</style>
