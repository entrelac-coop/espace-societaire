<script lang="ts">
  export let label: string;
  export let name: string;
  export let value: string;
</script>

<div>
  <label for={name} class="font-medium">{label}</label>

  <input required type="email" id={name} {name} bind:value />
</div>

<style lang="postcss">
  input {
    @apply block border-2 border-shadow w-full max-w-md;
  }
</style>
