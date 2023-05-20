import { w as writable } from "./index.js";
const needToFetchDataInFolder = writable(false);
const slugData = writable({
  defaultVal: "",
  current: "",
  previous: []
});
export {
  needToFetchDataInFolder as n,
  slugData as s
};
