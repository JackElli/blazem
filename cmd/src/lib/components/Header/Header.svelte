<script lang="ts">
    import Tab from "../Tabs/Tab.svelte";
</script>

<div class="w-full h-12">
    <div class="flex justify-between items-center h-full">
        <div class="flex items-center">
            <a
                href="/folders"
                class="font-medium text-3xl text-[#3b82f6] cursor-pointer"
            >
                Blazem
            </a>

            <input
                class="ml-4 rounded-md w-52 text-sm px-1 pl-4 py-1 mt-[2px] focus:py-2 shadow-md bg-gray-100 focus:bg-white outline-none focus:w-64"
                placeholder="Search..."
            />
        </div>

        <Tab href="/help">Need help?</Tab>
    </div>
</div>
