<script lang="ts">
  import { apiFetch } from '$lib/api/client';
  let username = '';
  let password = '';
  let error = '';
  async function login() {
    try {
      await apiFetch('/api/admin/login', { method: 'POST', body: JSON.stringify({ username, password }) });
      location.href = '/admin';
    } catch (err) {
      error = err instanceof Error ? err.message : 'No se pudo iniciar sesión';
    }
  }
</script>

<section class="login">
  <form on:submit|preventDefault={login}>
    <p class="eyebrow">Panel administrativo</p>
    <h1 class="display">ELIXIR</h1>
    <label class="field"><span>Usuario</span><input class="input" required bind:value={username} /></label>
    <label class="field"><span>Contraseña</span><input class="input" type="password" required bind:value={password} /></label>
    {#if error}<p class="error">Demasiados intentos o credenciales inválidas. Intente nuevamente.</p>{/if}
    <button class="btn primary" type="submit">Ingresar</button>
  </form>
</section>
<style>
  .login { min-height: 100vh; display: grid; place-items: center; padding: 24px; }
  form { width: min(420px, 100%); display: grid; gap: 18px; border: 1px solid var(--color-border); padding: 32px; background: var(--color-surface); }
  h1 { margin: 0 0 10px; font-size: 4rem; }
  .error { color: #e0a39a; }
</style>
