<script>
    import { slugData } from "$lib/stores";
    import spinner from "$lib/svg/spinner.svg";

    $: data = $slugData;
</script>

{#if $slugData.current && $slugData.defaultVal && $slugData.current != "" && $slugData.defaultVal != ""}
    <div class="flex gap-1 {$$props.class} bg-gray-200 p-1 px-2 rounded-md">
        {#if data.defaultVal != data.current}
            <a class="text-lg underline text-[#3d3d75]" href="/folders"
                >{data.defaultVal}</a
            >
            <p class="font-medium text-gray-600 text-lg">/</p>
        {/if}

        {#if data.previous.length > 0}
            {#each data.previous as previousFolder}
                <a
                    class="text-lg text-[#3d3d75] underline"
                    href={`/folder/${previousFolder.key}`}
                    >{previousFolder.name}</a
                >
                <p class=" text-gray-600 text-lg">/</p>
            {/each}
        {/if}
        <p class="font-medium text-gray-600 text-lg">{data.current}</p>
    </div>
{:else}
    <img src={spinner} alt="spinner" class="animate-spin w-7" />
{/if}
