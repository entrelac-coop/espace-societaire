<script lang="ts">
  export let label: string;
  export let name: string;
  export let value: string;
</script>

<div>
  <label for={name} class="font-medium">{label}</label>

  <select required id={name} {name} bind:value>
    <slot />
  </select>
</div>

<style lang="postcss">
  select {
    @apply block border-2 border-shadow w-full max-w-md;
  }
</style>
