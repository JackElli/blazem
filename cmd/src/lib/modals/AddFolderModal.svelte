<script lang="ts">
    import ActionButton from "$lib/components/ActionButton.svelte";
    import Modal from "$lib/components/Modal/Modal.svelte";
    import { createEventDispatcher, tick } from "svelte";
    import { generateKey } from "$lib/funcs";
    import { hostName } from "../../global";
    import { networkRequest } from "$lib/network/request";

    export let visible = false;

    let folderValue: string;
    let folderNameTxt: HTMLInputElement;
    let privateFolder = true;

    const dispatch = createEventDispatcher();

    const addFolder = async () => {
        let folder = window.location.href.split("folder/")[1];
        let key = generateKey();

        if (folderValue == "") return;
        await networkRequest(`http://${hostName}:3100/folder`, {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                folder: folder,
                name: folderValue,
                docCount: 0,
                value: "",
                key: key,
                type: "folder",
                global: !privateFolder,
            }),
        });
        dispatch("hideModal");
        dispatch("getData");
    };

    $: visible && nowVisible();
    const nowVisible = async () => {
        if (visible) {
            await tick();
            folderNameTxt.focus();
        }
    };
</script>

<Modal title="Add Folder" bind:visible class="w-96">
    <p class="text-gray-300 text-xs mt-4">Folder name</p>
    <input
        class="border border-gray-300 rounded-sm w-80 h-7 pl-2"
        type="text"
        bind:value={folderValue}
        bind:this={folderNameTxt}
    />

    <div class="flex items-center gap-2 mt-2">
        Private
        <input type="checkbox" bind:checked={privateFolder} />
    </div>

    <ActionButton class="mt-4" on:click={() => addFolder()}>
        <p>Add</p>
    </ActionButton>
</Modal>
