import { c as create_ssr_component, v as validate_component } from "../../../../chunks/index2.js";
import { C as ColourPanel, P as PanelItem } from "../../../../chunks/PanelItem.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${$$result.head += `<!-- HEAD_svelte-18qfo0l_START -->${$$result.title = `<title>Blazem | Stats</title>`, ""}<!-- HEAD_svelte-18qfo0l_END -->`, ""}
<div class="grid grid-cols-4 mt-5 gap-2"><div class="flex flex-col gap-2">${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Number of files backed up</h1>
            <div class="mt-1">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `831`;
        }
      })}</div>`;
    }
  })}

        ${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Number of folders being watched</h1>
            <div class="mt-1">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `412`;
        }
      })}</div>`;
    }
  })}
        ${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Your billing</h1>
            <div class="mt-1">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `$12`;
        }
      })}</div>`;
    }
  })}</div>
    <div class="flex flex-col gap-2">${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Most viewed files</h1>
            <div class="mt-2">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `testing.txt`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `testing.txt`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `testing.txt`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `testing.txt`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `testing.txt`;
        }
      })}</div>`;
    }
  })}</div>
    <div class="flex flex-col gap-2">${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Blazem uptime</h1>
            <div class="mt-1">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `1043 days`;
        }
      })}</div>`;
    }
  })}</div></div>`;
});
export {
  Page as default
};
