<script lang="ts">
  import { env } from '$env/dynamic/public';
  import AccordionFAQ from '$lib/components/AccordionFAQ.svelte';
  import { apiFetch } from '$lib/api/client';
  import { toast } from '$lib/stores/toast';
  let name = '';
  let email = '';
  let subject = 'Consulta general';
  let message = '';
  const faqs = [
    { question: '¿Los perfumes son originales?', answer: 'Sí. ELIXIR Exclusive comercializa fragancias seleccionadas y documentadas.' },
    { question: '¿Qué medios de pago aceptan?', answer: 'El checkout opera en ARS mediante MercadoPago.' },
    { question: '¿Hacen envíos?', answer: 'Sí, a CABA, GBA e Interior con seguimiento.' },
    { question: '¿Puedo consultar por WhatsApp?', answer: 'Sí. Recomendamos WhatsApp para asesoramiento rápido.' },
    { question: '¿Qué pasa si una fragancia está agotada?', answer: 'Podés escribirnos para lista de reposición.' },
    { question: '¿Cómo aplico un descuento?', answer: 'Ingresá el código en el carrito y presioná Aplicar.' }
  ];
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
    <AccordionFAQ items={faqs} />
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
  form { display: grid; gap: 18px; border: 1px solid var(--color-border); padding: 24px; background: var(--color-surface); }
  .email { color: var(--color-text-muted); margin: 22px 0 42px; }
  @media (max-width: 900px) { .contact { grid-template-columns: 1fr; } }
</style>
