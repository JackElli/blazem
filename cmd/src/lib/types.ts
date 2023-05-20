export type SidebarItem = {
    value: string;
    href: string;
    active: boolean;
}

export interface NetworkResponse {
    code: number;
    msg: string;
    data?: object;
}