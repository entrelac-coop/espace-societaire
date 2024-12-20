<script lang="ts">
  import { currentToast } from "../toast";
</script>

{#if $currentToast}
  <div
    class="fixed bottom-0 inset-x-0 mx-auto mb-4 bg-red-700 text-white max-w-md p-3 border-2 border-shadow"
    class:bg-green-700={$currentToast.type === "success"}
  >
    <p>{$currentToast.message}</p>
  </div>
{/if}
