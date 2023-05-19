import { c as create_ssr_component, b as add_attribute, e as escape } from "../../../chunks/index2.js";
import { v as versionNum } from "../../../chunks/global.js";
const Page = create_ssr_component(($$result, $$props, $$bindings, slots) => {
  let username;
  let password;
  return `<div class="w-96 h-96 px-10 pt-4 shadow-md mx-auto mt-20 bg-gray-200 rounded-md"><div class="w-72 h-96"><h1 class="font-medium text-[70px] text-[#3b82f6] text-center">Blazem
        </h1>
        <p class="font-medium">Log in</p>
        <form><input placeholder="Username" class="py-2 w-full pl-4 rounded-md mt-3"${add_attribute("value", username, 0)}>
            <input placeholder="Password" class="py-2 w-full pl-4 rounded-md mt-3" type="password"${add_attribute("value", password, 0)}>
            <br>
            <button type="submit" class="bg-gray-300 py-1 px-3 mt-3 hover:bg-gray-100 rounded-sm">Log in</button></form>
        <p class="mt-6 text-gray-500">Pre Alpha</p></div>
    <p class="fixed right-2 bottom-2">v${escape(versionNum)}</p></div>`;
});
export {
  Page as default
};
