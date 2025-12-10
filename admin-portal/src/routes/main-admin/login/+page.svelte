<script>
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    let email = "";
    let password = "";
    let error = "";
    async function login() {
        try {
            await api.post("/auth/main-admin/login", { email, password });
            goto("/main-admin/dashboard");
        } catch (e) {
            error = e.message;
        }
    }
</script>

<h1>Main Admin Login</h1>

<form on:submit|preventDefault={login}>
    <label for="email">Email</label>
    <input type="email" id="email" bind:value={email} />
    <label for="password">Password</label>
    <input type="password" id="password" bind:value={password} />
    <button type="submit">Login</button>
</form>

{#if error}
<p>{error}</p>
{/if}

<style>
    h1 {
        color: #333;
    }
    form {
        display: flex;
        flex-direction: column;
        width: 300px;
        margin-top: 20px;
    }
    label {
        margin-top: 10px;
    }
    input {
        padding: 5px;
        margin-top: 5px;
    }
    button {
        margin-top: 20px;
        padding: 10px;
        background-color: #333;
        color: #fff;
        border: none;
        cursor: pointer;
    }
</style>