import { c as create_ssr_component, a as subscribe, b as add_attribute, e as escape, d as each, v as validate_component } from "../../../../chunks/index2.js";
import { s as slugData, n as needToFetchDataInFolder } from "../../../../chunks/stores2.js";
import { s as spinner } from "../../../../chunks/spinner.js";
import { A as AddObjectModal } from "../../../../chunks/AddObjectModal.js";
const Breadcrumb = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let data;
  let $slugData, $$unsubscribe_slugData;
  $$unsubscribe_slugData = subscribe(slugData, (value) => $slugData = value);
  data = $slugData;
  $$unsubscribe_slugData();
  return `${$slugData.current && $slugData.defaultVal && $slugData.current != "" && $slugData.defaultVal != "" ? `<div${add_attribute("class", `flex gap-1 ${$$props.class}`, 0)}>${data.defaultVal != data.current ? `<a class="text-lg underline text-[#3d3d75]" href="/folders">${escape(data.defaultVal)}</a>
            <p class="font-medium text-gray-600 text-lg">/</p>` : ``}

        ${data.previous.length > 0 ? `${each(data.previous, (previousFolder) => {
    return `<a class="text-lg text-[#3d3d75] underline"${add_attribute("href", `/folder/${previousFolder.key}`, 0)}>${escape(previousFolder.folderName)}</a>
                <p class="text-gray-600 text-lg">/</p>`;
  })}` : ``}
        <p class="font-medium text-gray-600 text-lg">${escape(data.current)}</p></div>` : `<img${add_attribute("src", spinner, 0)} alt="spinner" class="animate-spin w-7">`}`;
});
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $$unsubscribe_needToFetchDataInFolder;
  $$unsubscribe_needToFetchDataInFolder = subscribe(needToFetchDataInFolder, (value) => value);
  let addObjectVisible = false;
  let $$settled;
  let $$rendered;
  do {
    $$settled = true;
    $$rendered = `<div class="w-10/12 z-10">${validate_component(Breadcrumb, "Breadcrumb").$$render($$result, {}, {}, {})}
    ${validate_component(AddObjectModal, "AddObjectModal").$$render(
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
    <button class="flex justify-center mt-2 items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"><p class="ml-2 mr-2">Add object</p></button></div>
<div class="mt-4">${slots.default ? slots.default({}) : ``}</div>`;
  } while (!$$settled);
  $$unsubscribe_needToFetchDataInFolder();
  return $$rendered;
});
export {
  Layout as default
};
