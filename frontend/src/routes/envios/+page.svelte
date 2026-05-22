<script lang="ts">
  import AccordionFAQ from '$lib/components/AccordionFAQ.svelte';
  import { formatARS } from '$lib/utils/currency';
  import type { PageData } from './$types';
  export let data: PageData;
</script>

<svelte:head><title>Envíos | ELIXIR Exclusive</title></svelte:head>
<section class="container page-pad">
  <p class="eyebrow">Envíos</p>
  <h1 class="display section-title">Envíos a todo el país</h1>
  <div class="gold-rule"></div>
  <table class="table">
    <thead><tr><th>Zona</th><th>Costo</th><th>Tiempo estimado</th></tr></thead>
    <tbody>{#each data.zones as zone}<tr><td>{zone.zone_name}</td><td>{formatARS(zone.base_cost_cents)}</td><td>{zone.estimated_days_min} a {zone.estimated_days_max} días</td></tr>{/each}</tbody>
  </table>
  <div class="copy"><p>Trabajamos con operadores nacionales y referencia de Correo Argentino para cobertura federal. El tracking se informa por email y WhatsApp cuando la orden pasa a despacho.</p></div>
  <AccordionFAQ items={data.settings.faq_items} />
</section>

<style>.copy{margin:42px 0;color:var(--color-text-muted);font-size:1.1rem;line-height:1.8;max-width:760px}</style>
