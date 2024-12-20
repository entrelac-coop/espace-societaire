<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { meta, router } from "tinro";

  import { TextField, PasswordField } from "../lib/fields";
  import Button from "../lib/Button.svelte";
  import Form from "../lib/Form.svelte";

  import toast from "../toast";
  import auth from "../auth";
  import { resetUser, startResetUser } from "../api";

  const route = meta();
  const email = decodeURIComponent(route.query.email);

  let code = "";
  let password = "";

  const resetMutation = useMutation(resetUser, {
    onSuccess({ token }) {
      toast.success("Votre mot de passe a bien été mis à jour.");
      auth.setToken(token);
      router.goto("/");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  const startResetMutation = useMutation(startResetUser, {
    onSuccess() {
      toast.success("Un nouveau code vient de vous être envoyé.");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    if ($resetMutation.isLoading || $resetMutation.isSuccess) {
      return;
    }

    $resetMutation.mutate({
      email,
      password,
      token: code,
    });
  }

  async function startReset() {
    if ($startResetMutation.isLoading) {
      return;
    }

    $startResetMutation.mutate({ email });
  }
</script>

<h1>J'ai oublié mon mot de passe</h1>

<p class="mb-3">
  Un email a été envoyé à l'adresse <strong>{email}</strong>. Cet email contient
  un code que vous devez entrer ici pour choisir un nouveau mot de passe.
</p>

<Form
  on:submit={submit}
  loading={$resetMutation.isLoading || $resetMutation.isSuccess}
>
  <TextField name="token" label="Code" bind:value={code} />
  <PasswordField
    name="password"
    label="Nouveau mot de passe"
    bind:value={password}
  />
</Form>

<Button
  link
  loading={$startResetMutation.isLoading}
  on:click={startReset}
  class="mt-3">Renvoyer un email</Button
>
