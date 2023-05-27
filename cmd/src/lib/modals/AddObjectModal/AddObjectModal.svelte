<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import AddData from "./AddData.svelte";
    import AddFolder from "./AddFolder.svelte";
    import Modal from "$lib/components/Modal/Modal.svelte";

    export let visible = false;
    let active: string = "data";

    const dispatch = createEventDispatcher();

    const hideModal = () => {
        visible = false;
        dispatch("getData");
    };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<Modal title="Add object" bind:visible class="w-96">
    <div class="flex gap-1 text-sm mt-2">
        <h1
            class={`hover:underline cursor-pointer ${
                active == "data" ? "underline" : ""
            }`}
            on:click={() => (active = "data")}
        >
            Data
        </h1>
        <h1
            class={`hover:underline cursor-pointer ${
                active == "folder" ? "underline" : ""
            }`}
            on:click={() => (active = "folder")}
        >
            Folder
        </h1>
    </div>
    {#if active == "data"}
        <AddData {active} on:hideModal={hideModal} />
    {:else}
        <AddFolder {active} on:hideModal={hideModal} />
    {/if}
</Modal>
