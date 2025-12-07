<script>
    import { onMount } from "svelte";
    import {
        getPayments,
        getCancellations,
        processRefund,
        updatePayment,
        getDailySales,
        getMonthlySales,
    } from "$lib/api.js";

    let payments = [];
    let cancellations = [];
    let loading = true;
    let statusFilter = "";
    let activeTab = "payments";

    let showModal = false;
    let modalType = "";
    let selectedItem = null;
    let processStatus = "";
    let processNotes = "";
    let processing = false;
    let dailySales = null;
    let monthlySales = null;

    onMount(async () => {
        await loadData();
    });

    async function loadData() {
        loading = true;
        try {
            const filter = statusFilter ? { status: statusFilter } : {};
            const [
                paymentsData,
                cancellationsData,
                dailyReport,
                monthlyReport,
            ] = await Promise.all([
                getPayments(filter),
                getCancellations(),
                getDailySales().catch(() => null),
                getMonthlySales().catch(() => null),
            ]);

            payments = paymentsData;
            cancellations = cancellationsData;
            dailySales = dailyReport ?? {
                total_sales: 0,
                total_orders: 0,
                items_sold: 0,
            };
            monthlySales = monthlyReport ?? {
                total_sales: 0,
                total_orders: 0,
                items_sold: 0,
            };
        } catch (e) {
            console.error("Failed to load data:", e);
            if (!dailySales) {
                dailySales = { total_sales: 0, total_orders: 0, items_sold: 0 };
            }
            if (!monthlySales) {
                monthlySales = {
                    total_sales: 0,
                    total_orders: 0,
                    items_sold: 0,
                };
            }
        }
        loading = false;
    }

    function filterByStatus(status) {
        statusFilter = status;
        loadData();
    }

    function openUpdateModal(payment) {
        selectedItem = payment;
        modalType = "update";
        processStatus = payment.payment_status;
        showModal = true;
    }

    function openRefundModal(cancellation) {
        selectedItem = cancellation;
        modalType = "refund";
        processStatus = "Completed";
        processNotes = "";
        showModal = true;
    }

    function closeModal() {
        showModal = false;
        selectedItem = null;
        processStatus = "";
        processNotes = "";
    }

    async function handleUpdate() {
        processing = true;
        try {
            if (modalType === "update") {
                await updatePayment(selectedItem.id, {
                    payment_status: processStatus,
                    notes: processNotes,
                });
            } else {
                await processRefund(selectedItem.id, {
                    cancellation_id: selectedItem.id,
                    status: processStatus,
                    notes: processNotes,
                });
            }
            closeModal();
            await loadData();
        } catch (e) {
            console.error("Failed to process:", e);
        }
        processing = false;
    }

    function formatDate(dateStr) {
        return new Date(dateStr).toLocaleDateString("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        });
    }

    function formatCurrency(amount) {
        return new Intl.NumberFormat("en-US", {
            style: "currency",
            currency: "USD",
        }).format(amount);
    }

    function formatNumber(value) {
        return new Intl.NumberFormat("en-US").format(value ?? 0);
    }

    function handleOverlayKeydown(event) {
        if (["Escape", "Enter", " "].includes(event.key)) {
            event.preventDefault();
            closeModal();
        }
    }

    function handleOverlayClick(event) {
        if (event.target === event.currentTarget) {
            closeModal();
        }
    }
</script>

<div class="container">
    <header>
        <h1>Admin Dashboard</h1>
        <p class="subtitle">Manage payments, refunds, and store settings</p>
    </header>

    <!-- Stats Cards -->
    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-icon pending">⏳</div>
            <div class="stat-info">
                <span class="stat-value"
                    >{payments.filter((p) => p.payment_status === "Pending")
                        .length}</span
                >
                <span class="stat-label">Pending</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon completed">✓</div>
            <div class="stat-info">
                <span class="stat-value"
                    >{payments.filter((p) => p.payment_status === "Completed")
                        .length}</span
                >
                <span class="stat-label">Completed</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon refunds">↩</div>
            <div class="stat-info">
                <span class="stat-value">{cancellations.length}</span>
                <span class="stat-label">Pending Refunds</span>
            </div>
        </div>
        <div class="stat-card">
            <div class="stat-icon total">$</div>
            <div class="stat-info">
                <span class="stat-value"
                    >{formatCurrency(
                        payments
                            .filter((p) => p.payment_status === "Completed")
                            .reduce((sum, p) => sum + p.amount, 0),
                    )}</span
                >
                <span class="stat-label">Total Revenue</span>
            </div>
        </div>
        <div class="stat-card highlight">
            <div class="stat-icon daily">📅</div>
            <div class="stat-info">
                <div class="stat-heading">
                    <span class="stat-value"
                        >{formatCurrency(dailySales?.total_sales ?? 0)}</span
                    >
                    <span class="badge subtle">Today</span>
                </div>
                <div class="stat-meta">
                    <span
                        >{formatNumber(dailySales?.total_orders ?? 0)} orders</span
                    >
                    <span>•</span>
                    <span
                        >{formatNumber(dailySales?.items_sold ?? 0)} items</span
                    >
                </div>
            </div>
        </div>
        <div class="stat-card highlight">
            <div class="stat-icon monthly">🗓️</div>
            <div class="stat-info">
                <div class="stat-heading">
                    <span class="stat-value"
                        >{formatCurrency(monthlySales?.total_sales ?? 0)}</span
                    >
                    <span class="badge subtle">This month</span>
                </div>
                <div class="stat-meta">
                    <span
                        >{formatNumber(monthlySales?.total_orders ?? 0)} orders</span
                    >
                    <span>•</span>
                    <span
                        >{formatNumber(monthlySales?.items_sold ?? 0)} items</span
                    >
                </div>
            </div>
        </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
        <button
            class="tab"
            class:active={activeTab === "payments"}
            on:click={() => (activeTab = "payments")}
        >
            Payments
        </button>
        <button
            class="tab"
            class:active={activeTab === "refunds"}
            on:click={() => (activeTab = "refunds")}
        >
            Refund Requests
            {#if cancellations.length > 0}
                <span class="badge">{cancellations.length}</span>
            {/if}
        </button>
    </div>

    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading...</p>
        </div>
    {:else if activeTab === "payments"}
        <!-- Filter Pills -->
        <div class="filter-pills">
            <button
                class:active={statusFilter === ""}
                on:click={() => filterByStatus("")}>All</button
            >
            <button
                class:active={statusFilter === "Pending"}
                on:click={() => filterByStatus("Pending")}>Pending</button
            >
            <button
                class:active={statusFilter === "Completed"}
                on:click={() => filterByStatus("Completed")}>Completed</button
            >
            <button
                class:active={statusFilter === "Cancelled"}
                on:click={() => filterByStatus("Cancelled")}>Cancelled</button
            >
            <button
                class:active={statusFilter === "Refunded"}
                on:click={() => filterByStatus("Refunded")}>Refunded</button
            >
        </div>

        <!-- Payments Table -->
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Order</th>
                        <th>Amount</th>
                        <th>Method</th>
                        <th>Status</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each payments as payment}
                        <tr>
                            <td>#{payment.id}</td>
                            <td>Order #{payment.order_id}</td>
                            <td class="amount"
                                >{formatCurrency(payment.amount)}</td
                            >
                            <td>{payment.payment_method}</td>
                            <td>
                                <span
                                    class="status-badge {payment.payment_status.toLowerCase()}"
                                >
                                    {payment.payment_status}
                                </span>
                            </td>
                            <td>{formatDate(payment.created_at)}</td>
                            <td>
                                <button
                                    class="btn-action"
                                    on:click={() => openUpdateModal(payment)}
                                >
                                    Edit
                                </button>
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="7" class="empty">No payments found</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {:else}
        <!-- Refunds Table -->
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Order</th>
                        <th>Amount</th>
                        <th>Reason</th>
                        <th>Type</th>
                        <th>Date</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each cancellations as cancel}
                        <tr>
                            <td>#{cancel.id}</td>
                            <td>Order #{cancel.order_id}</td>
                            <td class="amount"
                                >{formatCurrency(cancel.refund_amount || 0)}</td
                            >
                            <td class="reason">{cancel.reason}</td>
                            <td>{cancel.cancellation_type}</td>
                            <td>{formatDate(cancel.created_at)}</td>
                            <td>
                                <button
                                    class="btn-approve"
                                    on:click={() => openRefundModal(cancel)}
                                >
                                    Process
                                </button>
                            </td>
                        </tr>
                    {:else}
                        <tr>
                            <td colspan="7" class="empty">No pending refunds</td
                            >
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

{#if showModal}
    <button
        type="button"
        class="modal-overlay"
        aria-label="Close dialog"
        on:click={handleOverlayClick}
        on:keydown={handleOverlayKeydown}
    >
        <div class="modal" role="dialog" aria-modal="true" tabindex="-1">
            <h2>
                {modalType === "update" ? "Update Payment" : "Process Refund"}
            </h2>

            <div class="form-group">
                <label for="process-status">Status</label>
                <select id="process-status" bind:value={processStatus}>
                    {#if modalType === "update"}
                        <option value="Pending">Pending</option>
                        <option value="Completed">Completed</option>
                        <option value="Failed">Failed</option>
                        <option value="Cancelled">Cancelled</option>
                    {:else}
                        <option value="Completed">Approve Refund</option>
                        <option value="Rejected">Reject Refund</option>
                    {/if}
                </select>
            </div>

            <div class="form-group">
                <label for="process-notes">Notes</label>
                <textarea
                    id="process-notes"
                    bind:value={processNotes}
                    placeholder="Add notes..."
                ></textarea>
            </div>

            <div class="modal-actions">
                <button class="btn-cancel" on:click={closeModal}>Cancel</button>
                <button
                    class="btn-submit"
                    on:click={handleUpdate}
                    disabled={processing}
                >
                    {processing ? "Processing..." : "Save"}
                </button>
            </div>
        </div>
    </button>
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
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
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

    .stat-icon.pending {
        background: #fef3c7;
    }
    .stat-icon.completed {
        background: #d1fae5;
    }
    .stat-icon.refunds {
        background: #fee2e2;
    }
    .stat-icon.total {
        background: #dbeafe;
    }

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
        transition: all 0.2s;
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

    .filter-pills {
        display: flex;
        gap: 0.5rem;
        margin-bottom: 1rem;
    }

    .filter-pills button {
        padding: 0.5rem 1rem;
        border: 1px solid #e5e7eb;
        background: white;
        border-radius: 20px;
        cursor: pointer;
        font-size: 0.875rem;
        transition: all 0.2s;
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
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    th,
    td {
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

    .amount {
        font-weight: 600;
        color: #059669;
    }

    .reason {
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .status-badge {
        padding: 0.25rem 0.75rem;
        border-radius: 20px;
        font-size: 0.75rem;
        font-weight: 500;
        text-transform: uppercase;
    }

    .status-badge.pending {
        background: #fef3c7;
        color: #d97706;
    }
    .status-badge.completed {
        background: #d1fae5;
        color: #059669;
    }
    .status-badge.cancelled {
        background: #fee2e2;
        color: #dc2626;
    }
    .status-badge.refunded {
        background: #dbeafe;
        color: #2563eb;
    }
    .status-badge.failed {
        background: #fecaca;
        color: #b91c1c;
    }

    .btn-action,
    .btn-approve {
        padding: 0.5rem 1rem;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        font-size: 0.875rem;
    }

    .btn-action {
        background: #f3f4f6;
        color: #374151;
    }

    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
    }

    .modal {
        background: white;
        padding: 2rem;
        border-radius: 12px;
        width: 400px;
        max-width: 90%;
    }

    .modal h2 {
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

    .form-group select,
    .form-group textarea {
        width: 100%;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 8px;
        font-size: 1rem;
    }

    .form-group textarea {
        min-height: 80px;
        resize: vertical;
    }

    .modal-actions {
        display: flex;
        gap: 1rem;
        justify-content: flex-end;
        margin-top: 1.5rem;
    }

    .btn-cancel {
        padding: 0.75rem 1.5rem;
        border: none;
        background: #f3f4f6;
        border-radius: 8px;
        cursor: pointer;
    }

    .btn-submit {
        padding: 0.75rem 1.5rem;
        border: none;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        border-radius: 8px;
        cursor: pointer;
    }

    .btn-submit:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    @media (max-width: 768px) {
        .stats-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }
</style>
