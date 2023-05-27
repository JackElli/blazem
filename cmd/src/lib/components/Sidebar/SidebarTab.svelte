<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";

    export let href: any;
</script>

<button
    class="block w-full text-left mt-2 rounded-md py-1 pl-2 hover:bg-gray-300 cursor-default {$page
        .url.pathname == href
        ? 'bg-gray-300 font-medium shadow-lg'
        : ''}"
    on:click={() => {
        goto(href);
    }}
>
    <slot />
</button>
