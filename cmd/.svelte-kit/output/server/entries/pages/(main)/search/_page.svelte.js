import { c as create_ssr_component, b as add_attribute, v as validate_component } from "../../../../chunks/index2.js";
import { D as DataContainer } from "../../../../chunks/DataContainer.js";
import { L as Loading } from "../../../../chunks/Loading.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let searchTxt;
  let allData;
  let loading = false;
  return `${$$result.head += `<!-- HEAD_svelte-122hss4_START -->${$$result.title = `<title>Blazem | Search</title>`, ""}<!-- HEAD_svelte-122hss4_END -->`, ""}
<div class="search_container"><textarea class="block mx-auto mt-5 border border-gray-300 font-medium h-20 p-2 resize-none w-full text-xl rounded-md shadow-md outline-none"${add_attribute("this", searchTxt, 0)}></textarea>
    <button class="flex justify-center items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative mt-2"><p class="ml-2 mr-2">Search</p></button>
    ${validate_component(Loading, "Loading").$$render($$result, { loading }, {}, {
    default: () => {
      return `${validate_component(DataContainer, "DataContainer").$$render($$result, { class: "mt-6", allData: allData?.docs }, {}, {})}`;
    }
  })}</div>`;
});
export {
  Page as default
};
