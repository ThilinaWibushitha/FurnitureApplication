<script>
    import { api } from "$lib/api";
    import { page } from "$app/stores";
    let password = "";
    let confirmPassword = "";
    let error = "";
    let success = "";
    async function resetPassword() {
        if (password !== confirmPassword) {
            error = "Passwords do not match";
            return;
        }
        try {
            const token = $page.url.searchParams.get("token");
            await api.post("/auth/main-admin/reset-password", { token, password });
            success = "Password reset successfully";
        } catch (e) {
            error = e.message;
        }
    }
</script>

<h1>Reset Password</h1>

<form on:submit|preventDefault={resetPassword}>
    <label for="password">Password</label>
    <input type="password" id="password" bind:value={password} />
    <label for="confirm-password">Confirm Password</label>
    <input type="password" id="confirm-password" bind:value={confirmPassword} />
    <button type="submit">Reset Password</button>
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