<script lang="ts">
    import { createEventDispatcher, tick } from "svelte";
    import { generateKey } from "$lib/funcs";
    import { hostName } from "../../../global";
    import { networkRequest } from "$lib/network/request";

    export let active: string;
    let folderName: any;

    const dispatch = createEventDispatcher();

    const addFolder = async () => {
        let folder = window.location.href.split("folder/")[1];
        let key = generateKey();
        let folderValue = folderName.value;

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
            }),
        });
        dispatch("hideModal");
        dispatch("getData");
    };

    $: active && activeChanged();
    const activeChanged = async () => {
        await tick();
        if (active == "folder") {
            folderName.focus();
            return;
        }
    };
</script>

<p class="text-gray-300 text-xs mt-2">Folder name</p>
<input
    class="border border-gray-300 rounded-sm w-80 h-7 pl-2"
    type="text"
    bind:this={folderName}
/>
<br />
<button
    class="flex justify-center items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative mt-2"
    on:click={() => addFolder()}
>
    <p class="ml-2 mr-2">Add</p>
</button>
