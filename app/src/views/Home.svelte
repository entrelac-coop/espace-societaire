<script lang="ts">
  import { token } from "../auth";

  import SignedOutHome from "../lib/SignedOutHome.svelte";
  import SignedInHome from "../lib/SignedInHome.svelte";
</script>

{#if $token}
  <SignedInHome />
{:else}
  <SignedOutHome />
{/if}
