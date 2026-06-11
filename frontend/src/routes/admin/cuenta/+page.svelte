<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { AdminUser } from '$lib/types/admin';

  let users: AdminUser[] = [];
  let current_password = '';
  let new_password = '';
  let repeat_password = '';
  let username = '';
  let password = '';
  let resetPassword: Record<string, string> = {};
  let loading = true;
  let deletingUsername = '';
  let error = '';
  let userError = '';

  async function loadUsers() {
    loading = true;
    try {
      users = (await apiFetch<{ items: AdminUser[] }>('/api/admin/users')).items;
    } catch (err) {
      userError = err instanceof Error ? err.message : 'No se pudieron cargar usuarios';
    } finally {
      loading = false;
    }
  }

  async function changeOwnPassword() {
    error = '';
    if (new_password !== repeat_password) {
      error = 'Las contraseñas nuevas no coinciden.';
      return;
    }
    try {
      await apiFetch('/api/admin/password', {
        method: 'POST',
        body: JSON.stringify({ current_password, new_password })
      });
      current_password = new_password = repeat_password = '';
      toast.push('Contraseña actualizada');
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo cambiar la contraseña';
    }
  }

  async function createUser() {
    userError = '';
    try {
      await apiFetch('/api/admin/users', { method: 'POST', body: JSON.stringify({ username, password }) });
      username = password = '';
      toast.push('Usuario creado');
      await loadUsers();
    } catch (err) {
      userError = err instanceof Error ? err.message : 'No se pudo crear el usuario';
    }
  }

  async function resetUserPassword(user: AdminUser) {
    const next = resetPassword[user.username] ?? '';
    if (!next.trim()) {
      userError = 'Ingresá una contraseña nueva para ese usuario.';
      return;
    }
    try {
      await apiFetch(`/api/admin/users/${encodeURIComponent(user.username)}/password`, {
        method: 'POST',
        body: JSON.stringify({ new_password: next })
      });
      resetPassword[user.username] = '';
      toast.push('Contraseña restablecida');
    } catch (err) {
      userError = err instanceof Error ? err.message : 'No se pudo restablecer la contraseña';
    }
  }

  async function deleteUser(user: AdminUser) {
    if (deletingUsername || !confirm(`¿Eliminar el usuario ${user.username}?`)) return;
    deletingUsername = user.username;
    userError = '';
    try {
      await apiFetch(`/api/admin/users/${encodeURIComponent(user.username)}`, { method: 'DELETE' });
      toast.push('Usuario eliminado');
      await loadUsers();
    } catch (err) {
      userError = err instanceof Error ? err.message : 'No se pudo eliminar el usuario';
    } finally {
      deletingUsername = '';
    }
  }

  onMount(loadUsers);
</script>

<h1 class="display section-title">Cuenta</h1>

<div class="sections">
  <form class="panel" on:submit|preventDefault={changeOwnPassword}>
    <h2>Cambiar mi contraseña</h2>
    <label class="field"><span>Contraseña actual</span><input class="input" type="password" required bind:value={current_password} /></label>
    <div class="grid-2">
      <label class="field"><span>Nueva contraseña</span><input class="input" type="password" required minlength="10" bind:value={new_password} /></label>
      <label class="field"><span>Repetir nueva contraseña</span><input class="input" type="password" required minlength="10" bind:value={repeat_password} /></label>
    </div>
    <p class="hint">Usá al menos 10 caracteres con letras y números.</p>
    {#if error}<p class="error">{error}</p>{/if}
    <button class="btn primary" type="submit">Guardar contraseña</button>
  </form>

  <section class="panel">
    <div class="row-head"><h2>Usuarios administradores</h2></div>
    <form class="new-user" on:submit|preventDefault={createUser}>
      <label class="field"><span>Usuario</span><input class="input" required bind:value={username} /></label>
      <label class="field"><span>Contraseña inicial</span><input class="input" type="password" required minlength="10" bind:value={password} /></label>
      <button class="btn primary" type="submit">Crear usuario</button>
    </form>
    {#if userError}<p class="error">{userError}</p>{/if}
    {#if loading}
      <p class="hint">Cargando usuarios...</p>
    {:else if users.length === 0}
      <p class="empty">No hay usuarios cargados.</p>
    {:else}
      <div class="users">
        {#each users as user}
          <article>
            <div>
              <strong>{user.username}</strong>
              <span>{user.last_login_at ? `Último ingreso: ${new Date(user.last_login_at).toLocaleString('es-AR')}` : 'Todavía no ingresó'}</span>
            </div>
            <label class="field"><span>Nueva contraseña</span><input class="input" type="password" minlength="10" bind:value={resetPassword[user.username]} /></label>
            <button class="btn" type="button" disabled={Boolean(deletingUsername)} on:click={() => resetUserPassword(user)}>Restablecer</button>
            <button class="btn danger" type="button" disabled={Boolean(deletingUsername)} on:click={() => deleteUser(user)}>{deletingUsername === user.username ? 'Eliminando...' : 'Eliminar'}</button>
          </article>
        {/each}
      </div>
    {/if}
  </section>
</div>

<style>
  .sections { display: grid; gap: 28px; max-width: 1040px; }
  .panel { border-top: 1px solid var(--color-border); padding-top: 22px; display: grid; gap: 16px; }
  h2 { margin: 0; font-size: 1rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: .12em; }
  .hint, .empty { margin: 0; color: var(--color-text-muted); }
  .error { margin: 0; color: var(--color-danger-soft); }
  .new-user { display: grid; grid-template-columns: 1fr 1fr auto; gap: 12px; align-items: end; }
  .users { display: grid; gap: 12px; }
  article { display: grid; grid-template-columns: 1fr minmax(180px, 260px) auto auto; gap: 12px; align-items: end; border: 1px solid var(--color-border); background: var(--color-surface); padding: 14px; }
  article div { display: grid; gap: 5px; }
  article span { color: var(--color-text-muted); font-size: .86rem; }
  .danger { border-color: var(--color-danger-soft); color: var(--color-danger-soft); }
  @media (max-width: 900px) { .new-user, article { grid-template-columns: 1fr; } }
</style>
