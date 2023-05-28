<script>
    import PageTitle from "$lib/components/PageTitle.svelte";
    import ColourPanel from "$lib/components/Panel/ColourPanel.svelte";
    import PanelItem from "$lib/components/Panel/PanelItem.svelte";
</script>

<svelte:head>
    <title>Blazem | Stats</title>
</svelte:head>
<div class="flex items-center">
    <PageTitle>Stats</PageTitle>
</div>

<div class="grid grid-cols-4 mt-4 gap-2">
    <div class="flex flex-col gap-2">
        <ColourPanel>
            <h1 class="font-medium">Number of folders</h1>
            <div class="mt-1">
                <PanelItem>2</PanelItem>
            </div>
        </ColourPanel>

        <ColourPanel>
            <h1 class="font-medium">Number of files</h1>
            <div class="mt-1">
                <PanelItem>12</PanelItem>
            </div>
        </ColourPanel>
        <ColourPanel>
            <h1 class="font-medium">Number of JSON documents</h1>
            <div class="mt-1">
                <PanelItem>15</PanelItem>
            </div>
        </ColourPanel>
    </div>
    <div class="flex flex-col gap-2 col-span-2">
        <ColourPanel>
            <h1 class="font-medium">CPU usage</h1>
            <div class="mt-2">
                <PanelItem>Imagine a graph of CPU usage</PanelItem>
                <div class="h-40" />
            </div>
        </ColourPanel>
    </div>
    <div class="flex flex-col gap-2">
        <ColourPanel>
            <h1 class="font-medium">Blazem uptime</h1>
            <div class="mt-1">
                <PanelItem>1043 days</PanelItem>
            </div>
        </ColourPanel>
    </div>
</div>
