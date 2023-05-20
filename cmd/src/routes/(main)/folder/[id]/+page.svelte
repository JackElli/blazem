<script lang="ts">
    import { onDestroy, onMount } from "svelte";

    import type { NetworkResponse } from "$lib/types";
    import Loading from "$lib/components/Loading.svelte";
    import { networkRequest } from "$lib/network/request";
    import { hostName } from "../../../../global.js";
    import { needToFetchDataInFolder, slugData } from "$lib/stores.js";
    import DataContainer from "$lib/components/DataContainer/DataContainer.svelte";

    type DataInFolder = {
        data: object;
        folderName: string;
    };

    interface DataInFolderResponse extends NetworkResponse {
        data: DataInFolder;
    }
    interface ParentFolderResponse extends NetworkResponse {
        data: object[];
    }

    export let data;

    let folderId = data?.folder.id;
    let loading = true;
    let allData: NetworkResponse;
    let parentFolders: NetworkResponse;

    $: folderName = allData?.folderName;

    const getFolderData = async () => {
        let folderResp = await networkRequest(
            `http://${hostName}:3100/folder/${folderId}`,
            {
                method: "GET",
                credentials: "include",
            }
        );
        let folderData = await folderResp.data;
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
<Loading {loading}>
    <DataContainer on:getData={fetchData} allData={allData.data} />
</Loading>
