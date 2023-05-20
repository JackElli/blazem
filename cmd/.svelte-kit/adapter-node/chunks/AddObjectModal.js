import { c as create_ssr_component, h as createEventDispatcher, b as add_attribute, t as tick, v as validate_component } from "./index2.js";
const randomKeyVals = "0123456789abcdef";
const keyLength = 16;
const generateKey = () => {
  let randomKey = "";
  for (let i = 0; i < keyLength; i++) {
    let randIndex = Math.floor(Math.random() * randomKeyVals.length);
    randomKey += randomKeyVals[randIndex];
  }
  return randomKey;
};
const AddData = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { active } = $$props;
  let dataKey;
  let dataValue;
  createEventDispatcher();
  const activeChanged = async () => {
    await tick();
    if (active == "data") {
      dataKey.placeholder = generateKey();
      dataValue.focus();
    }
  };
  if ($$props.active === void 0 && $$bindings.active && active !== void 0)
    $$bindings.active(active);
  active && activeChanged();
  return `<p class="text-gray-300 text-xs mt-2">Document Key</p>
<input class="border border-gray-300 rounded-sm w-80 h-7 pl-2" type="text" placeholder="testkey"${add_attribute("this", dataKey, 0)}>
<br>
<p class="text-gray-300 text-xs mt-2">Document value</p>
<textarea class="border border-gray-300 h-24 pl-2 pt-2 resize-none w-80"${add_attribute("this", dataValue, 0)}></textarea>
<p class="text-[#3d3d75] hover:underline cursor-pointer">Upload a file?</p>
<br>
<button class="flex justify-center items-center bg-white border-l-4 border-l-[#3d3d75] h-8 border-gray-300 border hover:border-gray-400 relative"><p class="ml-2 mr-2">Add</p></button>`;
});
const AddObjectModal = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { visible = false } = $$props;
  let active = "data";
  createEventDispatcher();
  if ($$props.visible === void 0 && $$bindings.visible && visible !== void 0)
    $$bindings.visible(visible);
  return `
${visible ? `<div class="modal-position w-full h-full fixed top-0 left-0 z-20 bg-black bg-opacity-40"><div class="w-[400px] bg-white block mx-auto mt-8 rounded-sm"><div class="p-5"><h1 class="font-medium">Add object</h1>
                <div class="flex gap-1 text-sm mt-2"><h1${add_attribute("class", `hover:underline cursor-pointer ${"underline"}`, 0)}>Data
                    </h1>
                    <h1${add_attribute("class", `hover:underline cursor-pointer ${""}`, 0)}>Folder
                    </h1></div>
                ${`${validate_component(AddData, "AddData").$$render($$result, { active }, {}, {})}`}</div></div></div>` : ``}`;
});
export {
  AddObjectModal as A
};
