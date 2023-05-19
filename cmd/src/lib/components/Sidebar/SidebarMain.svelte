<script>
    import SidebarDivider from "./SidebarDivider.svelte";
    import SidebarTab from "./SidebarTab.svelte";
</script>

<div class="w-11/12 mx-auto">
    <SidebarDivider>Views</SidebarDivider>
    <SidebarTab href="/folders">Folders</SidebarTab>
    <SidebarTab href="/stats">Stats</SidebarTab>
    <SidebarDivider>Files</SidebarDivider>
    <SidebarTab href="/export">Export</SidebarTab>
    <SidebarTab href="/import">Import</SidebarTab>
    <SidebarDivider>Nodes</SidebarDivider>
    <SidebarTab href="/nodes">Nodes</SidebarTab>
    <SidebarDivider>Search</SidebarDivider>
    <SidebarTab href="/search">Advanced Search</SidebarTab>
</div>
