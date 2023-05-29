<script lang="ts">
    import { needToFetchDataInFolder } from "../../../lib/stores";
    import { hostName } from "../../../global";
    import { onMount } from "svelte";
    import Folder from "../../../lib/components/Folder/Folder.svelte";
    import Panel from "$lib/components/Panel/Panel.svelte";
    import Loading from "$lib/components/Loading.svelte";
    import { goto } from "$app/navigation";
    import { networkRequest } from "$lib/network/request";
    import ActionButton from "$lib/components/ActionButton.svelte";
    import PageTitle from "$lib/components/PageTitle.svelte";
    import AddFolderModal from "$lib/modals/AddFolderModal.svelte";

    export let data;
    $: service = data.service;

    let folders: any;
    let addFolderVisible = false;
    let loading = true;

    const fetchData = async () => {
        let resp = await networkRequest(`http://${hostName}:3100/folders`, {
            method: "GET",
            credentials: "include",
        });
        folders = resp.data;
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

<div>
    <AddFolderModal
        on:getData={fetchData}
        bind:visible={addFolderVisible}
        on:hideModal={() => (addFolderVisible = false)}
    />

    <div class="flex justify-between items-center">
        <PageTitle class="bg-gray-200 p-1 px-2 rounded-md">/ Folders</PageTitle>
        <div class="flex gap-2">
            <ActionButton on:click={() => (addFolderVisible = true)}>
                <p>New folder</p>
            </ActionButton>
        </div>
    </div>
</div>

<Loading {loading}>
    {#if Object.keys(folders).length != 0}
        <Panel class="mt-4 bg-white">
            {#each Object.entries(folders) as [_, folder]}
                <Folder data={folder} />
            {/each}
        </Panel>
    {:else}
        <p class="mt-7 text-3xl">You currently have no folders</p>
        <button
            on:click={() => (addFolderVisible = true)}
            class="mt-2 text-xl text-[#3b82f6] underline hover:text-blue-400"
            >Create one here</button
        >
    {/if}
</Loading>
