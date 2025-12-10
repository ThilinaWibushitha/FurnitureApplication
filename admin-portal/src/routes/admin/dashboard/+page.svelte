<script>
    import { onMount } from "svelte";
    import { api } from "$lib/api";

    let deleteRequests = [];

    onMount(async () => {
        deleteRequests = await api.get("/users/delete-requests");
    });

    async function approveDeleteRequest(id) {
        await api.put(`/users/delete-requests/${id}/approve`);
        deleteRequests = deleteRequests.filter((r) => r.id !== id);
    }

    async function cancelDeleteRequest(id) {
        await api.put(`/users/delete-requests/${id}/cancel`);
        deleteRequests = deleteRequests.filter((r) => r.id !== id);
    }
</script>

<h1>Admin Dashboard</h1>

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