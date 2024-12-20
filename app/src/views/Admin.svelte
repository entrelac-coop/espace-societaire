<script lang="ts">
  import { Route } from "tinro";
  import { isAdmin } from "../auth";
  import Button from "../lib/Button.svelte";
  import UsersView from "./admin/Users.svelte";
  import UserView from "./admin/User.svelte";
</script>

{#if $isAdmin}
  <p><a href="/admin" class="text-xl font-bold mb-2">Administration</a></p>

  <Route path="/" redirect="/admin/users" />

  <Route path="/users">
    <UsersView />
  </Route>

  <Route path="/users/:userID" let:meta>
    <UserView userID={meta.params.userID} />
  </Route>
{:else}
  <h1>Accès refusé</h1>

  <p>Vous devez être administrateur.ice pour accéder à cette page.</p>

  <Button class="mt-5" href="/">Aller à l'accueil</Button>
{/if}
