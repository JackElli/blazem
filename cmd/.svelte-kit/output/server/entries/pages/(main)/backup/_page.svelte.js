import { c as create_ssr_component, v as validate_component, b as add_attribute } from "../../../../chunks/index2.js";
import { P as Panel } from "../../../../chunks/Panel.js";
import { s as spinner } from "../../../../chunks/spinner.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let deploying;
  let { data } = $$props;
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  data.service;
  deploying = data.deploying;
  return `${$$result.head += `<!-- HEAD_svelte-1nqf2p6_START -->${$$result.title = `<title>Blazem | Backup</title>`, ""}<!-- HEAD_svelte-1nqf2p6_END -->`, ""}

<h1 class="mt-4 text-xl">Add a backup service</h1>

${!deploying ? `<h1 class="mt-4 text-md text-gray-500">Select a route</h1>
    <div class="flex gap-4 mt-4">${validate_component(Panel, "Panel").$$render(
    $$result,
    {
      class: `w-40 h-40 flex flex-col justify-center items-center hover:bg-gray-200 cursor-pointer ${"bg-gray-200"}`
    },
    {},
    {
      default: () => {
        return `Local`;
      }
    }
  )}
        ${validate_component(Panel, "Panel").$$render(
    $$result,
    {
      class: `w-40 h-40 flex flex-col justify-center items-center hover:bg-gray-200 cursor-pointer ${"bg-white"}`
    },
    {},
    {
      default: () => {
        return `AWS`;
      }
    }
  )}</div>

    ${`<div class="mt-4"><h1 class="text-md text-gray-500">Run these commands</h1>
            ${validate_component(Panel, "Panel").$$render($$result, { class: "mt-4 bg-white w-96" }, {}, {
    default: () => {
      return `<h1>docker pull blazem</h1>
                <h1>docker compose up</h1>`;
    }
  })}</div>`}
    ${``}` : `<img${add_attribute("src", spinner, 0)} alt="deploying" class="animate-spin">
    Backup Deploying`}`;
});
export {
  Page as default
};
