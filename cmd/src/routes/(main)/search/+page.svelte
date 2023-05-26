<script lang="ts">
    import { networkRequest } from "$lib/network/request";
    import type { NetworkResponse } from "$lib/types";
    import { onMount } from "svelte";
    import { hostName } from "../../../global";
    import DataContainer from "$lib/components/DataContainer/DataContainer.svelte";
    import Loading from "$lib/components/Loading.svelte";
    import ActionButton from "$lib/components/ActionButton.svelte";
    import PageTitle from "$lib/components/PageTitle.svelte";

    const specialWords: string[] = ["select", "delete", "where"];

    type QueryResult = {
        docs: object[];
    };

    let searchTxt: HTMLTextAreaElement;
    let allData: any;
    let loading = false;

    const search = async () => {
        const searchValue = searchTxt.value;
        loading = true;

        const queryResp = await networkRequest(
            `http://${hostName}:3100/query`,
            {
                method: "POST",
                credentials: "include",
                body: JSON.stringify({ query: searchValue }),
            }
        );
        allData = queryResp?.data;
        loading = false;
    };

    const checkForSpecial = (e: KeyboardEvent) => {
        if (e.key == "Enter") {
            e.preventDefault();
            search();
            return;
        }
        changeSpecialWords();
    };

    const changeSpecialWords = () => {
        let searchVal = searchTxt.value.split(" ");
        searchVal = searchVal.map((word: string) => {
            if (specialWords.includes(word)) {
                return word.toUpperCase();
            }
            return word;
        });
        searchTxt.value = searchVal.join(" ");
    };

    onMount(() => {
        searchTxt.focus();
    });
</script>

<svelte:head>
    <title>Blazem | Search</title>
</svelte:head>
<div>
    <PageTitle>Advanced Search</PageTitle>
    <textarea
        class="block mx-auto mt-5 border border-gray-300 font-medium h-20 p-2 resize-none w-full text-xl rounded-sm shadow-md outline-none"
        placeholder="SELECT all WHERE..."
        bind:this={searchTxt}
        on:keydown={checkForSpecial}
    />
    <ActionButton class="mt-2" on:click={() => search()}>
        <p class="ml-2 mr-2">Search</p>
    </ActionButton>

    <Loading {loading}>
        <DataContainer class="mt-6" allData={allData?.docs} />
    </Loading>
</div>
