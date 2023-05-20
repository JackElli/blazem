function guard(name) {
  return () => {
    throw new Error(`Cannot call ${name}(...) on the server`);
  };
}
const goto = guard("goto");
async function networkRequest(url, options) {
  try {
    let resp = await fetch(url, options).then((res) => {
      if (res.status == 401) {
        goto("/");
      }
      return {
        code: 200,
        msg: "lets go",
        data: res.json()
      };
    });
    return resp;
  } catch (e) {
    console.log(e);
    return {
      code: 500,
      msg: "OOps",
      data: void 0
    };
  }
}
export {
  networkRequest as n
};
