<script lang="ts">
  import AccordionFAQ from '$lib/components/AccordionFAQ.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head>
  <title>Preguntas frecuentes | ELIXIR Exclusive</title>
  <meta name="description" content="Preguntas frecuentes sobre fragancias, pagos, envíos y postventa de ELIXIR Exclusive." />
</svelte:head>

<section class="container page-pad">
  <p class="eyebrow">Ayuda</p>
  <h1 class="display section-title">Preguntas frecuentes</h1>
  <div class="gold-rule"></div>
  <AccordionFAQ items={data.settings.faq_items} />
</section>
