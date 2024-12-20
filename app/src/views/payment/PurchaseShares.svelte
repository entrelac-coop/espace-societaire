<script lang="ts">
  import { useMutation } from "@sveltestack/svelte-query";
  import { createCheckoutSession } from "../../api";
  import toast from "../../toast";

  import { NumberField } from "../../lib/fields";
  import Form from "../../lib/Form.svelte";

  let quantity = 1;
  let gift = false;

  const mutation = useMutation(createCheckoutSession, {
    onSuccess({ url }) {
      window.location.href = url;
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    if ($mutation.isLoading || $mutation.isSuccess) {
      return;
    }

    $mutation.mutate({
      quantity,
      gift,
    });
  }
</script>

<h1>Acheter des parts sociales</h1>

<p class="mb-3">Combien de parts sociales souhaitez-vous acheter ?</p>

<Form
  on:submit={submit}
  loading={$mutation.isLoading || $mutation.isSuccess}
  button="Continuer"
>
  <div>
    <NumberField
      bind:value={quantity}
      label="Quantité"
      name="quantity"
      min="1"
    />

    <div class="mt-4 text-lg">
      {#if quantity === 1}
        <p>Une part sociale coûte <strong>50 €</strong>.</p>
      {:else}
        <p>
          {quantity} parts sociales coûtent <strong>{quantity * 50} €</strong>
          ({quantity} × 50 €).
        </p>
      {/if}
    </div>
  </div>

  <div>
    <label class="mb-3 inline-flex items-center">
      <input type="checkbox" class="cursor-pointer mr-2" bind:checked={gift} />
      <span class="font-bold">Est-ce pour offrir ?</span>
    </label>

    {#if gift}
      <p class="mb-3">Vous obtiendrez une carte cadeau à imprimer contenant un code. Le ou la destinataire du cadeau pourra utiliser ce code lors de son inscription à l’espace sociétaire pour obtenir les parts sociales.</p>
    {/if}
  </div>
</Form>
