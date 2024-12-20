<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import Button from "./Button.svelte";

  const dispatch = createEventDispatcher<{ submit: void }>();

  let className: string = "";
  export { className as class };

  export let loading: boolean = false;
  export let button: string = "Envoyer";

  function submit(): void {
    if (loading) {
      return;
    }

    dispatch("submit");
  }
</script>

<form
  on:submit|preventDefault={submit}
  class={`flex flex-col gap-3 ${className}`}
>
  <slot />

  <div>
    <Button {loading} type="submit">{button}</Button>
  </div>
</form>
