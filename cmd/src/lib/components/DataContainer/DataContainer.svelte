<script lang="ts">
    import Data from "../Data/Data.svelte";
    import Folder from "../Folder/Folder.svelte";
    import Panel from "../Panel/Panel.svelte";

    export let allData: any;

    let folders: [{ [id: string]: any }];
    let datas: [{ [id: string]: any }];

    $: allData && load();
    const load = () => {
        folders = allData.filter(function (item: any) {
            return item.data.type == "folder";
        });

        datas = allData.filter(function (item: any) {
            return item.data.type != "folder";
        });
    };
</script>

{#if allData}
    {#if allData.length > 0}
        <div class={`flex flex-col gap-4 w-full mb-3 ${$$props.class}`}>
            {#if folders.length > 0}
                <Panel class="bg-white">
                    {#each folders as folder}
                        <Folder data={folder.data} />
                    {/each}
                </Panel>
            {/if}
            {#if datas.length > 0}
                <Panel class="bg-white">
                    {#each datas as data}
                        <Data on:getData {data} />
                    {/each}
                </Panel>
            {/if}
        </div>
    {:else}
        <p class="mt-2">No data found</p>
    {/if}
{/if}
