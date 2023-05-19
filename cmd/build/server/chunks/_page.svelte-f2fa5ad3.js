import { c as create_ssr_component, b as add_attribute, v as validate_component } from './index2-e2d3d016.js';
import { D as DataContainer } from './DataContainer-102101e7.js';
import { L as Loading } from './Loading-0ac87029.js';
import './Folder-164df24f.js';
import './Panel-37be2dea.js';
import './spinner-75d0e167.js';

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

export { Page as default };
//# sourceMappingURL=_page.svelte-f2fa5ad3.js.map
