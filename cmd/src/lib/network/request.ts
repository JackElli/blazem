import { goto } from "$app/navigation";
import type { NetworkResponse } from "$lib/types";

export async function networkRequest(url: string, options: object): Promise<NetworkResponse> {

    try {
        const response = await fetch(url, options);
        if (response.status == 403) {
            return {
                code: 403,
                msg: "Not good"
            }
        }
        if (response.status == 401) {
            goto("/login")
            return {
                code: 401,
                msg: "unauthorised",
                data: {}
            }
        }
        return await response.json();
    }
    catch (error) {
        return {
            code: 500,
            msg: "Error occurred",
        }
    }
}





