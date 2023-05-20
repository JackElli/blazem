import { c as create_ssr_component, b as add_attribute } from './index2-e2d3d016.js';

const Panel = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `
<div${add_attribute("class", `rounded-sm p-2 scroll-pb-44 shadow-md shadow-gray-400 ${$$props.class}`, 0)}>${slots.default ? slots.default({}) : ``}</div>`;
});

export { Panel as P };
//# sourceMappingURL=Panel-37be2dea.js.map
