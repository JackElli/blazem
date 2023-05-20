import { writable } from 'svelte/store';

export const headerMessage = writable('');
export const previousPages = writable({})
export const needToFetchDataInFolder = writable(false)
export const sidebarActive = writable("")

type BreadcrumbData = {
    defaultVal: string;
    current: string;
    previous: object[];
}

export const slugData = writable<BreadcrumbData>({
    defaultVal: "",
    current: "",
    previous: []
})


