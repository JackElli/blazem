<script lang="ts">
    import SidebarMain from "./SidebarMain.svelte";
    import SidebarNoService from "./SidebarNoService.svelte";

    export let service = false;
</script>

<div class="bg-gray-200 w-40 h-screen border-r border-r-gray-300 sticky top-0">
    {#if service}
        <SidebarMain />
    {:else}
        <SidebarNoService />
    {/if}
    <p
        class="absolute bottom-0 left-0 right-0 text-center text-xs text-gray-400"
    >
        Jack Ellis 2023
    </p>
</div>
