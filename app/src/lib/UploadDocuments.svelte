<script lang="ts">
  import { FileField } from "./fields";
  import Form from "./Form.svelte";

  import toast from "../toast";
  import { uploadDocuments } from "../api";
  import { useMutation, useQueryClient } from "@sveltestack/svelte-query";

  let identityFrontFiles: FileList;
  let identityBackFiles: FileList;
  let addressProofFiles: FileList;

  const queryClient = useQueryClient();

  const mutation = useMutation(uploadDocuments, {
    onSuccess() {
      queryClient.invalidateQueries("me");
      toast.success("Vos documents ont bien été enregistrés.");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  function submit() {
    const data = new FormData();

    data.append("identity_front", identityFrontFiles[0]);
    if (identityBackFiles && identityBackFiles.length > 0) {
      data.append("identity_back", identityBackFiles[0]);
    }
    data.append("address_proof", addressProofFiles[0]);

    $mutation.mutate(data);
  }
</script>

<h1>Inscription</h1>

<h2>3. Téléversement de vos documents</h2>

<p class="mb-3">
  Afin de valider votre identité, nous devons vous demander une pièce d'identité
  et un justificatif de domicile.
</p>

<Form on:submit={submit} loading={$mutation.isLoading || $mutation.isSuccess}>
  <FileField
    bind:files={identityFrontFiles}
    label="Pièce d'identité (recto)"
    name="identityFront"
  />
  <FileField
    bind:files={identityBackFiles}
    label="Pièce d'identité (verso) (facultatif)"
    name="identityBack"
    optional
  />
  <FileField
    bind:files={addressProofFiles}
    label="Justificatif de domicile"
    name="addressProof"
  />
</Form>
