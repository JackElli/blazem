<script lang="ts">
    export let title = "";
    export let visible = false;

    const setVisible = (e: any) => {
        if (e.target.className.includes("modal-position")) {
            visible = false;
        }
    };
</script>

{#if visible}
    <div
        class="modal-position w-full h-full fixed top-0 left-0 z-20 bg-black bg-opacity-40"
        on:mousedown={setVisible}
    >
        <div
            class="bg-white relative block mx-auto mt-8 rounded-md shadow-xl {$$props.class}"
        >
            <div class="p-5">
                <h1 class="font-medium text-lg">{title}</h1>
                <slot />
            </div>
        </div>
    </div>
{/if}
