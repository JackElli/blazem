<script>
    import { goto } from "$app/navigation";
    import Panel from "$lib/components/Panel/Panel.svelte";
    import { onMount } from "svelte";
    import spinner from "$lib/svg/spinner.svg";

    export let data;

    let selectedRoute = "local";

    $: service = data.service;
    $: deploying = data.deploying;

    onMount(async () => {
        if (service === true) {
            goto("/folders");
        }
    });
</script>

<svelte:head>
    <title>Blazem | Backup</title>
</svelte:head>

<h1 class="mt-4 text-xl">Add a backup service</h1>

{#if !deploying}
    <h1 class="mt-4 text-md text-gray-500">Select a route</h1>
    <div class="flex gap-4 mt-4">
        <Panel
            on:click={() => (selectedRoute = "local")}
            class={`w-40 h-40 flex flex-col justify-center items-center hover:bg-gray-200 cursor-pointer ${
                selectedRoute == "local" ? "bg-gray-200" : "bg-white"
            }`}>Local</Panel
        >
        <Panel
            on:click={() => (selectedRoute = "aws")}
            class={`w-40 h-40 flex flex-col justify-center items-center hover:bg-gray-200 cursor-pointer ${
                selectedRoute == "aws" ? "bg-gray-200" : "bg-white"
            }`}>AWS</Panel
        >
    </div>

    {#if selectedRoute == "local"}
        <div class="mt-4">
            <h1 class=" text-md text-gray-500">Run these commands</h1>
            <Panel class="mt-4 bg-white w-96">
                <h1>docker pull blazem</h1>
                <h1>docker compose up</h1>
            </Panel>
        </div>
    {/if}
    {#if selectedRoute == "aws"}
        <div class="mt-4">
            <button class="bg-white p-2">Deploy</button>
        </div>
    {/if}
{:else}
    <img src={spinner} alt="deploying" class="animate-spin" />
    Backup Deploying
{/if}
