<script lang="ts">
  import { useMutation, useQueryClient } from "@sveltestack/svelte-query";
  import { createCheckoutSession, useGiftCode } from "../api";
  import toast from "../toast";

  import { NumberField, TextField } from "./fields";
  import Form from "./Form.svelte";

  const queryClient = useQueryClient();

  let quantity = 1;
  let giftCode = "";

  const createCheckoutSessionMutation = useMutation(createCheckoutSession, {
    onSuccess({ url }) {
      window.location.href = url;
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });
  
  const useGiftCodeMutation = useMutation(useGiftCode, {
    onSuccess() {
      queryClient.invalidateQueries("me");
    },
    onError() {
      toast.error("Une erreur est survenue, veuillez réessayer.");
    },
  });

  async function submit() {
    if ($createCheckoutSessionMutation.isLoading || $createCheckoutSessionMutation.isSuccess || $useGiftCodeMutation.isLoading || $useGiftCodeMutation.isSuccess) {
      return;
    }

    if (giftCode === "") {
      $createCheckoutSessionMutation.mutate({ quantity });
    } else {
      $useGiftCodeMutation.mutate({ giftCode });
    }
  }
</script>

<h1>Inscription</h1>

<h2>4. Paiement</h2>

<p class="mb-3">Combien de parts sociales souhaitez-vous acheter ?</p>

<Form
  on:submit={submit}
  loading={$createCheckoutSessionMutation.isLoading || $createCheckoutSessionMutation.isSuccess || $useGiftCodeMutation.isLoading || $useGiftCodeMutation.isSuccess}
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

    <div class="mt-4 mb-3 text-lg">
      <p class="font-bold">Si vous posséder un code cadeau, vous pouvez l’entrer ici au lieu d’acheter des parts.</p>
    </div>

    <TextField
      bind:value={giftCode}
      label="Code cadeau"
      name="giftCode"
      required={false}
    />
  </div>
</Form>
