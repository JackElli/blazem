import { c as create_ssr_component, a as subscribe, v as validate_component, d as each, f as set_store_value } from "../../../../chunks/index2.js";
import { n as needToFetchDataInFolder } from "../../../../chunks/stores2.js";
import { h as hostName } from "../../../../chunks/global.js";
import { F as Folder } from "../../../../chunks/Folder.js";
import { P as Panel } from "../../../../chunks/Panel.js";
import { A as AddObjectModal } from "../../../../chunks/AddObjectModal.js";
import { L as Loading } from "../../../../chunks/Loading.js";
import { n as networkRequest } from "../../../../chunks/request.js";
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
export {
  Page as default
};
