<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { hostName } from "../../global";
    import { networkRequest } from "$lib/network/request";

    export let visible = false;
    export let data: any;

    const dispatch = createEventDispatcher();

    const setVisible = (e: any) => {
        if (e.target.className.includes("modal_position")) {
            visible = false;
        }
    };

    const deleteDoc = async () => {
        await networkRequest(`http://${hostName}:3100/doc/${data.key}`, {
            method: "DELETE",
            credentials: "include",
        });
        visible = false;
        dispatch("getData");
    };
    $: dataAttrs = data.data;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
{#if visible}
    <div class="modal_position" on:mousedown={setVisible}>
        <div class="view_data_modal">
            <div class="modal_container">
                <h1 class="title">{data.key}</h1>
                <div class="data_main">
                    {#each Object.entries(dataAttrs) as [key, dataAttr]}
                        {#if key != "key"}
                            <div class="data_attr w-4/5">
                                <div class="data_attr_key">
                                    <p>{key}</p>
                                </div>
                                <p>{dataAttr}</p>
                            </div>
                        {/if}
                    {/each}
                </div>
                <button
                    class="flex absolute right-5 bottom-1 justify-center items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400"
                    on:click={() => deleteDoc()}
                >
                    <p class="ml-2 mr-2">Delete</p>
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    h1 {
        margin: 0;
    }
    .modal_position {
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.6);
        position: fixed;
        top: 0;
        left: 0;
        z-index: 10;
    }

    .modal_container {
        position: relative;
        padding: 20px;
    }

    .title {
        font-weight: 500;
    }

    .view_data_modal {
        width: 500px;
        background: white;
        display: block;
        margin: 0 auto;
        margin-top: 30px;
        border-radius: 3px;
    }

    .data_main {
        margin-top: 20px;
    }

    .data_attr {
        align-items: center;
        align-items: stretch;
        border-bottom: 1px solid #bbb;
        display: flex;
        flex-grow: 1;
        gap: 5px;
        padding-right: 1px;
    }
    .data_attr p {
        flex: 1 1;
        margin: 0;
        padding: 2px;
    }

    .data_attr:last-child {
        border-bottom: none;
    }

    .data_attr_key {
        align-items: center;
        color: #2c2c54;
        display: flex;
        font-weight: 500;
        padding-left: 3px;
        padding-right: 2px;
        width: 70px;
    }
</style>
