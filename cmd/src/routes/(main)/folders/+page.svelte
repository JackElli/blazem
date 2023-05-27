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
    import ActionButton from "$lib/components/ActionButton.svelte";
    import PageTitle from "$lib/components/PageTitle.svelte";

    export let data;
    $: service = data.service;

    let folders: any;
    let addObjectVisible = false;

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
    <AddObjectModal on:getData={fetchData} bind:visible={addObjectVisible} />
    <div class="flex justify-between items-center">
        <PageTitle>/ Folders</PageTitle>
        <ActionButton on:click={() => (addObjectVisible = true)}>
            <p>Add object</p>
        </ActionButton>
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
        <p class="mt-4 text-3xl text-center">You currently have no folders</p>
        <button
            on:click={() => (addObjectVisible = true)}
            class="mt-2 text-xl block mx-auto text-center text-[#3b82f6] underline hover:text-blue-400"
            >Create one here</button
        >
    {/if}
</Loading>
