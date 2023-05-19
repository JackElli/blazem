import { c as create_ssr_component, a as subscribe, v as validate_component, d as each, f as set_store_value } from './index2-e2d3d016.js';
import { n as needToFetchDataInFolder } from './stores2-fdb6e3fc.js';
import { h as hostName } from './global-6b632810.js';
import { F as Folder } from './Folder-164df24f.js';
import { P as Panel } from './Panel-37be2dea.js';
import { A as AddObjectModal } from './AddObjectModal-f61a07b7.js';
import { L as Loading } from './Loading-0ac87029.js';
import { n as networkRequest } from './request-2d8f7fbb.js';
import './index-fdb9b5cd.js';
import './spinner-75d0e167.js';

const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $needToFetchDataInFolder, $$unsubscribe_needToFetchDataInFolder;
  $$unsubscribe_needToFetchDataInFolder = subscribe(needToFetchDataInFolder, (value) => $needToFetchDataInFolder = value);
  let { data } = $$props;
  let folderResponse;
  let addObjectVisible = false;
  let loading = true;
  const fetchData = async () => {
    let resp = await networkRequest(`http://${hostName}:3100/folders`, { method: "GET", credentials: "include" });
    let data2 = await resp.data;
    folderResponse = data2?.data;
    loading = false;
  };
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  let $$settled;
  let $$rendered;
  do {
    $$settled = true;
    data.service;
    {
      {
        if ($needToFetchDataInFolder) {
          fetchData();
          set_store_value(needToFetchDataInFolder, $needToFetchDataInFolder = false, $needToFetchDataInFolder);
        }
      }
    }
    $$rendered = `${$$result.head += `<!-- HEAD_svelte-bikkq7_START -->${$$result.title = `<title>Blazem | Folders</title>`, ""}<!-- HEAD_svelte-bikkq7_END -->`, ""}

<div class="w-10/12">${validate_component(AddObjectModal, "AddObjectModal").$$render(
      $$result,
      { visible: addObjectVisible },
      {
        visible: ($$value) => {
          addObjectVisible = $$value;
          $$settled = false;
        }
      },
      {}
    )}
    <p class="font-medium text-gray-600 text-lg">/ Folders</p>
    <button class="flex justify-center mt-2 items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"><p class="ml-2 mr-2">Add object</p></button></div>

${validate_component(Loading, "Loading").$$render($$result, { loading }, {}, {
      default: () => {
        return `${validate_component(Panel, "Panel").$$render($$result, { class: "mt-4 bg-white" }, {}, {
          default: () => {
            return `${each(Object.entries(folderResponse), ([_, folder]) => {
              return `${validate_component(Folder, "Folder").$$render($$result, { data: folder }, {}, {})}`;
            })}`;
          }
        })}`;
      }
    })}`;
  } while (!$$settled);
  $$unsubscribe_needToFetchDataInFolder();
  return $$rendered;
});

export { Page as default };
//# sourceMappingURL=_page.svelte-1614dde0.js.map
