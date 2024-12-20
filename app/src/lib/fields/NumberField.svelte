<script lang="ts">
  export let label: string;
  export let name: string;
  export let value: number;
  export let min: string;
</script>

<div>
  <label for={name} class="font-medium">{label}</label>

  <input required type="number" id={name} {name} bind:value {min} />
</div>

<style lang="postcss">
  input {
    @apply block border-2 border-shadow w-full max-w-md;
  }
</style>
