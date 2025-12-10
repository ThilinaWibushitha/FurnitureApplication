<script>
    import { api } from "$lib/api";
    let email = "";
    let error = "";
    let success = "";
    async function requestPasswordReset() {
        try {
            await api.post("/auth/main-admin/request-password-reset", { email });
            success = "Password reset email sent";
        } catch (e) {
            error = e.message;
        }
    }
</script>

<h1>Forgot Password</h1>

<form on:submit|preventDefault={requestPasswordReset}>
    <label for="email">Email</label>
    <input type="email" id="email" bind:value={email} />
    <button type="submit">Request Password Reset</button>
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