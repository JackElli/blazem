<script lang="ts">
    export let data: {
        key: string;
        folderName: string;
        docCount: number;
        backedUp?: string;
    };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    class={`folder w-full cursor-pointer pb-1 [&:nth-child(n+2)]:pt-1  border-l-4 ${
        data.backedUp ? "border-l-green-600" : "border-l-[#3d3d75]"
    }`}
    on:click={() => window.location.assign("/folder/" + data.key)}
>
    <div class="flex pl-3 relative group">
        <div class="w-20">
            <h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Folder</h1>
        </div>
        <div class="w-60">
            <h1 class="font-medium">Folder name</h1>
            <h1 class="folder_name font-normal text-lg group-hover:underline">
                {data.folderName}
            </h1>
        </div>

        <div class="w-40">
            <h1 class="font-medium">Doc count</h1>
            <h1 class="font-normal text-lg">{data.docCount ?? "?"}</h1>
        </div>

        <div>
            <h1 class="font-medium">Permission</h1>
            <h1 class="font-normal text-lg">Open</h1>
        </div>
    </div>
</div>

<style>
    .folder:hover {
        cursor: pointer;
    }
</style>
