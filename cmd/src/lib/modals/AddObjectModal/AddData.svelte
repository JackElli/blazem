<script lang="ts">
    import { createEventDispatcher, tick } from "svelte";
    import { generateKey } from "$lib/funcs";
    import { hostName } from "../../../global";
    import { networkRequest } from "$lib/network/request";
    import ActionButton from "$lib/components/ActionButton.svelte";

    export let active: string;
    let dataKey: any;
    let dataValue: any;

    const dispatch = createEventDispatcher();

    const addData = async () => {
        let folder = window.location.href.split("folder/")[1];
        let key = dataKey.value || dataKey.placeholder;
        let type = "text";
        let value = dataValue.value;

        if (value == "") return;

        await networkRequest(`http://${hostName}:3100/doc`, {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                folder: folder,
                key: key,
                type: type,
                value: value,
            }),
        });
        dispatch("hideModal");
        dispatch("getData");
    };

    $: active && activeChanged();
    const activeChanged = async () => {
        await tick();
        if (active == "data") {
            dataKey.placeholder = generateKey();
            dataValue.focus();
        }
    };
</script>

<p class="text-gray-300 text-xs mt-2">Document Key</p>
<input
    class="border border-gray-300 rounded-sm w-80 h-7 pl-2"
    type="text"
    placeholder="testkey"
    bind:this={dataKey}
/>
<br />
<p class="text-gray-300 text-xs mt-2">Document value</p>
<textarea
    class="border border-gray-300 h-24 pl-2 pt-2 resize-none w-80"
    bind:this={dataValue}
/>
<p class="text-[#3d3d75] hover:underline cursor-pointer">Upload a file?</p>
<br />
<ActionButton on:click={() => addData()}>
    <p>Add</p>
</ActionButton>
