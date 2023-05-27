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
    import ActionButton from "$lib/components/ActionButton.svelte";

    export let data;

    let allData: any;
    let parentFolders: any;
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
        if (folderResp.code == 403) {
            unauthorised = true;
            return;
        }
        allData = folderResp?.data;
    };

    const getBreadcrumbData = async () => {
        let parentResp = await networkRequest(
            `http://${hostName}:3100/parents/${folderId}`,
            {
                method: "GET",
                credentials: "include",
            }
        );
        parentFolders = parentResp?.data;
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
        $slugData.defaultVal = "Folders";
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

<div class=" z-10">
    {#if !unauthorised}
        <div class="flex justify-between items-center w-full">
            <Breadcrumb />
            <ActionButton on:click={() => (addObjectVisible = true)}>
                <p>Add object</p>
            </ActionButton>
        </div>
        <AddObjectModal on:getData={getData} bind:visible={addObjectVisible} />
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
                    class="mt-2 text-xl block mx-auto text-center text-[#3b82f6] underline hover:text-blue-400"
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
