<script>
    import { goto } from "$app/navigation";
</script>

<h1>Welcome to the Furniture Admin Portal</h1>

<button on:click={() => goto("/main-admin/login")}>Main Admin Login</button>
<button on:click={() => goto("/admin/login")}>Admin Login</button>

<style>
    h1 {
        color: #333;
        text-align: center;
    }
    button {
        display: block;
        width: 200px;
        padding: 10px;
        margin: 20px auto;
        background-color: #333;
        color: #fff;
        border: none;
        cursor: pointer;
    }
</style>