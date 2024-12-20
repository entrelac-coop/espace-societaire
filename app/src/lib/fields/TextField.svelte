<script lang="ts">
  export let label: string;
  export let name: string;
  export let value: string;
  export let rows: number = 1;
  export let required: boolean = true;
</script>

<div>
  <label for={name} class="font-medium">{label}</label>

  {#if rows <= 1}
    <input {required} type="text" id={name} {name} bind:value />
  {:else}
    <textarea {required} id={name} {name} bind:value {rows} />
  {/if}
</div>

<style lang="postcss">
  input,
  textarea {
    @apply block border-2 border-shadow w-full max-w-md;
  }
</style>
