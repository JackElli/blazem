<script>
    import ColourPanel from "$lib/components/Panel/ColourPanel.svelte";
    import PanelItem from "$lib/components/Panel/PanelItem.svelte";
</script>

<svelte:head>
    <title>Blazem | Stats</title>
</svelte:head>
<div class="flex gap-2 items-center h-8">
    <h1 class="font-medium text-gray-600 text-lg">Stats</h1>
    <p class="text-xs">(Not real for now)</p>
</div>

<div class="grid grid-cols-4 mt-4 gap-2">
    <div class="flex flex-col gap-2">
        <ColourPanel>
            <h1 class="font-medium">Number of files backed up</h1>
            <div class="mt-1">
                <PanelItem>831</PanelItem>
            </div>
        </ColourPanel>

        <ColourPanel>
            <h1 class="font-medium">Number of folders being watched</h1>
            <div class="mt-1">
                <PanelItem>412</PanelItem>
            </div>
        </ColourPanel>
        <ColourPanel>
            <h1 class="font-medium">Your billing</h1>
            <div class="mt-1">
                <PanelItem>$12</PanelItem>
            </div>
        </ColourPanel>
    </div>
    <div class="flex flex-col gap-2">
        <ColourPanel>
            <h1 class="font-medium">Most viewed files</h1>
            <div class="mt-2">
                <PanelItem>testing.txt</PanelItem>
                <PanelItem>testing.txt</PanelItem>
                <PanelItem>testing.txt</PanelItem>
                <PanelItem>testing.txt</PanelItem>
                <PanelItem>testing.txt</PanelItem>
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
