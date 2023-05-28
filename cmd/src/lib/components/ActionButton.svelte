<script lang="ts">
    import LoadingSpinner from "./LoadingSpinner.svelte";
    export let loading = false;
    export let colour = "bg-blue-500 hover:bg-blue-400";
</script>

<button
    class="flex justify-center gap-2 py-1 px-2 text-white font-medium rounded-md items-center {$$props.class} {colour}"
    on:click
>
    {#if loading}
        <LoadingSpinner class="w-4" />
    {/if}
    <slot />
</button>
