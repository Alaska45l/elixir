<script lang="ts">
  type FAQ = { question: string; answer: string };
  export let items: FAQ[] = [];
</script>

<div class="faq">
  {#each items as item}
    <details>
      <summary>{item.question}</summary>
      <p>{item.answer}</p>
    </details>
  {/each}
</div>

<style>
  details { border-top: 1px solid var(--color-border); padding: 18px 0; }
  summary { cursor: pointer; color: var(--color-text); }
  p { color: var(--color-text-muted); max-width: 70ch; line-height: 1.7; }
</style>
