<script>
    import { api } from "$lib/api";
    let email = "";
    let password = "";
    let error = "";
    let success = "";

    async function createAdmin() {
        try {
            await api.post("/auth/admin/create", { email, password });
            success = "Admin account created successfully";
            email = "";
            password = "";
        } catch (e) {
            error = e.message;
        }
    }
</script>

<h1>Create Admin Account</h1>

<form on:submit|preventDefault={createAdmin}>
    <label for="email">Email</label>
    <input type="email" id="email" bind:value={email} />
    <label for="password">Password</label>
    <input type="password" id="password" bind:value={password} />
    <button type="submit">Create Admin</button>
</form>

{#if error}
<p>{error}</p>
{/if}

{#if success}
<p>{success}</p>
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