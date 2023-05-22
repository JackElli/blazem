<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import type { NetworkResponse } from "$lib/types";
    import Loading from "$lib/components/Loading.svelte";
    import { networkRequest } from "$lib/network/request";
    import { hostName } from "../../../../global.js";
    import { needToFetchDataInFolder, slugData } from "$lib/stores.js";
    import DataContainer from "$lib/components/DataContainer/DataContainer.svelte";
    import Breadcrumb from "$lib/components/Breadcrumb/Breadcrumb.svelte";
    import AddObjectModal from "$lib/modals/AddObjectModal/AddObjectModal.svelte";

    export let data;

    let allData: NetworkResponse;
    let parentFolders: NetworkResponse;
    let folderId = data?.folder.id;
    let loading = true;
    let addObjectVisible = false;
    let unauthorised = false;

    const getData = () => {
        $needToFetchDataInFolder = true;
    };

    const getFolderData = async () => {
        let folderResp = await networkRequest(
            `http://${hostName}:3100/folder/${folderId}`,
            {
                method: "GET",
                credentials: "include",
            }
        );
        let folderData = await folderResp.data;
        if (folderData?.code == 403) {
            unauthorised = true;
            return;
        }
        allData = folderData?.data;
    };

    const getBreadcrumbData = async () => {
        let parentResp = await networkRequest(
            `http://${hostName}:3100/parents/${folderId}`,
            {
                method: "GET",
                credentials: "include",
            }
        );
        let parentData = await parentResp.data;
        parentFolders = parentData?.data;
    };

    const fetchData = async () => {
        await getFolderData();
        await getBreadcrumbData();

        loading = false;
    };

    $: data && pageChange();
    const pageChange = async () => {
        folderId = data?.folder.id;
        await fetchData();
    };

    $: folderName = allData?.folderName;
    $: parentFolders && parentChange();
    const parentChange = () => {
        $slugData.current = folderName;
        $slugData.previous = parentFolders ?? [];
    };

    $: {
        if ($needToFetchDataInFolder) {
            fetchData();
            $needToFetchDataInFolder = false;
        }
    }

    onMount(async () => {
        $slugData.defaultVal = "Your folders";
        $slugData.current = folderName;
    });

    onDestroy(() => {
        $slugData.defaultVal = "";
        $slugData.current = "";
        $slugData.previous = [];
    });
</script>

<!-- svelte-ignore missing-declaration -->
<svelte:head>
    <title>Blazem | {folderName ?? ""}</title>
</svelte:head>

<div class="w-10/12 z-10">
    {#if !unauthorised}
        <Breadcrumb />
        <AddObjectModal on:getData={getData} bind:visible={addObjectVisible} />
        <button
            class="flex justify-center mt-2 items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"
            on:click={() => (addObjectVisible = true)}
        >
            <p class="ml-2 mr-2">Add object</p>
        </button>
    {/if}
</div>

<div class="mt-4">
    {#if !unauthorised}
        <Loading {loading}>
            {#if Object.keys(allData?.data ?? {}).length != 0}
                <DataContainer on:getData={fetchData} allData={allData.data} />
            {:else}
                <p class="mt-4 text-3xl text-center">
                    You currently have no data in this folder
                </p>
                <button
                    on:click={() => (addObjectVisible = true)}
                    class="mt-2 text-xl block mx-auto text-center text-[#3b82f6] underline"
                    >Create some here</button
                >
            {/if}
        </Loading>
    {:else}
        <p class="mt-20 text-3xl text-center">
            You do not have permission to view this folder.
        </p>
    {/if}
</div>
