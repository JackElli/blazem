import { c as create_ssr_component, e as escape, j as null_to_empty } from './index2-e2d3d016.js';

const css = {
  code: ".folder.svelte-4gck9d:hover{cursor:pointer}",
  map: null
};
const Folder = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { data } = $$props;
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  $$result.css.add(css);
  return `
<div class="${escape(
    null_to_empty(`folder w-full cursor-pointer pb-1 [&:nth-child(n+2)]:pt-1  border-l-4 ${data.backedUp ? "border-l-green-600" : "border-l-[#3d3d75]"}`),
    true
  ) + " svelte-4gck9d"}"><div class="flex pl-3 relative group"><div class="w-20"><h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Folder</h1></div>
        <div class="w-60"><h1 class="font-medium">Folder name</h1>
            <h1 class="folder_name font-normal text-lg group-hover:underline">${escape(data.folderName)}</h1></div>

        <div class="w-40"><h1 class="font-medium">Doc count</h1>
            <h1 class="font-normal text-lg">${escape(data.docCount ?? "?")}</h1></div>

        <div><h1 class="font-medium">Permission</h1>
            <h1 class="font-normal text-lg">Open</h1></div></div>
</div>`;
});

export { Folder as F };
//# sourceMappingURL=Folder-164df24f.js.map
