<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import AddData from "./AddData.svelte";
    import AddFolder from "./AddFolder.svelte";

    export let visible = false;

    let active: string = "data";

    const dispatch = createEventDispatcher();

    const setVisible = (e: any) => {
        if (e.target.className.includes("modal-position")) {
            visible = false;
        }
    };

    const hideModal = () => {
        visible = false;
        dispatch("getData");
    };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
{#if visible}
    <div
        class="modal-position w-full h-full fixed top-0 left-0 z-20 bg-black bg-opacity-40"
        on:mousedown={setVisible}
    >
        <div class="w-[400px] bg-white block mx-auto mt-8 rounded-sm">
            <div class="p-5">
                <h1 class="font-medium">Add object</h1>
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
            </div>
        </div>
    </div>
{/if}
