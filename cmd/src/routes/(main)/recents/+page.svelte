<script lang="ts">
    import ColourPanel from "$lib/components/Panel/ColourPanel.svelte";
    import PanelItem from "$lib/components/Panel/PanelItem.svelte";
</script>

<svelte:head>
    <title>Blazem | Recents</title>
</svelte:head>
<div class="grid grid-cols-3 mt-5 gap-2">
    <ColourPanel>
        <h1 class="font-medium">Recently backed up files</h1>
        <div class="mt-2">
            <PanelItem>testing.txt</PanelItem>
            <PanelItem>testing.txt</PanelItem>
            <PanelItem>testing.txt</PanelItem>
            <PanelItem>testing.txt</PanelItem>
            <PanelItem>testing.txt</PanelItem>
        </div>
    </ColourPanel>
    <div class="flex flex-col gap-2">
        <ColourPanel>
            <h1 class="font-medium">Recently viewed files</h1>
            <div class="mt-2">
                <PanelItem>darthdad.mp4</PanelItem>
                <PanelItem>hello.jpg</PanelItem>
            </div>
        </ColourPanel>
        <ColourPanel>
            <h1 class="font-medium">Recently deleted files</h1>
            <div class="mt-2">
                <PanelItem>passwds.txt</PanelItem>
                <PanelItem>amazing.txt</PanelItem>
            </div>
        </ColourPanel>
    </div>
    <ColourPanel>
        <h1 class="font-medium">Recent server logs</h1>
        <div class="mt-2">
            <PanelItem>01/05/23 Good!</PanelItem>
            <PanelItem>01/05/23 Good!</PanelItem>
            <PanelItem>01/05/23 Bad!</PanelItem>
            <PanelItem>01/05/23 Bad!</PanelItem>
            <PanelItem>01/05/23 Bad!</PanelItem>
        </div>
    </ColourPanel>
</div>
