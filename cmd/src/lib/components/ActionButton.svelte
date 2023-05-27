<script lang="ts">
    import LoadingSpinner from "./LoadingSpinner.svelte";

    export let positive = true;
    export let loading = false;
</script>

<button
    class={`flex justify-center gap-2 py-1 px-2 text-white font-medium rounded-md  items-center ${
        positive
            ? "bg-[#3b82f6] hover:bg-blue-400"
            : "bg-red-500 hover:bg-red-400"
    } 
    ${$$props.class}`}
    on:click
>
    {#if loading}
        <LoadingSpinner class="w-4" />
    {/if}
    <slot />
</button>
