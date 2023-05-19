<script lang="ts">
    import ViewDataModal from "$lib/modals/ViewDataModal.svelte";
    export let data: any;

    let viewDataVisible = false;

    const clamp = (value: string) => {
        if (value.length > 10) {
            return value.substring(0, 10) + "...";
        }
        return value;
    };

    const valueClamp = (value: unknown) => {
        let newValue = value as string;
        if (newValue.length > 60) {
            return newValue.substring(0, 60) + "...";
        }
        return newValue;
    };

    $: dataAttrs = data.data;
</script>

<ViewDataModal {data} on:getData bind:visible={viewDataVisible} />
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    class="data w-full pb-1 border-l-4 border-l-blue-500 cursor-pointer"
    on:click={() => {
        viewDataVisible = true;
    }}
>
    <div class="flex pl-3 relative h-full">
        <div class="min-w-[80px]">
            <h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Data</h1>
        </div>

        <div class="min-w-[120px]">
            <h1 class="font-medium">Key</h1>
            <h1 class="folder_key font-normal text-lg">{clamp(data.key)}</h1>
        </div>

        <div class="flex gap-20">
            {#if dataAttrs.type == "text"}
                {#each Object.entries(dataAttrs) as [key, dataAttr]}
                    {#if !["key", "folder", "type"].includes(key)}
                        <div>
                            <h1 class="font-medium">{key}</h1>
                            <h1 class="font-normal text-lg">
                                {valueClamp(dataAttr)}
                            </h1>
                        </div>
                    {/if}
                {/each}
            {:else}
                <!-- svelte-ignore a11y-missing-attribute -->
                <img src={dataAttrs.value} />
            {/if}
        </div>
    </div>
</div>

<style>
    .data:hover .folder_key {
        text-decoration: underline;
    }
</style>
