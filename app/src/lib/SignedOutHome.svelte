<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { router } from "tinro";
  import { setToken } from "../auth";
  import { createToken } from "../api";
  import toast from "../toast";
  import { EmailField, PasswordField } from "../lib/fields";
  import Form from "../lib/Form.svelte";
  import Button from "../lib/Button.svelte";

  const form = {
    email: "",
    password: "",
  };

  const mutation = useMutation(createToken, {
    onSuccess({ token }) {
      setToken(token);
      toast.success("Vous êtes bien connecté.e.");
    },
    onError(error: any, { email }) {
      if (error.code === "not-confirmed") {
        toast.error("Veuillez confirmer votre compte.");
        router.goto("/confirm");
        router.location.query.set("email", encodeURIComponent(email));
        return;
      }

      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    $mutation.mutate(form);
  }
</script>

<h1>Espace sociétaire</h1>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
  <div>
    <h2>Je suis déjà sociétaire</h2>

    <Form
      on:submit={submit}
      loading={$mutation.isLoading || $mutation.isSuccess}
      button="Connexion"
    >
      <EmailField name="email" label="Adresse email" bind:value={form.email} />
      <PasswordField
        name="password"
        label="Mot de passe"
        bind:value={form.password}
      />
    </Form>

    <a class="mt-3 inline-block" href="/reset/start"
      >J'ai oublié mon mot de passe</a
    >
  </div>

  <div>
    <h2>Je ne suis pas encore sociétaire</h2>

    <p class="mb-6">
      Si vous avez commencé votre inscription sans la terminer, vous pouvez
      continuer celle-ci en vous connectant.
    </p>

    <Button href="/sign-up">Devenir sociétaire</Button>
  </div>
</div>
