<script lang="ts">
    import { page } from "$app/stores";

    export let href: any;
</script>

<button
    class="block w-full text-left rounded-md py-1 pl-2 mt-1 hover:bg-gray-300 cursor-pointer {$page
        .url.pathname == href
        ? 'bg-gray-300 font-medium'
        : ''}"
    on:click={() => {
        window.location.assign(href);
    }}
>
    <slot />
</button>
