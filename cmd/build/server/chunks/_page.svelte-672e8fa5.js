import { c as create_ssr_component, v as validate_component } from './index2-e2d3d016.js';
import { C as ColourPanel, P as PanelItem } from './PanelItem-27b833ad.js';

const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `${$$result.head += `<!-- HEAD_svelte-127koss_START -->${$$result.title = `<title>Blazem | Recents</title>`, ""}<!-- HEAD_svelte-127koss_END -->`, ""}
<div class="grid grid-cols-3 mt-5 gap-2">${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Recently backed up files</h1>
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
  })}
    <div class="flex flex-col gap-2">${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Recently viewed files</h1>
            <div class="mt-2">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `darthdad.mp4`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `hello.jpg`;
        }
      })}</div>`;
    }
  })}
        ${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Recently deleted files</h1>
            <div class="mt-2">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `passwds.txt`;
        }
      })}
                ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `amazing.txt`;
        }
      })}</div>`;
    }
  })}</div>
    ${validate_component(ColourPanel, "ColourPanel").$$render($$result, {}, {}, {
    default: () => {
      return `<h1 class="font-medium">Recent server logs</h1>
        <div class="mt-2">${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `01/05/23 Good!`;
        }
      })}
            ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `01/05/23 Good!`;
        }
      })}
            ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `01/05/23 Bad!`;
        }
      })}
            ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `01/05/23 Bad!`;
        }
      })}
            ${validate_component(PanelItem, "PanelItem").$$render($$result, {}, {}, {
        default: () => {
          return `01/05/23 Bad!`;
        }
      })}</div>`;
    }
  })}</div>`;
});

export { Page as default };
//# sourceMappingURL=_page.svelte-672e8fa5.js.map
