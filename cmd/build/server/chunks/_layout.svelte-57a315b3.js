import { c as create_ssr_component, a as subscribe, v as validate_component, b as add_attribute, e as escape, d as each } from './index2-e2d3d016.js';
import { n as needToFetchDataInFolder, s as slugData } from './stores2-fdb6e3fc.js';
import { s as spinner } from './spinner-75d0e167.js';
import { A as AddObjectModal } from './AddObjectModal-f61a07b7.js';
import './index-fdb9b5cd.js';

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

export { Layout as default };
//# sourceMappingURL=_layout.svelte-57a315b3.js.map
