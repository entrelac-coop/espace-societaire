<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { router } from "tinro";

  import { EmailField } from "../lib/fields";
  import Form from "../lib/Form.svelte";

  import toast from "../toast";
  import { startResetUser } from "../api";

  let email = "";

  const mutation = useMutation(startResetUser, {
    onSuccess() {
      router.goto("/reset");
      router.location.query.set("email", encodeURIComponent(email));
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    if ($mutation.isLoading || $mutation.isSuccess) {
      return;
    }

    $mutation.mutate({ email });
  }
</script>

<h1>J'ai oublié mon mot de passe</h1>

<p class="mb-3">
  Nous allons vous envoyer un email contenant un code. Ce code vous permettra de
  choisir un nouveau mot de passe.
</p>

<Form on:submit={submit} loading={$mutation.isLoading || $mutation.isSuccess}>
  <EmailField name="email" label="Adresse email" bind:value={email} />
</Form>
