<script>
    import { onMount } from 'svelte';
    import { getPasswordChangeRequests, resolvePasswordChangeRequest } from '$lib/api.js';
    
    let requests = [];
    let loading = true;
    let user = null;
    
    let showModal = false;
    let selectedRequest = null;
    let action = '';
    let processing = false;
    let successMessage = '';
    let errorMessage = '';
    
    onMount(async () => {
        const userStr = localStorage.getItem('user');
        if (userStr) {
            user = JSON.parse(userStr);
        }
        
        if (user?.account_type !== 'main_admin') {
            window.location.href = '/';
            return;
        }
        
        await loadData();
    });
    
    async function loadData() {
        loading = true;
        try {
            requests = await getPasswordChangeRequests();
        } catch (e) {
            console.error('Failed to load password requests:', e);
            requests = [];
        }
        loading = false;
    }
    
    function openModal(request, actionType) {
        selectedRequest = request;
        action = actionType;
        showModal = true;
    }
    
    function closeModal() {
        showModal = false;
        selectedRequest = null;
        action = '';
    }
    
    async function handleResolve() {
        if (!selectedRequest) return;
        processing = true;
        errorMessage = '';
        try {
            await resolvePasswordChangeRequest(selectedRequest.id, action === 'approve');
            successMessage = action === 'approve' 
                ? 'Password change request approved successfully.'
                : 'Password change request rejected.';
            closeModal();
            await loadData();
        } catch (e) {
            console.error('Failed to resolve request:', e);
            errorMessage = 'Failed to process request. Please try again.';
        }
        processing = false;
    }
    
    function formatDate(dateStr) {
        if (!dateStr) return 'N/A';
        return new Date(dateStr).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }
</script>

<div class="container">
    <header>
        <h1>Password Change Requests</h1>
        <p class="subtitle">Review and approve client password change requests</p>
    </header>
    
    {#if successMessage}
        <div class="alert success">{successMessage}</div>
    {/if}
    
    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-icon pending">!</div>
            <div class="stat-info">
                <span class="stat-value">{requests.filter(r => r.status === 'Pending').length}</span>
                <span class="stat-label">Pending Requests</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon approved">Check</div>
            <div class="stat-info">
                <span class="stat-value">{requests.filter(r => r.status === 'Approved').length}</span>
                <span class="stat-label">Approved</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon rejected">X</div>
            <div class="stat-info">
                <span class="stat-value">{requests.filter(r => r.status === 'Rejected').length}</span>
                <span class="stat-label">Rejected</span>
            </div>
        </div>
    </div>
    
    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading requests...</p>
        </div>
    {:else}
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>Request ID</th>
                        <th>Customer</th>
                        <th>Email</th>
                        <th>Requested At</th>
                        <th>Status</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each requests as request}
                        <tr>
                            <td>#{request.id}</td>
                            <td>{request.customer_name || 'Unknown'}</td>
                            <td>{request.customer_email || 'N/A'}</td>
                            <td>{formatDate(request.created_at)}</td>
                            <td>
                                <span class="status-badge {request.status?.toLowerCase()}">
                                    {request.status}
                                </span>
                            </td>
                            <td class="actions">
                                {#if request.status === 'Pending'}
                                    <button class="btn-approve" on:click={() => openModal(request, 'approve')}>
                                        Approve
                                    </button>
                                    <button class="btn-reject" on:click={() => openModal(request, 'reject')}>
                                        Reject
                                    </button>
                                {:else}
                                    <span class="resolved">Resolved</span>
                                {/if}
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="6" class="empty">No password change requests</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if showModal && selectedRequest}
    <div class="modal-overlay" on:click={closeModal} on:keypress={closeModal}>
        <div class="modal" on:click|stopPropagation on:keypress|stopPropagation>
            <h2>{action === 'approve' ? 'Approve Password Change?' : 'Reject Password Change?'}</h2>
            
            {#if errorMessage}
                <div class="alert error">{errorMessage}</div>
            {/if}
            
            <div class="request-details">
                <p><strong>Customer:</strong> {selectedRequest.customer_name}</p>
                <p><strong>Email:</strong> {selectedRequest.customer_email}</p>
                <p><strong>Requested:</strong> {formatDate(selectedRequest.created_at)}</p>
            </div>
            
            {#if action === 'approve'}
                <p class="info-text">This will allow the customer to change their password.</p>
            {:else}
                <p class="warning-text">The customer's password will remain unchanged and they will be notified.</p>
            {/if}
            
            <div class="modal-actions">
                <button class="btn-cancel" on:click={closeModal}>Cancel</button>
                <button 
                    class={action === 'approve' ? 'btn-submit' : 'btn-danger'}
                    on:click={handleResolve}
                    disabled={processing}
                >
                    {processing ? 'Processing...' : (action === 'approve' ? 'Approve' : 'Reject')}
                </button>
            </div>
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
        color: #fff;
    }
    
    .subtitle {
        color: rgba(255,255,255,0.7);
        margin: 0.5rem 0 2rem;
    }
    
    .alert {
        padding: 1rem;
        border-radius: 8px;
        margin-bottom: 1.5rem;
    }
    
    .alert.success {
        background: #dcfce7;
        color: #166534;
    }
    
    .alert.error {
        background: #fee2e2;
        color: #991b1b;
    }
    
    .stats-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
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
    
    .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 1.25rem;
        font-weight: 600;
    }
    
    .stat-icon.pending { background: #fef3c7; color: #d97706; }
    .stat-icon.approved { background: #d1fae5; color: #059669; }
    .stat-icon.rejected { background: #fee2e2; color: #dc2626; }
    
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
    
    .table-container {
        background: white;
        border-radius: 12px;
        overflow: hidden;
        box-shadow: 0 2px 8px rgba(0,0,0,0.06);
    }
    
    table {
        width: 100%;
        border-collapse: collapse;
    }
    
    th, td {
        padding: 1rem 1.5rem;
        text-align: left;
    }
    
    th {
        background: #f9fafb;
        font-weight: 600;
        color: #374151;
        font-size: 0.875rem;
        text-transform: uppercase;
    }
    
    tr:not(:last-child) td {
        border-bottom: 1px solid #f3f4f6;
    }
    
    .status-badge {
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
        text-transform: uppercase;
    }
    
    .status-badge.pending { background: #fef3c7; color: #d97706; }
    .status-badge.approved { background: #d1fae5; color: #059669; }
    .status-badge.rejected { background: #fee2e2; color: #dc2626; }
    
    .actions {
        display: flex;
        gap: 0.5rem;
    }
    
    .btn-approve, .btn-reject {
        padding: 0.4rem 0.75rem;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.75rem;
        font-weight: 500;
    }
    
    .btn-approve {
        background: #d1fae5;
        color: #059669;
    }
    
    .btn-reject {
        background: #fee2e2;
        color: #dc2626;
    }
    
    .resolved {
        color: #666;
        font-size: 0.875rem;
    }
    
    .empty {
        text-align: center;
        color: #666;
        padding: 3rem !important;
    }
    
    .loading {
        text-align: center;
        padding: 3rem;
        color: #fff;
    }
    
    .spinner {
        width: 40px;
        height: 40px;
        border: 4px solid rgba(255,255,255,0.2);
        border-top-color: #50c9c3;
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
        margin: 0 0 1.5rem;
    }
    
    .request-details {
        background: #f9fafb;
        border-radius: 8px;
        padding: 1rem;
        margin-bottom: 1rem;
    }
    
    .request-details p {
        margin: 0.25rem 0;
    }
    
    .info-text {
        color: #059669;
        background: #d1fae5;
        padding: 0.75rem;
        border-radius: 8px;
    }
    
    .warning-text {
        color: #dc2626;
        background: #fee2e2;
        padding: 0.75rem;
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
        font-weight: 500;
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
    
    .btn-submit:disabled, .btn-danger:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }
</style>
