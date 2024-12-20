<script lang="ts">
  let className = "";
  export let type: string = "button";
  export let href: string | null = null;
  export let loading: boolean = false;
  export let link: boolean = false;

  export { className as class };
</script>

{#if loading}
  <div class={`button loading ${className}`} class:link>
    <span>Chargement...</span>
  </div>
{:else if href}
  <a class={`button ${className}`} {href} class:link>
    <slot />
  </a>
{:else}
  <button class={`button ${className}`} on:click {type} class:link>
    <slot />
  </button>
{/if}

<style lang="postcss">
  .button:not(.link) {
    @apply inline-block px-7 py-3 text-sm font-bold uppercase bg-button text-shadow border-2 border-shadow transition;
    box-shadow: 6px 6px 0px 0px #161511;
  }

  .button.loading:not(.link) {
    @apply bg-white;
    box-shadow: none;
  }

  .button:not(.loading):not(.link):hover {
    @apply bg-primary text-white;
  }

  a.button {
    @apply no-underline;
  }
</style>
