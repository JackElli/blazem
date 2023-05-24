<script lang="ts">
    import Panel from "$lib/components/Panel/Panel.svelte";
    import Node from "$lib/components/Node/Node.svelte";
    import Loading from "$lib/components/Loading.svelte";
    import { networkRequest } from "$lib/network/request";
    import type { NetworkResponse } from "$lib/types";
    import { onMount } from "svelte";
    import { hostName } from "../../../global";

    type Node = {
        ip: string;
        active: boolean;
    };

    let nodes: Node[] = [];
    let loading = true;

    onMount(async () => {
        const nodeResp = await networkRequest(
            `http://${hostName}:3100/nodemap`,
            {
                method: "GET",
                credentials: "include",
            }
        );
        nodes = nodeResp?.data;
        loading = false;
    });
</script>

<svelte:head>
    <title>Blazem | Nodes</title>
</svelte:head>

<Loading {loading}>
    <div class="mt-5">
        <Panel class="bg-white">
            {#each nodes as node}
                <Node {node} />
            {/each}
        </Panel>
    </div>
</Loading>
