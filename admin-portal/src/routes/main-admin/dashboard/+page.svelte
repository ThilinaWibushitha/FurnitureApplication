<script>
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";
    let salesToday = 0;
    let salesThisMonth = 0;
    let paymentRequests = [];
    let deleteRequests = [];
    let passwordRequests = [];
    onMount(async () => {
        const sales = await api.get("/reports/sales-summary");
        salesToday = sales.today;
        salesThisMonth = sales.this_month;
        paymentRequests = await api.get("/payments/requests");
        deleteRequests = await api.get("/users/delete-requests");
        passwordRequests = await api.get("/users/password-requests");
    });
    async function approveDeleteRequest(id) {
        await api.put(`/users/delete-requests/${id}/approve`);
        deleteRequests = deleteRequests.filter((r) => r.id !== id);
    }
    async function cancelDeleteRequest(id) {
        await api.put(`/users/delete-requests/${id}/cancel`);
        deleteRequests = deleteRequests.filter((r) => r.id !== id);
    }
    async function approvePasswordRequest(id) {
        await api.put(`/users/password-requests/${id}/approve`);
        passwordRequests = passwordRequests.filter((r) => r.id !== id);
    }
</script>

<h1>Main Admin Dashboard</h1>

<h2>Sales</h2>
<p>Today: ${salesToday}</p>
<p>This month: ${salesThisMonth}</p>

<h2>Payment Requests</h2>
<table>
    <thead>
        <tr>
            <th>Order ID</th>
            <th>Amount</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        {#each paymentRequests as request}
        <tr>
            <td>{request.order_id}</td>
            <td>{request.amount}</td>
            <td>
                <button on:click={() => goto(`/orders/${request.order_id}`)}>View Order</button>
            </td>
        </tr>
        {/each}
    </tbody>
</table>

<h2>Client Delete Requests</h2>
<table>
    <thead>
        <tr>
            <th>User ID</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        {#each deleteRequests as request}
        <tr>
            <td>{request.user_id}</td>
            <td>
                <button on:click={() => approveDeleteRequest(request.id)}>Approve</button>
                <button on:click={() => cancelDeleteRequest(request.id)}>Cancel</button>
            </td>
        </tr>
        {/each}
    </tbody>
</table>

<h2>Password Change Requests</h2>
<table>
    <thead>
        <tr>
            <th>User ID</th>
            <th>Action</th>
        </tr>
    </thead>
    <tbody>
        {#each passwordRequests as request}
        <tr>
            <td>{request.user_id}</td>
            <td>
                <button on:click={() => approvePasswordRequest(request.id)}>Approve</button>
            </td>
        </tr>
        {/each}
    </tbody>
</table>

<style>
    h1, h2 {
        color: #333;
    }
    table {
        width: 100%;
        border-collapse: collapse;
        margin-top: 20px;
    }
    th, td {
        border: 1px solid #ddd;
        padding: 8px;
        text-align: left;
    }
    th {
        background-color: #f2f2f2;
    }
</style>