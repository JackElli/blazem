<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { hostName } from "../../global";
    import { networkRequest } from "$lib/network/request";
    import ActionButton from "$lib/components/ActionButton.svelte";
    import Modal from "$lib/components/Modal/Modal.svelte";

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

<Modal bind:visible title={data.key} class="w-[400px]">
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
    <ActionButton
        positive={false}
        class="absolute right-5 bottom-1"
        on:click={() => deleteDoc()}
    >
        <p>Delete</p>
    </ActionButton>
</Modal>

<style>
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
