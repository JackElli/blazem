<script lang="ts">
    type Node = {
        ip: string;
        active: boolean;
    };
    export let node: Node;
</script>

<div
    class={`pt-1 pb-1 border-b border-b-gray-300 bg-white border-l-4 ${
        node.active == true ? "border-l-green-400" : "border-l-red-500"
    }`}
>
    <div class="flex gap-10 pl-3 relative h-full">
        <div class="w-32">
            <h1 class="font-medium">IP</h1>
            <h1 class="font-normal text-lg">{node.ip}</h1>
        </div>

        <div>
            <h1 class="font-medium">CPU</h1>
            <h1 class="font-normal text-lg">??</h1>
        </div>

        <div>
            <h1 class="font-medium">RAM</h1>
            <h1 class="font-normal text-lg">??</h1>
        </div>

        <div class="flex gap-20" />
    </div>
</div>
