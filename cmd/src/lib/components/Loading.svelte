<script>
    import LoadingSpinner from "./LoadingSpinner.svelte";
    export let loading = false;
</script>

{#if loading}
    <LoadingSpinner
        class="fixed top-0 bottom-0 left-0 right-0 my-auto mx-auto w-10"
    />
{:else}
    <slot />
{/if}
