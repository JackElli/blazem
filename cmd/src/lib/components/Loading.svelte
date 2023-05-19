<script>
    import spinner from "../svg/spinner.svg";
    export let loading = false;
</script>

{#if loading}
    <img
        src={spinner}
        alt="spinner"
        class={`animate-spin w-10 ${$$props.class} fixed top-0 bottom-0 left-0 right-0 my-auto mx-auto`}
    />
{:else}
    <slot />
{/if}
