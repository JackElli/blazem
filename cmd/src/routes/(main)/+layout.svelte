<script lang="ts">
    import "../../App.css";
    import { versionNum } from "../../global";
    import Header from "../../lib/components/Header/Header.svelte";
    import Sidebar from "../../lib/components/Sidebar/Sidebar.svelte";

    export let data;
</script>

<div class="flex gap-6 w-full">
    <Sidebar service={data.service} />
    <div class="w-10/12">
        <Header />
        <div class="mt-2">
            <slot />
        </div>
    </div>
    <p class="fixed bottom-2 right-2">v{versionNum}</p>
</div>
