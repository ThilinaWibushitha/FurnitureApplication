<script>
    import { onMount } from 'svelte';
    import { getClients, getPendingDeletions, cancelDeletion, confirmDeletion, updateClient } from '$lib/api.js';
    
    let clients = [];
    let pendingDeletions = [];
    let loading = true;
    let searchQuery = '';
    let statusFilter = '';
    let activeTab = 'all';
    
    let showModal = false;
    let modalType = '';
    let selectedClient = null;
    let editData = {};
    let processing = false;
    
    onMount(async () => {
        await loadData();
    });
    
    async function loadData() {
        loading = true;
        try {
            const filter = {};
            if (searchQuery) filter.search = searchQuery;
            if (statusFilter) filter.status = statusFilter;
            
            [clients, pendingDeletions] = await Promise.all([
                getClients(filter),
                getPendingDeletions()
            ]);
        } catch (e) {
            console.error('Failed to load data:', e);
        }
        loading = false;
    }
    
    function handleSearch() {
        loadData();
    }
    
    function openEditModal(client) {
        selectedClient = client;
        editData = {
            phone: client.phone || '',
            address: client.address || '',
            city: client.city || '',
            gender: client.gender || '',
            newsletter_opt_in: client.newsletter_opt_in || false
        };
        modalType = 'edit';
        showModal = true;
    }
    
    function openDeleteModal(client, action) {
        selectedClient = client;
        modalType = action; // 'cancel' or 'confirm'
        showModal = true;
    }
    
    function closeModal() {
        showModal = false;
        selectedClient = null;
        editData = {};
    }
    
    async function handleSave() {
        processing = true;
        try {
            await updateClient(selectedClient.customer_id, editData);
            closeModal();
            await loadData();
        } catch (e) {
            console.error('Failed to update:', e);
        }
        processing = false;
    }
    
    async function handleDeletionAction() {
        processing = true;
        try {
            if (modalType === 'cancel') {
                await cancelDeletion(selectedClient.customer_id);
            } else {
                await confirmDeletion(selectedClient.customer_id);
            }
            closeModal();
            await loadData();
        } catch (e) {
            console.error('Failed to process deletion:', e);
        }
        processing = false;
    }
    
    $: filteredClients = activeTab === 'deletions' ? pendingDeletions : clients;
</script>

<div class="container">
    <header>
        <h1>👥 Client Management</h1>
        <p class="subtitle">Manage client profiles and account requests</p>
    </header>
    
    <!-- Stats -->
    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-icon clients">👤</div>
            <div class="stat-info">
                <span class="stat-value">{clients.length}</span>
                <span class="stat-label">Total Clients</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon active">✓</div>
            <div class="stat-info">
                <span class="stat-value">{clients.filter(c => c.account_status === 'Active').length}</span>
                <span class="stat-label">Active</span>
            </div>
        </div>
        <div class="stat-card warning">
            <div class="stat-icon deletion">⚠</div>
            <div class="stat-info">
                <span class="stat-value">{pendingDeletions.length}</span>
                <span class="stat-label">Pending Deletions</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon points">★</div>
            <div class="stat-info">
                <span class="stat-value">{clients.reduce((sum, c) => sum + (c.loyalty_points || 0), 0)}</span>
                <span class="stat-label">Total Points</span>
            </div>
        </div>
    </div>
    
    <!-- Tabs -->
    <div class="tabs">
        <button class="tab" class:active={activeTab === 'all'} on:click={() => activeTab = 'all'}>
            All Clients
        </button>
        <button class="tab" class:active={activeTab === 'deletions'} on:click={() => activeTab = 'deletions'}>
            Deletion Requests
            {#if pendingDeletions.length > 0}
                <span class="badge">{pendingDeletions.length}</span>
            {/if}
        </button>
    </div>
    
    {#if activeTab === 'all'}
        <!-- Search -->
        <div class="search-bar">
            <input 
                type="text" 
                bind:value={searchQuery} 
                placeholder="Search clients by name or email..."
                on:keyup={(e) => e.key === 'Enter' && handleSearch()}
            />
            <button on:click={handleSearch}>Search</button>
        </div>
    {/if}
    
    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading clients...</p>
        </div>
    {:else}
        <div class="clients-grid">
            {#each filteredClients as client}
                <div class="client-card" class:pending-deletion={client.account_status === 'PendingDeletion'}>
                    <div class="client-header">
                        <div class="avatar">
                            {client.first_name?.[0]}{client.last_name?.[0] || ''}
                        </div>
                        <div class="client-info">
                            <h3>{client.first_name} {client.last_name || ''}</h3>
                            <p>{client.email || 'No email'}</p>
                        </div>
                        <span class="status-badge {client.account_status.toLowerCase().replace('pendingdeletion', 'pending')}">
                            {client.account_status}
                        </span>
                    </div>
                    
                    <div class="client-details">
                        <div class="detail">
                            <span class="label">📱 Phone</span>
                            <span class="value">{client.phone || 'N/A'}</span>
                        </div>
                        <div class="detail">
                            <span class="label">📍 City</span>
                            <span class="value">{client.city || 'N/A'}</span>
                        </div>
                        <div class="detail">
                            <span class="label">⭐ Points</span>
                            <span class="value points">{client.loyalty_points || 0}</span>
                        </div>
                        <div class="detail">
                            <span class="label">📧 Newsletter</span>
                            <span class="value">{client.newsletter_opt_in ? 'Yes' : 'No'}</span>
                        </div>
                    </div>
                    
                    <div class="client-actions">
                        {#if client.account_status === 'PendingDeletion'}
                            <button class="btn-restore" on:click={() => openDeleteModal(client, 'cancel')}>
                                Restore Account
                            </button>
                            <button class="btn-delete" on:click={() => openDeleteModal(client, 'confirm')}>
                                Confirm Delete
                            </button>
                        {:else}
                            <button class="btn-edit" on:click={() => openEditModal(client)}>
                                Edit Profile
                            </button>
                        {/if}
                    </div>
                </div>
            {:else}
                <div class="empty-state">
                    <p>No clients found</p>
                </div>
            {/each}
        </div>
    {/if}
</div>

{#if showModal && selectedClient}
    <div class="modal-overlay" on:click={closeModal} on:keypress={closeModal}>
        <div class="modal" on:click|stopPropagation on:keypress|stopPropagation>
            {#if modalType === 'edit'}
                <h2>Edit Client Profile</h2>
                <p class="modal-subtitle">{selectedClient.first_name} {selectedClient.last_name}</p>
                
                <div class="form-group">
                    <label>Phone</label>
                    <input type="tel" bind:value={editData.phone} placeholder="Phone number" />
                </div>
                <div class="form-group">
                    <label>Address</label>
                    <input type="text" bind:value={editData.address} placeholder="Address" />
                </div>
                <div class="form-group">
                    <label>City</label>
                    <input type="text" bind:value={editData.city} placeholder="City" />
                </div>
                <div class="form-group">
                    <label>Gender</label>
                    <select bind:value={editData.gender}>
                        <option value="">Not specified</option>
                        <option value="Male">Male</option>
                        <option value="Female">Female</option>
                        <option value="Other">Other</option>
                    </select>
                </div>
                <div class="form-group checkbox">
                    <label>
                        <input type="checkbox" bind:checked={editData.newsletter_opt_in} />
                        Newsletter subscription
                    </label>
                </div>
                
                <div class="modal-actions">
                    <button class="btn-cancel" on:click={closeModal}>Cancel</button>
                    <button class="btn-submit" on:click={handleSave} disabled={processing}>
                        {processing ? 'Saving...' : 'Save Changes'}
                    </button>
                </div>
            {:else}
                <h2>{modalType === 'cancel' ? 'Restore Account?' : 'Confirm Deletion?'}</h2>
                
                {#if modalType === 'cancel'}
                    <p>This will restore the account for <strong>{selectedClient.first_name} {selectedClient.last_name}</strong> and cancel the deletion request.</p>
                {:else}
                    <p class="warning-text">⚠️ This action cannot be undone. The account for <strong>{selectedClient.first_name} {selectedClient.last_name}</strong> will be permanently deleted and their data will be anonymized.</p>
                {/if}
                
                <div class="modal-actions">
                    <button class="btn-cancel" on:click={closeModal}>Cancel</button>
                    <button 
                        class={modalType === 'cancel' ? 'btn-submit' : 'btn-danger'} 
                        on:click={handleDeletionAction}
                        disabled={processing}
                    >
                        {processing ? 'Processing...' : (modalType === 'cancel' ? 'Restore Account' : 'Delete Permanently')}
                    </button>
                </div>
            {/if}
        </div>
    </div>
{/if}

<style>
    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 2rem;
    }
    
    header h1 {
        font-size: 2rem;
        margin: 0;
        color: #1a1a2e;
    }
    
    .subtitle {
        color: #666;
        margin: 0.5rem 0 2rem;
    }
    
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 1.5rem;
        margin-bottom: 2rem;
    }
    
    .stat-card {
        background: white;
        border-radius: 12px;
        padding: 1.5rem;
        display: flex;
        align-items: center;
        gap: 1rem;
        box-shadow: 0 2px 8px rgba(0,0,0,0.06);
    }
    
    .stat-card.warning {
        border: 2px solid #fbbf24;
    }
    
    .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 1.5rem;
    }
    
    .stat-icon.clients { background: #dbeafe; }
    .stat-icon.active { background: #d1fae5; }
    .stat-icon.deletion { background: #fef3c7; }
    .stat-icon.points { background: #fce7f3; }
    
    .stat-value {
        display: block;
        font-size: 1.5rem;
        font-weight: 700;
        color: #1a1a2e;
    }
    
    .stat-label {
        color: #666;
        font-size: 0.875rem;
    }
    
    .tabs {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
    }
    
    .tab {
        padding: 0.75rem 1.5rem;
        border: none;
        background: #f3f4f6;
        border-radius: 8px;
        cursor: pointer;
        font-weight: 500;
        color: #666;
        position: relative;
    }
    
    .tab.active {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
    }
    
    .badge {
        position: absolute;
        top: -8px;
        right: -8px;
        background: #ef4444;
        color: white;
        border-radius: 50%;
        width: 20px;
        height: 20px;
        font-size: 0.75rem;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    
    .search-bar {
        display: flex;
        gap: 1rem;
        margin-bottom: 1.5rem;
    }
    
    .search-bar input {
        flex: 1;
        padding: 0.75rem 1rem;
        border: 1px solid #e5e7eb;
        border-radius: 8px;
        font-size: 1rem;
    }
    
    .search-bar button {
        padding: 0.75rem 1.5rem;
        background: #1a1a2e;
        color: white;
        border: none;
        border-radius: 8px;
        cursor: pointer;
    }
    
    .clients-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
        gap: 1.5rem;
    }
    
    .client-card {
        background: white;
        border-radius: 12px;
        padding: 1.5rem;
        box-shadow: 0 2px 8px rgba(0,0,0,0.06);
    }
    
    .client-card.pending-deletion {
        border: 2px solid #fbbf24;
        background: #fffbeb;
    }
    
    .client-header {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 1rem;
    }
    
    .avatar {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 600;
    }
    
    .client-info {
        flex: 1;
    }
    
    .client-info h3 {
        margin: 0;
        font-size: 1.1rem;
    }
    
    .client-info p {
        margin: 0;
        color: #666;
        font-size: 0.875rem;
    }
    
    .status-badge {
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
    }
    
    .status-badge.active { background: #d1fae5; color: #059669; }
    .status-badge.pending { background: #fef3c7; color: #d97706; }
    .status-badge.deleted { background: #fee2e2; color: #dc2626; }
    
    .client-details {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 0.75rem;
        padding: 1rem 0;
        border-top: 1px solid #f3f4f6;
        border-bottom: 1px solid #f3f4f6;
    }
    
    .detail .label {
        display: block;
        font-size: 0.75rem;
        color: #666;
    }
    
    .detail .value {
        font-weight: 500;
    }
    
    .detail .value.points {
        color: #ec4899;
    }
    
    .client-actions {
        display: flex;
        gap: 0.5rem;
        margin-top: 1rem;
    }
    
    .btn-edit, .btn-restore, .btn-delete {
        flex: 1;
        padding: 0.5rem;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.875rem;
    }
    
    .btn-edit {
        background: #f3f4f6;
        color: #374151;
    }
    
    .btn-restore {
        background: #10b981;
        color: white;
    }
    
    .btn-delete {
        background: #ef4444;
        color: white;
    }
    
    .empty-state {
        grid-column: 1 / -1;
        text-align: center;
        padding: 3rem;
        color: #666;
    }
    
    .loading {
        text-align: center;
        padding: 3rem;
    }
    
    .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid #f3f4f6;
        border-top-color: #667eea;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin: 0 auto 1rem;
    }
    
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
    
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0,0,0,0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }
    
    .modal {
        background: white;
        padding: 2rem;
        border-radius: 12px;
        width: 450px;
        max-width: 90%;
    }
    
    .modal h2 {
        margin: 0 0 0.5rem;
    }
    
    .modal-subtitle {
        color: #666;
        margin: 0 0 1.5rem;
    }
    
    .form-group {
        margin-bottom: 1rem;
    }
    
    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    
    .form-group input, .form-group select {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 8px;
        font-size: 1rem;
    }
    
    .form-group.checkbox label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }
    
    .form-group.checkbox input {
        width: auto;
    }
    
    .warning-text {
        color: #dc2626;
        background: #fee2e2;
        padding: 1rem;
        border-radius: 8px;
    }
    
    .modal-actions {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        margin-top: 1.5rem;
    }
    
    .btn-cancel, .btn-submit, .btn-danger {
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 8px;
        cursor: pointer;
    }
    
    .btn-cancel {
        background: #f3f4f6;
    }
    
    .btn-submit {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
    }
    
    .btn-danger {
        background: #ef4444;
        color: white;
    }
    
    @media (max-width: 768px) {
        .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
</style>
