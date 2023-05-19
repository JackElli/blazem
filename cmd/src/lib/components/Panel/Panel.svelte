<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
    on:click
    class={`rounded-sm p-2 scroll-pb-44 shadow-md shadow-gray-400 ${$$props.class}`}
>
    <slot />
</div>
