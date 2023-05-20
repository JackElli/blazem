import { c as create_ssr_component, b as add_attribute, v as validate_component, d as each, e as escape, h as createEventDispatcher } from './index2-e2d3d016.js';
import { F as Folder } from './Folder-164df24f.js';
import { P as Panel } from './Panel-37be2dea.js';

const css$1 = {
  code: "h1.svelte-ues1oj.svelte-ues1oj{margin:0}.modal_position.svelte-ues1oj.svelte-ues1oj{width:100%;height:100%;background:rgba(0, 0, 0, 0.6);position:fixed;top:0;left:0;z-index:10}.modal_container.svelte-ues1oj.svelte-ues1oj{position:relative;padding:20px}.title.svelte-ues1oj.svelte-ues1oj{font-weight:500}.view_data_modal.svelte-ues1oj.svelte-ues1oj{width:500px;background:white;display:block;margin:0 auto;margin-top:30px;border-radius:3px}.data_main.svelte-ues1oj.svelte-ues1oj{margin-top:20px}.data_attr.svelte-ues1oj.svelte-ues1oj{align-items:center;align-items:stretch;border-bottom:1px solid #bbb;display:flex;flex-grow:1;gap:5px;padding-right:1px}.data_attr.svelte-ues1oj p.svelte-ues1oj{flex:1 1;margin:0;padding:2px}.data_attr.svelte-ues1oj.svelte-ues1oj:last-child{border-bottom:none}.data_attr_key.svelte-ues1oj.svelte-ues1oj{align-items:center;color:#2c2c54;display:flex;font-weight:500;padding-left:3px;padding-right:2px;width:70px}",
  map: null
};
const ViewDataModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let dataAttrs;
  let { visible = false } = $$props;
  let { data } = $$props;
  createEventDispatcher();
  if ($$props.visible === void 0 && $$bindings.visible && visible !== void 0)
    $$bindings.visible(visible);
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  $$result.css.add(css$1);
  dataAttrs = data.data;
  return `
${visible ? `<div class="modal_position svelte-ues1oj"><div class="view_data_modal svelte-ues1oj"><div class="modal_container svelte-ues1oj"><h1 class="title svelte-ues1oj">${escape(data.key)}</h1>
                <div class="data_main svelte-ues1oj">${each(Object.entries(dataAttrs), ([key, dataAttr]) => {
    return `${key != "key" ? `<div class="data_attr w-4/5 svelte-ues1oj"><div class="data_attr_key svelte-ues1oj"><p class="svelte-ues1oj">${escape(key)}</p></div>
                                <p class="svelte-ues1oj">${escape(dataAttr)}</p>
                            </div>` : ``}`;
  })}</div>
                <button class="flex absolute right-5 bottom-1 justify-center items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400"><p class="ml-2 mr-2">Delete</p></button></div></div></div>` : ``}`;
});
const css = {
  code: ".data.svelte-ae7gwi:hover .folder_key.svelte-ae7gwi{text-decoration:underline}",
  map: null
};
const Data = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let dataAttrs;
  let { data } = $$props;
  let viewDataVisible = false;
  const clamp = (value) => {
    if (value.length > 10) {
      return value.substring(0, 10) + "...";
    }
    return value;
  };
  const valueClamp = (value) => {
    let newValue = value;
    if (newValue.length > 60) {
      return newValue.substring(0, 60) + "...";
    }
    return newValue;
  };
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  $$result.css.add(css);
  let $$settled;
  let $$rendered;
  do {
    $$settled = true;
    dataAttrs = data.data;
    $$rendered = `${validate_component(ViewDataModal, "ViewDataModal").$$render(
      $$result,
      { data, visible: viewDataVisible },
      {
        visible: ($$value) => {
          viewDataVisible = $$value;
          $$settled = false;
        }
      },
      {}
    )}

<div class="data w-full pb-1 border-l-4 border-l-blue-500 cursor-pointer svelte-ae7gwi"><div class="flex pl-3 relative h-full"><div class="min-w-[80px]"><h1 class="font-medium">Type</h1>
            <h1 class="font-normal text-lg">Data</h1></div>

        <div class="min-w-[120px]"><h1 class="font-medium">Key</h1>
            <h1 class="folder_key font-normal text-lg svelte-ae7gwi">${escape(clamp(data.key))}</h1></div>

        <div class="flex gap-20">${dataAttrs.type == "text" ? `${each(Object.entries(dataAttrs), ([key, dataAttr]) => {
      return `${!["key", "folder", "type"].includes(key) ? `<div><h1 class="font-medium">${escape(key)}</h1>
                            <h1 class="font-normal text-lg">${escape(valueClamp(dataAttr))}</h1>
                        </div>` : ``}`;
    })}` : `
                <img${add_attribute("src", dataAttrs.value, 0)}>`}</div></div>
</div>`;
  } while (!$$settled);
  return $$rendered;
});
const DataContainer = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { allData } = $$props;
  let folders;
  let datas;
  const load = () => {
    folders = allData.filter(function(item) {
      return item.data.type == "folder";
    });
    datas = allData.filter(function(item) {
      return item.data.type != "folder";
    });
  };
  if ($$props.allData === void 0 && $$bindings.allData && allData !== void 0)
    $$bindings.allData(allData);
  allData && load();
  return `${allData && allData.length != 0 ? `<div${add_attribute("class", `flex flex-col gap-4 w-full mb-3 ${$$props.class}`, 0)}>${folders.length > 0 ? `${validate_component(Panel, "Panel").$$render($$result, { class: "bg-white" }, {}, {
    default: () => {
      return `${each(folders, (folder) => {
        return `${validate_component(Folder, "Folder").$$render($$result, { data: folder.data }, {}, {})}`;
      })}`;
    }
  })}` : ``}
        ${datas.length > 0 ? `${validate_component(Panel, "Panel").$$render($$result, { class: "bg-white" }, {}, {
    default: () => {
      return `${each(datas, (data) => {
        return `${validate_component(Data, "Data").$$render($$result, { data }, {}, {})}`;
      })}`;
    }
  })}` : ``}</div>` : `<h1${add_attribute("class", `${$$props.class}`, 0)}>No data found.</h1>`}`;
});

export { DataContainer as D };
//# sourceMappingURL=DataContainer-102101e7.js.map
