<script>
    import { onMount } from "svelte";
    import {
        getPasswordChangeRequests,
        resolvePasswordChangeRequest,
    } from "$lib/api";

    let requests = [];
    let isLoading = true;
    let errorMessage = "";
    let successMessage = "";

    onMount(async () => {
        await loadRequests();
    });

    async function loadRequests() {
        isLoading = true;
        try {
            requests = await getPasswordChangeRequests();
        } catch (error) {
            console.error(error);
            errorMessage = "Failed to load requests.";
        } finally {
            isLoading = false;
        }
    }

    async function handleResolve(id, approved) {
        errorMessage = "";
        successMessage = "";

        try {
            await resolvePasswordChangeRequest(id, approved);
            successMessage = `Request ${approved ? "approved" : "rejected"} successfully!`;
            await loadRequests();
        } catch (error) {
            console.error(error);
            errorMessage = "Failed to resolve request: " + error.message;
        }
    }

    function formatDate(dateStr) {
        return new Date(dateStr).toLocaleString();
    }
</script>

<div class="container">
    <div class="header-actions">
        <h1>Password Change Requests</h1>
    </div>

    {#if errorMessage}
        <div class="alert error">{errorMessage}</div>
    {/if}
    {#if successMessage}
        <div class="alert success">{successMessage}</div>
    {/if}

    <div class="card-grid">
        {#if isLoading}
            <p class="loading">Loading...</p>
        {:else if requests.length === 0}
            <p class="empty">No pending requests.</p>
        {:else}
            {#each requests as req}
                <div class="card request-card">
                    <div class="card-header">
                        <h3>{req.username}</h3>
                        <span class="badge pending">Pending</span>
                    </div>
                    <div class="card-body">
                        <p><strong>Email:</strong> {req.email}</p>
                        <p>
                            <strong>Requested:</strong>
                            {formatDate(req.created_at)}
                        </p>
                    </div>
                    <div class="card-actions">
                        <button
                            class="btn-reject"
                            on:click={() => handleResolve(req.id, false)}
                            >Reject</button
                        >
                        <button
                            class="btn-approve"
                            on:click={() => handleResolve(req.id, true)}
                            >Approve</button
                        >
                    </div>
                </div>
            {/each}
        {/if}
    </div>
</div>

<style>
    .container {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .header-actions {
        margin-bottom: 2rem;
    }

    h1 {
        color: #fff;
        margin: 0;
    }

    .card-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
    }

    .card {
        background: white;
        border-radius: 12px;
        padding: 1.5rem;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid #eee;
    }

    .card-header h3 {
        margin: 0;
    }

    .badge {
        padding: 0.25rem 0.75rem;
        border-radius: 999px;
        font-size: 0.875rem;
        font-weight: 500;
        background: #fef3c7;
        color: #d97706;
    }

    .card-actions {
        margin-top: 1.5rem;
        display: flex;
        gap: 1rem;
    }

    .btn-approve {
        flex: 1;
        background: #10b981;
        color: white;
        border: none;
        padding: 0.5rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: 600;
    }

    .btn-reject {
        flex: 1;
        background: #fee2e2;
        color: #b91c1c;
        border: none;
        padding: 0.5rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: 600;
    }

    .alert {
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 1.5rem;
    }

    .alert.error {
        background: #fee2e2;
        color: #991b1b;
    }

    .alert.success {
        background: #dcfce7;
        color: #166534;
    }

    .loading,
    .empty {
        color: white;
        text-align: center;
        width: 100%;
        grid-column: 1 / -1;
    }
</style>
