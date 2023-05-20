import { c as create_ssr_component, v as validate_component, e as escape, a as subscribe, b as add_attribute } from './index2-e2d3d016.js';
import { v as versionNum } from './global-6b632810.js';
import { p as page } from './stores-4309ae18.js';

const Tab = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $page, $$unsubscribe_page;
  $$unsubscribe_page = subscribe(page, (value) => $page = value);
  let { href } = $$props;
  if ($$props.href === void 0 && $$bindings.href && href !== void 0)
    $$bindings.href(href);
  $$unsubscribe_page();
  return `<a${add_attribute("href", href, 0)}${add_attribute("class", `font-medium px-2 py-1 rounded-sm text-md cursor-pointer hover:bg-gray-100 ${href == $page.url.pathname ? "bg-gray-100" : ""}`, 0)}>${slots.default ? slots.default({}) : ``}</a>`;
});
const Header = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="w-full h-12"><div class="flex justify-between items-center h-full"><div class="flex items-center"><a href="/folders" class="font-medium text-3xl text-[#3b82f6] cursor-pointer">Blazem
            </a>

            <input class="ml-4 rounded-sm w-52 text-sm px-1 py-1 focus:py-2 shadow-md bg-gray-100 focus:bg-white outline-none focus:w-64" placeholder="Search..."></div>

        ${validate_component(Tab, "Tab").$$render($$result, { href: "/help" }, {}, {
    default: () => {
      return `Need help?`;
    }
  })}</div></div>`;
});
const SidebarDivider = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<h1 class="text-sm font-medium text-gray-400 pl-2 py-2 mt-2">${slots.default ? slots.default({}) : ``}</h1>`;
});
const SidebarTab = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let $page, $$unsubscribe_page;
  $$unsubscribe_page = subscribe(page, (value) => $page = value);
  let { href } = $$props;
  if ($$props.href === void 0 && $$bindings.href && href !== void 0)
    $$bindings.href(href);
  $$unsubscribe_page();
  return `<button class="${"block w-full text-left rounded-md py-1 pl-2 mt-1 hover:bg-gray-300 cursor-pointer " + escape(
    $page.url.pathname == href ? "bg-gray-300 font-medium" : "",
    true
  )}">${slots.default ? slots.default({}) : ``}</button>`;
});
const SidebarMain = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="w-11/12 mx-auto">${validate_component(SidebarDivider, "SidebarDivider").$$render($$result, {}, {}, {
    default: () => {
      return `Views`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/folders" }, {}, {
    default: () => {
      return `Folders`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/stats" }, {}, {
    default: () => {
      return `Stats`;
    }
  })}
    ${validate_component(SidebarDivider, "SidebarDivider").$$render($$result, {}, {}, {
    default: () => {
      return `Files`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/export" }, {}, {
    default: () => {
      return `Export`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/import" }, {}, {
    default: () => {
      return `Import`;
    }
  })}
    ${validate_component(SidebarDivider, "SidebarDivider").$$render($$result, {}, {}, {
    default: () => {
      return `Nodes`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/nodes" }, {}, {
    default: () => {
      return `Nodes`;
    }
  })}
    ${validate_component(SidebarDivider, "SidebarDivider").$$render($$result, {}, {}, {
    default: () => {
      return `Search`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/search" }, {}, {
    default: () => {
      return `Advanced Search`;
    }
  })}</div>`;
});
const SidebarNoService = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  return `<div class="w-11/12 mx-auto">${validate_component(SidebarDivider, "SidebarDivider").$$render($$result, {}, {}, {
    default: () => {
      return `Add Service`;
    }
  })}
    ${validate_component(SidebarTab, "SidebarTab").$$render($$result, { href: "/backup" }, {}, {
    default: () => {
      return `Add Backup`;
    }
  })}</div>`;
});
const Sidebar = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { service = false } = $$props;
  if ($$props.service === void 0 && $$bindings.service && service !== void 0)
    $$bindings.service(service);
  return `<div class="bg-gray-200 w-40 h-screen border-r border-r-gray-300 sticky top-0">${service ? `${validate_component(SidebarMain, "SidebarMain").$$render($$result, {}, {}, {})}` : `${validate_component(SidebarNoService, "SidebarNoService").$$render($$result, {}, {}, {})}`}
    <p class="absolute bottom-0 left-0 right-0 text-center text-xs text-gray-400">Jack Ellis 2023
    </p></div>`;
});
const Layout = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let { data } = $$props;
  if ($$props.data === void 0 && $$bindings.data && data !== void 0)
    $$bindings.data(data);
  return `<div class="flex gap-6 w-full">${validate_component(Sidebar, "Sidebar").$$render($$result, { service: data.service }, {}, {})}
    <div class="w-10/12">${validate_component(Header, "Header").$$render($$result, {}, {}, {})}
        <div class="mt-2">${slots.default ? slots.default({}) : ``}</div></div>
    <p class="fixed bottom-2 right-2">v${escape(versionNum)}</p></div>`;
});

export { Layout as default };
//# sourceMappingURL=_layout.svelte-91628c5d.js.map
