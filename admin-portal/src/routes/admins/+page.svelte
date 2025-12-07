<script>
    import { onMount } from "svelte";
    import { getAdmins, createAdmin } from "$lib/api";

    let admins = [];
    let isLoading = true;
    let showAddModal = false;
    let errorMessage = "";
    let successMessage = "";

    let newAdmin = {
        username: "",
        email: "",
        password: "",
        full_name: "",
        phone: "",
    };

    onMount(async () => {
        await loadAdmins();
    });

    async function loadAdmins() {
        isLoading = true;
        try {
            admins = await getAdmins();
        } catch (error) {
            console.error(error);
            errorMessage = "Failed to load admins.";
        } finally {
            isLoading = false;
        }
    }

    async function handleCreateAdmin() {
        errorMessage = "";
        successMessage = "";

        try {
            await createAdmin(newAdmin);
            successMessage = "Admin created successfully!";
            showAddModal = false;
            newAdmin = {
                username: "",
                email: "",
                password: "",
                full_name: "",
                phone: "",
            };
            await loadAdmins();
        } catch (error) {
            console.error(error);
            errorMessage = "Failed to create admin: " + error.message;
        }
    }
</script>

<div class="container">
    <div class="header-actions">
        <h1>Admin Management</h1>
        <button class="btn-primary" on:click={() => (showAddModal = true)}>
            + Add New Admin
        </button>
    </div>

    {#if errorMessage}
        <div class="alert error">{errorMessage}</div>
    {/if}
    {#if successMessage}
        <div class="alert success">{successMessage}</div>
    {/if}

    <div class="card-grid">
        {#if isLoading}
            <p>Loading...</p>
        {:else if admins.length === 0}
            <p>No admins found.</p>
        {:else}
            {#each admins as admin}
                <div class="card admin-card">
                    <div class="card-header">
                        <h3>{admin.full_name}</h3>
                        <span class="badge {admin.account_type}"
                            >{admin.account_type === "main_admin"
                                ? "Main Admin"
                                : "Admin"}</span
                        >
                    </div>
                    <div class="card-body">
                        <p><strong>Username:</strong> {admin.username}</p>
                        <p><strong>Email:</strong> {admin.email}</p>
                        <p>
                            <strong>Status:</strong>
                            {admin.is_active ? "Active" : "Inactive"}
                        </p>
                    </div>
                </div>
            {/each}
        {/if}
    </div>
</div>

{#if showAddModal}
    <div class="modal-backdrop">
        <div class="modal">
            <h2>Add New Admin</h2>
            <form on:submit|preventDefault={handleCreateAdmin}>
                <div class="form-group">
                    <label for="username">Username</label>
                    <input
                        type="text"
                        id="username"
                        bind:value={newAdmin.username}
                        required
                    />
                </div>
                <div class="form-group">
                    <label for="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        bind:value={newAdmin.email}
                        required
                    />
                </div>
                <div class="form-group">
                    <label for="full_name">Full Name</label>
                    <input
                        type="text"
                        id="full_name"
                        bind:value={newAdmin.full_name}
                        required
                    />
                </div>
                <div class="form-group">
                    <label for="password">Password</label>
                    <input
                        type="password"
                        id="password"
                        bind:value={newAdmin.password}
                        required
                    />
                </div>

                <div class="modal-actions">
                    <button
                        type="button"
                        class="btn-secondary"
                        on:click={() => (showAddModal = false)}>Cancel</button
                    >
                    <button type="submit" class="btn-primary"
                        >Create Admin</button
                    >
                </div>
            </form>
        </div>
    </div>
{/if}

<style>
    .container {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .header-actions {
        display: flex;
        justify-content: space-between;
        align-items: center;
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

    .badge {
        padding: 0.25rem 0.75rem;
        border-radius: 999px;
        font-size: 0.875rem;
        font-weight: 500;
        background: #e2e8f0;
        color: #475569;
    }

    .badge.main_admin {
        background: #dbeafe;
        color: #1e40af;
    }

    .btn-primary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-weight: 600;
        cursor: pointer;
        transition: opacity 0.2s;
    }

    .btn-secondary {
        background: #e2e8f0;
        color: #475569;
        border: none;
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-weight: 600;
        cursor: pointer;
    }

    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 100;
    }

    .modal {
        background: white;
        padding: 2rem;
        border-radius: 16px;
        width: 100%;
        max-width: 500px;
    }

    .form-group {
        margin-bottom: 1rem;
    }

    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }

    .form-group input {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #e2e8f0;
        border-radius: 8px;
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 1rem;
        margin-top: 2rem;
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
</style>
