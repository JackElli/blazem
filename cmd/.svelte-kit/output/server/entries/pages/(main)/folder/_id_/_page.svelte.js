import { c as create_ssr_component, a as subscribe, o as onDestroy, e as escape, v as validate_component, f as set_store_value } from "../../../../../chunks/index2.js";
import { L as Loading } from "../../../../../chunks/Loading.js";
import { n as networkRequest } from "../../../../../chunks/request.js";
import { h as hostName } from "../../../../../chunks/global.js";
import { s as slugData, n as needToFetchDataInFolder } from "../../../../../chunks/stores2.js";
import { D as DataContainer } from "../../../../../chunks/DataContainer.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let folderName;
  let $slugData, $$unsubscribe_slugData;
  let $needToFetchDataInFolder, $$unsubscribe_needToFetchDataInFolder;
  $$unsubscribe_slugData = subscribe(slugData, (value) => $slugData = value);
  $$unsubscribe_needToFetchDataInFolder = subscribe(needToFetchDataInFolder, (value) => $needToFetchDataInFolder = value);
  let { data } = $$props;
  let folderId = data?.folder.id;
  let loading = true;
  let allData;
  let parentFolders;
  const getFolderData = async () => {
    let folderResp = await networkRequest(`http://${hostName}:3100/folder/${folderId}`, { method: "GET", credentials: "include" });
    let folderData = await folderResp.data;
    allData = folderData?.data;
  };
  const getBreadcrumbData = async () => {
    let parentResp = await networkRequest(`http://${hostName}:3100/parents/${folderId}`, { method: "GET", credentials: "include" });
    let parentData = await parentResp.data;
    parentFolders = parentData?.data;
  };
  const fetchData = async () => {
    await getFolderData();
    await getBreadcrumbData();
    loading = false;
  };
  const pageChange = async () => {
    folderId = data?.folder.id;
    await fetchData();
  };
  const parentChange = () => {
    set_store_value(slugData, $slugData.current = folderName, $slugData);
    set_store_value(slugData, $slugData.previous = parentFolders ?? [], $slugData);
  };
  onDestroy(() => {
    set_store_value(slugData, $slugData.defaultVal = "", $slugData);
    set_store_value(slugData, $slugData.current = "", $slugData);
    set_store_value(slugData, $slugData.previous = [], $slugData);
  });
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  folderName = allData?.folderName;
  data && pageChange();
  parentFolders && parentChange();
  {
    {
      if ($needToFetchDataInFolder) {
        fetchData();
        set_store_value(needToFetchDataInFolder, $needToFetchDataInFolder = false, $needToFetchDataInFolder);
      }
    }
  }
  $$unsubscribe_slugData();
  $$unsubscribe_needToFetchDataInFolder();
  return `
${$$result.head += `<!-- HEAD_svelte-m0idwz_START -->${$$result.title = `<title>Blazem | ${escape(folderName ?? "")}</title>`, ""}<!-- HEAD_svelte-m0idwz_END -->`, ""}
${validate_component(Loading, "Loading").$$render($$result, { loading }, {}, {
    default: () => {
      return `${validate_component(DataContainer, "DataContainer").$$render($$result, { allData: allData.data }, {}, {})}`;
    }
  })}`;
});
export {
  Page as default
};
