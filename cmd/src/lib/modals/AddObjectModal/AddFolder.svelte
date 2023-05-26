<script lang="ts">
    import { createEventDispatcher, tick } from "svelte";
    import { generateKey } from "$lib/funcs";
    import { hostName } from "../../../global";
    import { networkRequest } from "$lib/network/request";
    import ActionButton from "$lib/components/ActionButton.svelte";

    export let active: string;
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

    $: active && activeChanged();
    const activeChanged = async () => {
        await tick();
        if (active == "folder") {
            folderNameTxt.focus();
            return;
        }
    };
</script>

<p class="text-gray-300 text-xs mt-2">Folder name</p>
<input
    class="border border-gray-300 rounded-sm w-80 h-7 pl-2"
    type="text"
    bind:value={folderValue}
    bind:this={folderNameTxt}
/>
<br />
<div class="flex items-center gap-2 mt-5">
    Private
    <input type="checkbox" bind:checked={privateFolder} />
</div>

<ActionButton class="mt-2" on:click={() => addFolder()}>
    <p class="ml-2 mr-2">Add</p>
</ActionButton>
