<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { router } from "tinro";
  import toast from "../toast";
  import { createUser } from "../api";
  import categories from "../categories";

  import {
    TextField,
    PasswordField,
    EmailField,
    SelectField,
  } from "../lib/fields";
  import Form from "../lib/Form.svelte";

  let form = {
    email: "",
    password: "",
    lastName: "",
    firstName: "",
    birthDate: "",
    address: "",
    postalCode: "",
    city: "",
    country: "",
    phoneNumber: "",
    category: "supporters",
    reason: "",
  };

  const mutation = useMutation(createUser, {
    onSuccess(_, { email }) {
      toast.success("Votre inscription a bien été enregistrée.");
      router.goto(`/confirm?email=${encodeURIComponent(email)}`);
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  function submit() {
    $mutation.mutate({
      ...form,
      reason: form.category === "supporters" ? null : form.reason,
    });
  }
</script>

<h1>Inscription</h1>

<h2>1. Création de votre compte</h2>

<p class="mb-3">
  L'inscription se déroule en quatre étapes : création de votre compte,
  confirmation de votre compte, téléversement d'une pièce d'identité et d'un
  justificatif de domicile et enfin paiement.
</p>

<p class="mb-3">
  Vous n'êtes pas obligé.e de compléter les quatre étapes d'un seul coup.
</p>

<Form on:submit={submit} loading={$mutation.isLoading || $mutation.isSuccess}>
  <div class="grid grid-cols-1 md:grid-cols-2 md:gap-12">
    <div class="flex flex-col gap-2">
      <EmailField label="Adresse email" name="email" bind:value={form.email} />
      <PasswordField
        label="Mot de passe"
        name="password"
        bind:value={form.password}
      />
      <TextField label="Nom" name="lastName" bind:value={form.lastName} />
      <TextField label="Prénom" name="firstName" bind:value={form.firstName} />
      <TextField
        label="Numéro de téléphone"
        name="phoneNumber"
        bind:value={form.phoneNumber}
      />
      <TextField label="Adresse" name="address" bind:value={form.address} />
      <TextField
        label="Code postal"
        name="postalCode"
        bind:value={form.postalCode}
      />
      <TextField label="Ville" name="city" bind:value={form.city} />
      <TextField label="Pays" name="country" bind:value={form.country} />
    </div>

    <div class="flex flex-col gap-2 mt-2">
      <SelectField label="Catégorie" name="category" bind:value={form.category}>
        {#each Object.entries(categories) as [id, name]}
          <option value={id}>
            {name}
          </option>
        {/each}
      </SelectField>

      {#if form.category !== "supporters"}
        <TextField
          label="Raison de ce choix"
          name="reason"
          bind:value={form.reason}
          rows={5}
        />
      {/if}
    </div>
  </div>
</Form>
