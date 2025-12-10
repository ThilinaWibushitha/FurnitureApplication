<script>
    import { onMount } from 'svelte';
    import { getOrders, rejectPayment, sendClientEmail } from '$lib/api.js';
    
    let orders = [];
    let loading = true;
    let searchQuery = '';
    let statusFilter = '';
    let user = null;
    
    let showLocationModal = false;
    let showRejectModal = false;
    let showEmailModal = false;
    let selectedOrder = null;
    let rejectReason = '';
    let emailSubject = '';
    let emailBody = '';
    let processing = false;
    let successMessage = '';
    let errorMessage = '';
    
    onMount(async () => {
        const userStr = localStorage.getItem('user');
        if (userStr) {
            user = JSON.parse(userStr);
        }
        await loadData();
    });
    
    async function loadData() {
        loading = true;
        try {
            const filter = {};
            if (searchQuery) filter.search = searchQuery;
            if (statusFilter) filter.status = statusFilter;
            orders = await getOrders(filter);
        } catch (e) {
            console.error('Failed to load orders:', e);
            orders = [];
        }
        loading = false;
    }
    
    function handleSearch() {
        loadData();
    }
    
    function openLocationModal(order) {
        selectedOrder = order;
        showLocationModal = true;
    }
    
    function openRejectModal(order) {
        selectedOrder = order;
        rejectReason = '';
        showRejectModal = true;
    }
    
    function openEmailModal(order) {
        selectedOrder = order;
        emailSubject = '';
        emailBody = '';
        showEmailModal = true;
    }
    
    function closeModals() {
        showLocationModal = false;
        showRejectModal = false;
        showEmailModal = false;
        selectedOrder = null;
    }
    
    async function handleReject() {
        if (!selectedOrder) return;
        processing = true;
        errorMessage = '';
        try {
            await rejectPayment(selectedOrder.payment_id, {
                order_id: selectedOrder.id,
                reason: rejectReason
            });
            successMessage = 'Payment rejected and order cancelled successfully.';
            closeModals();
            await loadData();
        } catch (e) {
            console.error('Failed to reject:', e);
            errorMessage = 'Failed to reject payment. Please try again.';
        }
        processing = false;
    }
    
    async function handleSendEmail() {
        if (!selectedOrder || !emailSubject || !emailBody) return;
        processing = true;
        errorMessage = '';
        try {
            await sendClientEmail(selectedOrder.customer_id, {
                subject: emailSubject,
                body: emailBody
            });
            successMessage = 'Email sent successfully.';
            closeModals();
        } catch (e) {
            console.error('Failed to send email:', e);
            errorMessage = 'Failed to send email. Please try again.';
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
    
    function formatCurrency(amount) {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(amount ?? 0);
    }
    
    $: isMainAdmin = user?.account_type === 'main_admin';
</script>

<div class="container">
    <header>
        <h1>Orders Management</h1>
        <p class="subtitle">View orders, check delivery locations, and manage payments</p>
    </header>
    
    {#if successMessage}
        <div class="alert success">{successMessage}</div>
    {/if}
    
    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-icon orders">Package</div>
            <div class="stat-info">
                <span class="stat-value">{orders.length}</span>
                <span class="stat-label">Total Orders</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon pending">Clock</div>
            <div class="stat-info">
                <span class="stat-value">{orders.filter(o => o.order_status === 'Pending').length}</span>
                <span class="stat-label">Pending</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon delivered">Check</div>
            <div class="stat-info">
                <span class="stat-value">{orders.filter(o => o.order_status === 'Delivered').length}</span>
                <span class="stat-label">Delivered</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon revenue">Dollar</div>
            <div class="stat-info">
                <span class="stat-value">{formatCurrency(orders.reduce((sum, o) => sum + (o.total_amount || 0), 0))}</span>
                <span class="stat-label">Total Value</span>
            </div>
        </div>
    </div>
    
    <div class="toolbar">
        <div class="search-bar">
            <input 
                type="text" 
                bind:value={searchQuery} 
                placeholder="Search by order ID or customer..."
                on:keyup={(e) => e.key === 'Enter' && handleSearch()}
            />
            <button on:click={handleSearch}>Search</button>
        </div>
        <div class="filter-pills">
            <button class:active={statusFilter === ''} on:click={() => { statusFilter = ''; loadData(); }}>All</button>
            <button class:active={statusFilter === 'Pending'} on:click={() => { statusFilter = 'Pending'; loadData(); }}>Pending</button>
            <button class:active={statusFilter === 'Processing'} on:click={() => { statusFilter = 'Processing'; loadData(); }}>Processing</button>
            <button class:active={statusFilter === 'Shipped'} on:click={() => { statusFilter = 'Shipped'; loadData(); }}>Shipped</button>
            <button class:active={statusFilter === 'Delivered'} on:click={() => { statusFilter = 'Delivered'; loadData(); }}>Delivered</button>
        </div>
    </div>
    
    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading orders...</p>
        </div>
    {:else}
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>Order ID</th>
                        <th>Customer</th>
                        <th>Items</th>
                        <th>Total</th>
                        <th>Status</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each orders as order}
                        <tr>
                            <td>#{order.id}</td>
                            <td>
                                <div class="customer-info">
                                    <span class="customer-name">{order.customer_name || 'Unknown'}</span>
                                    <span class="customer-email">{order.customer_email || ''}</span>
                                </div>
                            </td>
                            <td>{order.item_count || 0} items</td>
                            <td class="amount">{formatCurrency(order.total_amount)}</td>
                            <td>
                                <span class="status-badge {order.order_status?.toLowerCase()}">
                                    {order.order_status}
                                </span>
                            </td>
                            <td>{formatDate(order.created_at)}</td>
                            <td class="actions">
                                <button class="btn-action" on:click={() => openLocationModal(order)} title="View Delivery Location">
                                    Location
                                </button>
                                <button class="btn-email" on:click={() => openEmailModal(order)} title="Send Email">
                                    Email
                                </button>
                                {#if isMainAdmin && order.order_status !== 'Cancelled' && order.order_status !== 'Delivered'}
                                    <button class="btn-reject" on:click={() => openRejectModal(order)} title="Reject Payment">
                                        Reject
                                    </button>
                                {/if}
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="7" class="empty">No orders found</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if showLocationModal && selectedOrder}
    <div class="modal-overlay" on:click={closeModals} on:keypress={closeModals}>
        <div class="modal" on:click|stopPropagation on:keypress|stopPropagation>
            <h2>Delivery Location</h2>
            <p class="modal-subtitle">Order #{selectedOrder.id}</p>
            
            <div class="location-details">
                <div class="detail-row">
                    <span class="label">Customer Name</span>
                    <span class="value">{selectedOrder.customer_name || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Delivery Address</span>
                    <span class="value">{selectedOrder.delivery_address || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">City</span>
                    <span class="value">{selectedOrder.delivery_city || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">State/Province</span>
                    <span class="value">{selectedOrder.delivery_state || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Postal Code</span>
                    <span class="value">{selectedOrder.delivery_postal_code || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Phone</span>
                    <span class="value">{selectedOrder.delivery_phone || 'N/A'}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Delivery Notes</span>
                    <span class="value">{selectedOrder.delivery_notes || 'None'}</span>
                </div>
            </div>
            
            <div class="modal-actions">
                <button class="btn-close" on:click={closeModals}>Close</button>
            </div>
        </div>
    </div>
{/if}

{#if showRejectModal && selectedOrder}
    <div class="modal-overlay" on:click={closeModals} on:keypress={closeModals}>
        <div class="modal" on:click|stopPropagation on:keypress|stopPropagation>
            <h2>Reject Payment & Cancel Order</h2>
            <p class="modal-subtitle warning">This will reject the payment and cancel Order #{selectedOrder.id}</p>
            
            {#if errorMessage}
                <div class="alert error">{errorMessage}</div>
            {/if}
            
            <div class="form-group">
                <label for="reject-reason">Reason for Rejection</label>
                <textarea 
                    id="reject-reason" 
                    bind:value={rejectReason} 
                    placeholder="Enter the reason for rejecting this payment..."
                    rows="4"
                ></textarea>
            </div>
            
            <div class="modal-actions">
                <button class="btn-cancel" on:click={closeModals}>Cancel</button>
                <button class="btn-danger" on:click={handleReject} disabled={processing || !rejectReason}>
                    {processing ? 'Processing...' : 'Reject Payment'}
                </button>
            </div>
        </div>
    </div>
{/if}

{#if showEmailModal && selectedOrder}
    <div class="modal-overlay" on:click={closeModals} on:keypress={closeModals}>
        <div class="modal" on:click|stopPropagation on:keypress|stopPropagation>
            <h2>Send Email to Customer</h2>
            <p class="modal-subtitle">Order #{selectedOrder.id} - {selectedOrder.customer_name}</p>
            
            {#if errorMessage}
                <div class="alert error">{errorMessage}</div>
            {/if}
            
            <div class="form-group">
                <label for="email-subject">Subject</label>
                <input 
                    type="text" 
                    id="email-subject" 
                    bind:value={emailSubject} 
                    placeholder="Email subject..."
                />
            </div>
            
            <div class="form-group">
                <label for="email-body">Message</label>
                <textarea 
                    id="email-body" 
                    bind:value={emailBody} 
                    placeholder="Enter your message..."
                    rows="6"
                ></textarea>
            </div>
            
            <div class="modal-actions">
                <button class="btn-cancel" on:click={closeModals}>Cancel</button>
                <button class="btn-submit" on:click={handleSendEmail} disabled={processing || !emailSubject || !emailBody}>
                    {processing ? 'Sending...' : 'Send Email'}
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
    
    .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 0.75rem;
        font-weight: 600;
    }
    
    .stat-icon.orders { background: #dbeafe; color: #1e40af; }
    .stat-icon.pending { background: #fef3c7; color: #d97706; }
    .stat-icon.delivered { background: #d1fae5; color: #059669; }
    .stat-icon.revenue { background: #e0e7ff; color: #4f46e5; }
    
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
    
    .toolbar {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        margin-bottom: 1.5rem;
        align-items: center;
    }
    
    .search-bar {
        display: flex;
        gap: 0.5rem;
        flex: 1;
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
    
    .filter-pills {
        display: flex;
        gap: 0.5rem;
    }
    
    .filter-pills button {
        padding: 0.5rem 1rem;
        border: 1px solid #e5e7eb;
        background: white;
        border-radius: 20px;
        cursor: pointer;
        font-size: 0.875rem;
    }
    
    .filter-pills button.active {
        background: #1a1a2e;
        color: white;
        border-color: #1a1a2e;
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
    
    .customer-info {
        display: flex;
        flex-direction: column;
    }
    
    .customer-name {
        font-weight: 500;
    }
    
    .customer-email {
        font-size: 0.875rem;
        color: #666;
    }
    
    .amount {
        font-weight: 600;
        color: #059669;
    }
    
    .status-badge {
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
        text-transform: uppercase;
    }
    
    .status-badge.pending { background: #fef3c7; color: #d97706; }
    .status-badge.processing { background: #dbeafe; color: #2563eb; }
    .status-badge.shipped { background: #e0e7ff; color: #4f46e5; }
    .status-badge.delivered { background: #d1fae5; color: #059669; }
    .status-badge.cancelled { background: #fee2e2; color: #dc2626; }
    
    .actions {
        display: flex;
        gap: 0.5rem;
    }
    
    .btn-action, .btn-email, .btn-reject {
        padding: 0.4rem 0.75rem;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.75rem;
        font-weight: 500;
    }
    
    .btn-action {
        background: #f3f4f6;
        color: #374151;
    }
    
    .btn-email {
        background: #dbeafe;
        color: #2563eb;
    }
    
    .btn-reject {
        background: #fee2e2;
        color: #dc2626;
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
        width: 500px;
        max-width: 90%;
        max-height: 90vh;
        overflow-y: auto;
    }
    
    .modal h2 {
        margin: 0 0 0.5rem;
    }
    
    .modal-subtitle {
        color: #666;
        margin: 0 0 1.5rem;
    }
    
    .modal-subtitle.warning {
        color: #dc2626;
    }
    
    .location-details {
        background: #f9fafb;
        border-radius: 8px;
        padding: 1rem;
        margin-bottom: 1.5rem;
    }
    
    .detail-row {
        display: flex;
        justify-content: space-between;
        padding: 0.5rem 0;
        border-bottom: 1px solid #e5e7eb;
    }
    
    .detail-row:last-child {
        border-bottom: none;
    }
    
    .detail-row .label {
        color: #666;
        font-size: 0.875rem;
    }
    
    .detail-row .value {
        font-weight: 500;
    }
    
    .form-group {
        margin-bottom: 1rem;
    }
    
    .form-group label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
    }
    
    .form-group input, .form-group textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 8px;
        font-size: 1rem;
        resize: vertical;
    }
    
    .modal-actions {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        margin-top: 1.5rem;
    }
    
    .btn-cancel, .btn-close, .btn-submit, .btn-danger {
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 8px;
        cursor: pointer;
        font-weight: 500;
    }
    
    .btn-cancel, .btn-close {
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
    
    @media (max-width: 768px) {
        .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        }
        
        .toolbar {
            flex-direction: column;
        }
        
        .search-bar {
            width: 100%;
        }
        
        .filter-pills {
            flex-wrap: wrap;
        }
    }
</style>
