<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";
  import { token } from "../../auth";
  import { baseURL, getUser } from "../../api";
  import categories from "../../categories";

  export let userID: string;

  const result = useQuery(["user", userID], () => getUser(userID));

  let user: any | null;
  $: user = $result.data;
</script>

{#if $result.isLoading}
  <span>Chargement...</span>
{:else if $result.isError}
  <span>Une erreur est survenue.</span>
{:else}
  <h1>
    {user.firstName}
    {user.lastName.toUpperCase()}
  </h1>

  <div class="flex flex-col gap-2">
    <div>
      <h3>Parts sociales</h3>
      <p>{user.shares}</p>
    </div>

    <div>
      <h3>Adresse email</h3>
      <p>{user.email}</p>
      {#if user.confirmed}
        <p>L'utilisateur.ice a confirmé son compte.</p>
      {:else}
        <p>L'utilisateur.ice n'a <strong>pas</strong> confirmé son compte.</p>
      {/if}
    </div>

    <div>
      <h3>Numéro de téléphone</h3>
      <p>{user.phoneNumber}</p>
    </div>

    <div>
      <h3>Adresse</h3>
      <p>{user.address}</p>
      <p>{user.postalCode} {user.city}</p>
      <p>{user.country}</p>
    </div>

    <div>
      <h3>Catégorie</h3>
      <p>{categories[user.category]}</p>

      {#if user.reason}
        <div class="border-l border-primary pl-2 mt-2">
          <h4 class="underline">Raison de ce choix</h4>
          <pre class="font-sans">{user.reason}</pre>
        </div>
      {/if}
    </div>

    <div>
      <h3>Documents</h3>
      {#if user.identityFront}
        <ul class="list-disc list-inside">
          <li>
            <a
              href={`${baseURL}admin/users/${user.id}/documents/${user.identityFront}?token=${$token}`}
              >Pièce d'identité (recto)</a
            >
          </li>
          {#if user.identityBack}
            <li>
              <a
                href={`${baseURL}admin/users/${user.id}/documents/${user.identityBack}?token=${$token}`}
                >Pièce d'identité (verso)</a
              >
            </li>
          {/if}
          <li>
            <a
              href={`${baseURL}admin/users/${user.id}/documents/${user.addressProof}?token=${$token}`}
              >Justificatif de domicile</a
            >
          </li>
        </ul>
      {:else}
        <p>
          L'utilisateur.ice n'a <strong>pas</strong> téléversé ses documents.
        </p>
      {/if}
    </div>
  </div>
{/if}
