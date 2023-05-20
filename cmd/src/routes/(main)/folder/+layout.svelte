<script>
    import Breadcrumb from "$lib/components/Breadcrumb/Breadcrumb.svelte";
    import AddObjectModal from "$lib/modals/AddObjectModal/AddObjectModal.svelte";
    import { needToFetchDataInFolder } from "$lib/stores";

    let addObjectVisible = false;
    const getData = () => {
        $needToFetchDataInFolder = true;
    };
</script>

<div class="w-10/12 z-10">
    <Breadcrumb />
    <AddObjectModal on:getData={getData} bind:visible={addObjectVisible} />
    <button
        class="flex justify-center mt-2 items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"
        on:click={() => (addObjectVisible = true)}
    >
        <p class="ml-2 mr-2">Add object</p>
    </button>
</div>
<div class="mt-4">
    <slot />
</div>
