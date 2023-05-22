import { goto } from "$app/navigation";
import type { NetworkResponse } from "$lib/types";

export async function networkRequest(url: string, options: object): Promise<NetworkResponse> {
    try {
        let resp = await fetch(url, options).then((res) => {
            if (res.status == 401) {
                goto("/");
            }

            return {
                code: 200,
                msg: "lets go",
                data: res.json(),
            }
        });
        return resp;
    } catch (e) {
        console.log(e);
        return {
            code: 500,
            msg: "OOps",
            data: undefined
        }
    }
}