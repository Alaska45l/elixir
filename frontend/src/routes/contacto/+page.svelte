<script lang="ts">
  import { env } from '$env/dynamic/public';
  import AccordionFAQ from '$lib/components/AccordionFAQ.svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  import type { PageData } from './$types';
  export let data: PageData;
  let name = '';
  let email = '';
  let subject = 'Consulta general';
  let message = '';
  async function submit() {
    await apiFetch('/api/contact', { method: 'POST', body: JSON.stringify({ name, email, subject, message }) });
    toast.push('Mensaje enviado');
    name = email = message = '';
  }
</script>

<svelte:head><title>Contacto | ELIXIR Exclusive</title></svelte:head>
<section class="container page-pad contact">
  <div>
    <p class="eyebrow">Contacto</p>
    <h1 class="display section-title">Asesoramiento privado</h1>
    <div class="gold-rule"></div>
    <a class="btn primary" href={`https://wa.me/${env.PUBLIC_WHATSAPP_NUMBER ?? '5491100000000'}`} target="_blank" rel="noreferrer">Consultar por WhatsApp</a>
    <p class="email"><bdo dir="rtl">ra.moc.evisulcxerixile@otcatnoc</bdo></p>
    <AccordionFAQ items={data.settings.faq_items} />
  </div>
  <form on:submit|preventDefault={submit}>
    <label class="field"><span>Nombre</span><input class="input" required bind:value={name} /></label>
    <label class="field"><span>Email</span><input class="input" type="email" required bind:value={email} /></label>
    <label class="field"><span>Asunto</span><select class="select" bind:value={subject}><option>Consulta general</option><option>Envíos</option><option>Pagos</option><option>Postventa</option></select></label>
    <label class="field"><span>Mensaje</span><textarea class="textarea" required bind:value={message}></textarea></label>
    <button class="btn primary" type="submit">Enviar mensaje</button>
  </form>
</section>

<style>
  .contact { display: grid; grid-template-columns: 1fr 430px; gap: 64px; align-items: start; }
  form { display: grid; gap: 18px; border-radius: 14px; padding: 24px; background: var(--color-surface); box-shadow: 0 4px 24px rgba(0,0,0,0.15); }
  .email { color: var(--color-text-muted); margin: 22px 0 42px; }
  @media (max-width: 900px) { .contact { grid-template-columns: 1fr; } }
</style>
