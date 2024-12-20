<script lang="ts">
  import { useQuery } from "@sveltestack/svelte-query";
  import { router } from "tinro";
  import { getUsers } from "../../api";
  import categories from "../../categories";

  let tab: "accepted" | "not-accepted" | "incomplete" = "not-accepted";

  const result = useQuery("users", getUsers);

  let users: any[];
  $: users = $result.data?.filter(
    (user: any) =>
      (user.accepted && tab === "accepted") ||
      (user.shares > 0 && !user.accepted && tab === "not-accepted") ||
      (user.shares === 0 && !user.accepted && tab === "incomplete")
  );
</script>

<h1>Sociétaires</h1>

<div class="mb-3">
  <button
    class="tab"
    class:active={tab === "not-accepted"}
    on:click={() => (tab = "not-accepted")}>En attente</button
  >
  <button
    class="tab"
    class:active={tab === "accepted"}
    on:click={() => (tab = "accepted")}>Validé.e.s</button
  >
  <button
    class="tab"
    class:active={tab === "incomplete"}
    on:click={() => (tab = "incomplete")}>Incomplets</button
  >
</div>

{#if tab === "not-accepted"}
  <p class="mb-3">
    Parcours d'inscription complété et parts sociales achetées, mais demande pas
    encore validée.
  </p>
{:else if tab === "accepted"}
  <p class="mb-3">
    Parcours d'inscription complété, parts sociales achetées et demande validée.
  </p>
{:else}
  <p class="mb-3">
    Parcours d’inscription incomplet ou pas de parts sociales achetées.
  </p>
{/if}

{#if $result.isLoading}
  <span>Chargement...</span>
{:else if $result.isError}
  <span>Une erreur est survenue.</span>
{:else}
  <table class="table-auto border-collapse border-primary text-left w-full">
    <thead>
      <tr>
        <th class="border border-primary p-2">Adresse email</th>
        <th class="border border-primary p-2">Prénom NOM</th>
        <th class="border border-primary p-2">Catégorie</th>
      </tr>
    </thead>

    <tbody>
      {#each users as user}
        <tr
          class="hover:bg-primary hover:text-white cursor-pointer"
          on:click={() => router.goto(`/admin/users/${user.id}`)}
        >
          <td class="border border-primary p-2">{user.email}</td>
          <td class="border border-primary p-2">
            {user.firstName}
            {user.lastName.toUpperCase()}
          </td>
          <td class="border border-primary p-2">
            {categories[user.category]}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style lang="postcss">
  .tab {
    @apply bg-button px-2 py-1 text-shadow border-shadow border-2 transition;
  }

  .tab.active,
  .tab:hover {
    @apply bg-primary text-white;
  }
</style>
