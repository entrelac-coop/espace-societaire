<script lang="ts">
  export let label: string;
  export let name: string;
  export let files: FileList;
  export let optional: boolean = false;
</script>

<div>
  <label for={name} class="font-medium">{label}</label>

  <input
    required={!optional}
    type="file"
    id={name}
    {name}
    class="p-2"
    bind:files
  />
</div>

<style lang="postcss">
  input {
    @apply block border-2 border-shadow w-full max-w-md;
  }
</style>
