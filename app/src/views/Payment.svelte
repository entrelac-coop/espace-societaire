<script lang="ts">
  import { meta, Route } from "tinro";
  import Button from "../lib/Button.svelte";
  import PurchaseSharesView from "./payment/PurchaseShares.svelte";
  
  const route = meta();
  const giftID = route.query.giftID;
</script>

<Route path="/">
  <PurchaseSharesView />
</Route>

<Route path="/success">
  <h1>C'est tout bon !</h1>

  <p>Votre achat de parts sociales est confirmé.</p>
  <p class="mb-3">Une facture vient d'être envoyée à votre adresse email.</p>

  {#if giftID}
    <p><a href={`https://societaire.entrelac.coop/api/gifts/${giftID}`}>Veuillez cliquer sur ici lien pour télécharger la carte cadeau.</a></p>
  {:else}
    <p>
      Votre entrée au sociétariat sera validée lors de la prochaine réunion
      mensuelle du conseil d'administration d'Entrelac.coop, vous serez averti.e
      par email.
    </p>
  {/if}

  <div class="mt-6">
    <Button href="/">Aller à l'accueil</Button>
  </div>
</Route>

<Route path="/cancel">
  <h1>C'est dommage...</h1>

  <p>Le paiement a été annulé.</p>

  <div class="mt-6">
    <Button href="/">Aller à l'accueil</Button>
  </div>
</Route>
