<script lang="ts">
    import { page } from "$app/stores";

    export let href: string;
</script>

<a
    {href}
    class={`font-medium px-2 py-1 rounded-md text-md cursor-pointer hover:bg-gray-100 ${
        href == $page.url.pathname ? "bg-gray-100" : ""
    }`}
>
    <slot /></a
>
