<script lang="ts">
  import logo from "../assets/logo.svg";

  import { QueryClient, QueryClientProvider } from "@sveltestack/svelte-query";
  import { router, Route } from "tinro";
  import { token, signOut } from "./auth";

  import Button from "./lib/Button.svelte";
  import Toast from "./lib/Toast.svelte";
  import Footer from "./lib/Footer.svelte";

  import HomeView from "./views/Home.svelte";
  import SignUpView from "./views/SignUp.svelte";
  import ConfirmView from "./views/Confirm.svelte";
  import ResetView from "./views/Reset.svelte";
  import ResetStartView from "./views/ResetStart.svelte";
  import PaymentView from "./views/Payment.svelte";
  import AdminView from "./views/Admin.svelte";

  router.subscribe((_) => window.scrollTo(0, 0));

  const queryClient = new QueryClient();
</script>

<main class="flex-grow w-full max-w-4xl mx-auto p-4 text-lg lg:mt-16">
  <div class="sm:flex sm:flex-row sm:justify-between mb-6 items-start">
    <a href="/" class="border-none">
      <img src={logo} alt="Logo entrelac.coop" class="w-40" />
    </a>

    {#if $token}
      <Button on:click={signOut}>Déconnexion</Button>
    {/if}
  </div>

  <QueryClientProvider client={queryClient}>
    <Route path="/">
      <HomeView />
    </Route>

    <Route path="/sign-up">
      <SignUpView />
    </Route>

    <Route path="/confirm">
      <ConfirmView />
    </Route>

    <Route path="/reset/*">
      <Route path="/">
        <ResetView />
      </Route>

      <Route path="/start">
        <ResetStartView />
      </Route>
    </Route>

    <Route path="/payment/*">
      <PaymentView />
    </Route>

    <Route path="/admin/*">
      <AdminView />
    </Route>
  </QueryClientProvider>
</main>

<Footer />

<Toast />

<style lang="postcss">
  :global(h1) {
    @apply text-5xl font-bold mb-6;
  }

  :global(h2) {
    @apply text-2xl font-bold mb-3;
  }

  :global(h3) {
    @apply text-xl font-bold mb-1;
  }

  :global(a),
  :global(.link) {
    @apply font-semibold inline-block;
  }

  :global(a):hover,
  :global(.link):hover {
    @apply underline;
  }
</style>
