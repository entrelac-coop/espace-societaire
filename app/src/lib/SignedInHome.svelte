<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";
  import { getCurrentUser } from "../api";
  import { isAdmin } from "../auth";

  import UploadDocuments from "./UploadDocuments.svelte";
  import PurchaseShares from "./PurchaseShares.svelte";
  import Button from "./Button.svelte";

  const result = useQuery("me", getCurrentUser);
</script>

{#if $result.isLoading}
  <span>Loading...</span>
{:else if $result.isError}
  <span>Une erreur est survenue.</span>
{:else if $result.data.mustUploadDocuments}
  <UploadDocuments />
{:else if $result.data.shares === 0}
  <PurchaseShares />
{:else}
  <h1>Espace sociétaire</h1>

  <p>Vous êtes connecté.e en tant que {$result.data.email}.</p>

  {#if $result.data.shares === 1}
    <p>Vous avez une part sociale.</p>
  {:else}
    <p>Vous avez {$result.data.shares} parts sociales.</p>
  {/if}

  {#if $isAdmin}
    <p>
      <a href="/admin">Cliquez ici pour accéder au panel d'administration.</a>
    </p>
  {/if}

  <Button class="mt-3" href="/payment">Acheter des parts sociales</Button>
{/if}
