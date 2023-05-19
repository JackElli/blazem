import { c as create_ssr_component, b as add_attribute } from './index2-e2d3d016.js';

const ColourPanel = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div${add_attribute("class", `bg-white rounded-sm p-2 scroll-pb-44 shadow-md shadow-gray-400 ${$$props.class}`, 0)}><div${add_attribute("class", `border-l-4 border-l-blue-400 pl-2`, 0)}>${slots.default ? slots.default({}) : ``}</div></div>`;
});
const PanelItem = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<p class="py-2">${slots.default ? slots.default({}) : ``}</p>`;
});

export { ColourPanel as C, PanelItem as P };
//# sourceMappingURL=PanelItem-27b833ad.js.map
