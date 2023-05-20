<script lang="ts">
    import { needToFetchDataInFolder } from "../../../lib/stores";
    import { hostName } from "../../../global";
    import { onMount } from "svelte";
    import type { NetworkResponse } from "$lib/types";
    import Folder from "../../../lib/components/Folder/Folder.svelte";
    import Panel from "$lib/components/Panel/Panel.svelte";
    import AddObjectModal from "$lib/modals/AddObjectModal/AddObjectModal.svelte";
    import Loading from "$lib/components/Loading.svelte";
    import { goto } from "$app/navigation";
    import { networkRequest } from "$lib/network/request";

    export let data;
    $: service = data.service;

    let folderResponse: NetworkResponse;
    let addObjectVisible = false;
    let loading = true;

    const fetchData = async () => {
        let resp = await networkRequest(`http://${hostName}:3100/folders`, {
            method: "GET",
            credentials: "include",
        });
        let data = await resp.data;
        folderResponse = data?.data;
        loading = false;
    };

    $: {
        if ($needToFetchDataInFolder) {
            fetchData();
            $needToFetchDataInFolder = false;
        }
    }

    onMount(async () => {
        if (service === false) {
            goto("/backup");
        }
        await fetchData();
    });
</script>

<svelte:head>
    <title>Blazem | Folders</title>
</svelte:head>

<div class="w-10/12">
    <AddObjectModal on:getData={fetchData} bind:visible={addObjectVisible} />
    <p class="font-medium text-gray-600 text-lg">/ Folders</p>
    <button
        class="flex justify-center mt-2 items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"
        on:click={() => (addObjectVisible = true)}
    >
        <p class="ml-2 mr-2">Add object</p>
    </button>
</div>

<Loading {loading}>
    {#if Object.keys(folderResponse).length != 0}
        <Panel class="mt-4 bg-white">
            {#each Object.entries(folderResponse) as [_, folder]}
                <Folder data={folder} />
            {/each}
        </Panel>
    {:else}
        <p class="mt-4">No folders available</p>
    {/if}
</Loading>
