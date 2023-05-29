<script lang="ts">
    import { onDestroy, onMount } from "svelte";
    import Loading from "$lib/components/Loading.svelte";
    import { networkRequest } from "$lib/network/request";
    import { hostName } from "../../../../global.js";
    import { needToFetchDataInFolder, slugData } from "$lib/stores.js";
    import DataContainer from "$lib/components/DataContainer/DataContainer.svelte";
    import Breadcrumb from "$lib/components/Breadcrumb/Breadcrumb.svelte";
    import ActionButton from "$lib/components/ActionButton.svelte";
    import AddFolderModal from "$lib/modals/AddFolderModal.svelte";
    import AddDataModal from "$lib/modals/AddDataModal.svelte";

    export let data;

    let allData: any;
    let parentFolders: any;
    let folderId = data?.folder.id;
    let loading = true;
    let addFolderVisible = false;
    let addDataVisible = false;
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
            <div class="flex gap-2">
                <ActionButton on:click={() => (addFolderVisible = true)}>
                    <p>New folder</p>
                </ActionButton>
                <ActionButton on:click={() => (addDataVisible = true)}>
                    <p>Add data</p>
                </ActionButton>
            </div>
        </div>
        <AddFolderModal
            on:getData={getData}
            bind:visible={addFolderVisible}
            on:hideModal={() => (addFolderVisible = false)}
        />
        <AddDataModal
            on:getData={getData}
            bind:visible={addDataVisible}
            on:hideModal={() => (addDataVisible = false)}
        />
    {/if}
</div>

<div class="mt-4">
    {#if !unauthorised}
        <Loading {loading}>
            {#if Object.keys(allData?.data ?? {}).length != 0}
                <DataContainer on:getData={fetchData} allData={allData.data} />
            {:else}
                <p class="mt-7 text-3xl">
                    You currently have no data in this folder
                </p>
                <button
                    on:click={() => (addDataVisible = true)}
                    class="mt-2 text-xl text-[#3b82f6] underline hover:text-blue-400"
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
