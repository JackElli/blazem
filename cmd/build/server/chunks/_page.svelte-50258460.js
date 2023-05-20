import { c as create_ssr_component, v as validate_component, d as each, b as add_attribute, e as escape } from './index2-e2d3d016.js';
import { P as Panel } from './Panel-37be2dea.js';
import { L as Loading } from './Loading-0ac87029.js';
import './spinner-75d0e167.js';

const Node = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { node } = $$props;
  if ($$props.node === void 0 && $$bindings.node && node !== void 0)
    $$bindings.node(node);
  return `
<div${add_attribute("class", ` w-full cursor-pointer pb-1 [&:nth-child(n+2)]:pt-1 border-l-4 ${node.active ? "border-l-green-600" : "border-l-red-500"}`, 0)}><div class="flex pl-3 relative"><div class="w-20"><h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Node</h1></div>
        <div class="w-60"><h1 class="font-medium">Ip</h1>
            <h1 class="folder_name font-normal text-lg">${escape(node.ip)}</h1></div></div></div>`;
});
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let nodes = [];
  let loading = true;
  return `${$$result.head += `<!-- HEAD_svelte-d22l93_START -->${$$result.title = `<title>Blazem | Nodes</title>`, ""}<!-- HEAD_svelte-d22l93_END -->`, ""}

${validate_component(Loading, "Loading").$$render($$result, { loading }, {}, {
    default: () => {
      return `<div class="mt-5">${validate_component(Panel, "Panel").$$render($$result, { class: "bg-white" }, {}, {
        default: () => {
          return `${each(nodes, (node) => {
            return `${validate_component(Node, "Node").$$render($$result, { node }, {}, {})}`;
          })}`;
        }
      })}</div>`;
    }
  })}`;
});

export { Page as default };
//# sourceMappingURL=_page.svelte-50258460.js.map
