import { c as create_ssr_component, b as add_attribute } from './index2-e2d3d016.js';
import { s as spinner } from './spinner-75d0e167.js';

const Loading = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { loading = false } = $$props;
  if ($$props.loading === void 0 && $$bindings.loading && loading !== void 0)
    $$bindings.loading(loading);
  return `${loading ? `<img${add_attribute("src", spinner, 0)} alt="spinner"${add_attribute("class", `animate-spin w-10 ${$$props.class} fixed top-0 bottom-0 left-0 right-0 my-auto mx-auto`, 0)}>` : `${slots.default ? slots.default({}) : ``}`}`;
});

export { Loading as L };
//# sourceMappingURL=Loading-0ac87029.js.map
