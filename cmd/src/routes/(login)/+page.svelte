<script lang="ts">
    import { goto } from "$app/navigation";
    import { prevent_default } from "svelte/internal";
    import { hostName, versionNum } from "../../global";

    let username: string;
    let password: string;

    const login = async () => {
        let res = await fetch(`http://${hostName}:3100/auth`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({
                username: username,
                password: password,
            }),
        });

        if (res.status != 200) {
            return;
        }

        goto("/folders");
    };
</script>

<div
    class="w-96 h-96 px-10 pt-4 shadow-md mx-auto mt-20 bg-gray-200 rounded-md"
>
    <div class="w-72 h-96">
        <h1 class="font-medium text-[70px] text-[#3b82f6] text-center">
            Blazem
        </h1>
        <p class="font-medium">Log in</p>
        <form on:submit|preventDefault={login}>
            <input
                bind:value={username}
                placeholder="Username"
                class="py-2 w-full pl-4 rounded-md mt-3"
            />
            <input
                bind:value={password}
                placeholder="Password"
                class="py-2 w-full pl-4 rounded-md mt-3"
                type="password"
            />
            <br />
            <button
                type="submit"
                class="bg-gray-300 py-1 px-3 mt-3 hover:bg-gray-100 rounded-sm"
                on:click={login}>Log in</button
            >
        </form>
        <p class="mt-6 text-gray-500">Pre Alpha</p>
    </div>
    <p class="fixed right-2 bottom-2">v{versionNum}</p>
</div>
