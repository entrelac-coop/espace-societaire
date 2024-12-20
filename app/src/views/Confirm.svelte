<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { meta, router } from "tinro";

  import { TextField } from "../lib/fields";
  import Button from "../lib/Button.svelte";
  import Form from "../lib/Form.svelte";

  import toast from "../toast";
  import auth from "../auth";
  import { confirmUser, startConfirmUser } from "../api";

  const route = meta();
  const email = decodeURIComponent(route.query.email);

  let code = "";

  const confirmMutation = useMutation(confirmUser, {
    onSuccess({ token }) {
      toast.success("Votre compte a bien été confirmé.");
      auth.setToken(token);
      router.goto("/");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  const startConfirmMutation = useMutation(startConfirmUser, {
    onSuccess() {
      toast.success("Un nouveau code vient de vous être envoyé.");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    $confirmMutation.mutate({
      email,
      token: code,
    });
  }

  async function startConfirm() {
    if ($startConfirmMutation.isLoading) {
      return;
    }

    $startConfirmMutation.mutate({ email });
  }
</script>

<h1>Inscription</h1>

<h2>2. Confirmation de votre compte</h2>

<p class="mb-3">
  Un email a été envoyé à l'adresse <strong>{email}</strong>. Cet email contient
  un code que vous devez entrer ici pour confirmer votre compte.
</p>

<Form
  on:submit={submit}
  loading={$confirmMutation.isLoading || $confirmMutation.isSuccess}
>
  <TextField name="token" label="Code" bind:value={code} />
</Form>

<Button
  class="mt-3"
  link
  loading={$startConfirmMutation.isLoading}
  on:click={startConfirm}>Renvoyer un email</Button
>
