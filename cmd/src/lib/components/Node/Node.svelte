<script lang="ts">
    export let node: {
        ip: string;
        active: boolean;
    };
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    class={`w-full pb-1 [&:nth-child(n+2)]:pt-1 border-l-4 ${
        node.active ? "border-l-green-600" : "border-l-red-500"
    }`}
>
    <div class="flex pl-3 relative">
        <div class="w-20">
            <h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Node</h1>
        </div>
        <div class="w-60">
            <h1 class="font-medium">Ip</h1>
            <h1 class="folder_name font-normal text-lg">{node.ip}</h1>
        </div>
    </div>
</div>
