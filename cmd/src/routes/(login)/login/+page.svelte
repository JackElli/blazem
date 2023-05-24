<script lang="ts">
    import { goto } from "$app/navigation";
    import { hostName, versionNum } from "../../../global";

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

<div class="flex w-1/2 min-w-[800px] gap-10 mx-auto mt-52">
    <div class="px-10 pt-4 pb-10 shadow-md bg-white rounded-sm">
        <div class="w-72">
            <h1 class="font-medium text-6xl text-[#3b82f6] text-center">
                Blazem
            </h1>
            <p class="font-medium mt-5">Log in</p>
            <form on:submit|preventDefault={login}>
                <input
                    bind:value={username}
                    placeholder="Username"
                    class="py-2 w-full pl-4 rounded-md mt-3 border border-gray-300 focus:bg-zinc-50 outline-none"
                />
                <input
                    bind:value={password}
                    placeholder="Password"
                    class="py-2 w-full pl-4 rounded-md mt-3 border border-gray-300 focus:bg-zinc-50 outline-none"
                    type="password"
                />
                <br />
                <button
                    type="submit"
                    class="bg-[#3b82f6] text-white py-1 px-3 mt-3 hover:bg-[#468afa] rounded-sm"
                    on:click={login}>Log in</button
                >
            </form>
        </div>
    </div>

    <div class="w-96 flex flex-col justify-center">
        <h1 class="text-4xl">You've been signed out.</h1>
        <p class="text-xl text-gray-500">Sign back in to view your data.</p>
    </div>
</div>

<p class="fixed right-2 bottom-2">v{versionNum}</p>
