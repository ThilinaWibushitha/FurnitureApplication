<script>
    import { onMount } from 'svelte';
    import {
        getItems,
        createItem,
        updateItem,
        deleteItem
    } from '$lib/api.js';

    let items = [];
    let loading = true;
    let errorMessage = '';
    let searchTerm = '';
    let categoryFilter = 'all';

    let showFormModal = false;
    let showDeleteModal = false;
    let isEditing = false;
    let processing = false;
    let selectedItem = null;

    let formData = {
        name: '',
        description: '',
        price: '',
        stock_quantity: '',
        category: '',
        image_url: '',
        is_active: true
    };

    let categories = [];

    $: categories = Array.from(
        new Set(items.map((item) => item.category).filter(Boolean))
    );

    const categoryPalette = {
        Seating: '#50c9c3',
        Storage: '#667eea',
        Lighting: '#f97316',
        Decor: '#ec4899',
        Bedroom: '#14b8a6',
        Outdoor: '#7c3aed'
    };

    onMount(async () => {
        await loadItems();
    });

    async function loadItems() {
        loading = true;
        errorMessage = '';
        try {
            items = await getItems();
        } catch (error) {
            console.error('Failed to load items:', error);
            errorMessage = 'Unable to load inventory. Please try again later.';
        }
        loading = false;
    }

    function filteredItems() {
        return items
            .filter((item) => {
                const matchesSearch = `${item.name} ${item.description ?? ''}`
                    .toLowerCase()
                    .includes(searchTerm.trim().toLowerCase());
                const matchesCategory =
                    categoryFilter === 'all' || item.category === categoryFilter;
                return matchesSearch && matchesCategory;
            })
            .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    }

    function openCreateModal() {
        resetForm();
        isEditing = false;
        showFormModal = true;
    }

    function openEditModal(item) {
        selectedItem = item;
        formData = {
            name: item.name || '',
            description: item.description || '',
            price: item.price != null ? String(item.price) : '',
            stock_quantity: item.stock_quantity != null ? String(item.stock_quantity) : '',
            category: item.category || '',
            image_url: item.image_url || '',
            is_active: item.is_active
        };
        isEditing = true;
        showFormModal = true;
    }

    function openDeleteModal(item) {
        selectedItem = item;
        showDeleteModal = true;
    }

    function closeModals() {
        showFormModal = false;
        showDeleteModal = false;
        selectedItem = null;
        processing = false;
    }

    function resetForm() {
        formData = {
            name: '',
            description: '',
            price: '',
            stock_quantity: '',
            category: '',
            image_url: '',
            is_active: true
        };
    }

    async function handleSave() {
        processing = true;
        try {
            const name = formData.name.trim();
            const description = formData.description.trim();
            const price = parseFloat(formData.price);
            const stock = parseInt(formData.stock_quantity, 10);
            const category = formData.category.trim();
            const imageUrl = formData.image_url.trim();

            if (!name || Number.isNaN(price) || Number.isNaN(stock)) {
                throw new Error('Please fill all required fields with valid values.');
            }

            const basePayload = {
                name,
                description: description ? description : null,
                price,
                stock_quantity: stock,
                category: category || 'Uncategorized',
                image_url: imageUrl ? imageUrl : null
            };

            if (isEditing && selectedItem) {
                await updateItem(selectedItem.id, {
                    ...basePayload,
                    is_active: formData.is_active
                });
            } else {
                await createItem(basePayload);
            }

            await loadItems();
            closeModals();
        } catch (error) {
            console.error('Failed to save item:', error);
            errorMessage = error.message || 'Unable to save item. Please try again.';
        }
        processing = false;
    }

    async function handleDelete() {
        if (!selectedItem) return;
        processing = true;
        try {
            await deleteItem(selectedItem.id);
            await loadItems();
            closeModals();
        } catch (error) {
            console.error('Failed to delete item:', error);
            errorMessage = 'Unable to delete item. Please try again later.';
        }
        processing = false;
    }

    function formatCurrency(amount) {
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD'
        }).format(amount ?? 0);
    }

    function formatDate(dateStr) {
        if (!dateStr) return '—';
        return new Date(dateStr).toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function handleOverlayKeydown(event) {
        if (event.key === 'Escape' || event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            closeModals();
        }
    }

    function handleGlobalKeydown(event) {
        if ((showFormModal || showDeleteModal) && event.key === 'Escape') {
            closeModals();
        }
    }

    function handleOverlayClick(event) {
        if (event.target === event.currentTarget) {
            closeModals();
        }
    }
</script>

<svelte:head>
    <title>Inventory Management - Admin Portal</title>
</svelte:head>

<svelte:window on:keydown={handleGlobalKeydown} />

<div class="container">
    <header>
        <div>
            <h1>🪑 Inventory Management</h1>
            <p class="subtitle">Add, update, and curate your furniture catalogue</p>
        </div>
        <button class="primary" on:click={openCreateModal}>
            <span>＋</span>
            <span>Add Item</span>
        </button>
    </header>

    <div class="toolbar">
        <div class="search">
            <input
                type="search"
                placeholder="Search by name or description..."
                bind:value={searchTerm}
            />
        </div>
        <div class="filters">
            <label>
                Category
                <select bind:value={categoryFilter}>
                    <option value="all">All</option>
                    {#each categories as category}
                        <option value={category}>{category}</option>
                    {/each}
                </select>
            </label>
            <button class="ghost" on:click={loadItems}>
                ↻ Refresh
            </button>
        </div>
    </div>

    {#if errorMessage}
        <div class="error-banner">{errorMessage}</div>
    {/if}

    {#if loading}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading inventory...</p>
        </div>
    {:else if filteredItems().length === 0}
        <div class="empty-state">
            <div class="empty-card">
                <span class="emoji">📦</span>
                <h2>No items found</h2>
                <p>Try adjusting your filters or add a new item to get started.</p>
                <button class="primary" on:click={openCreateModal}>Add New Item</button>
            </div>
        </div>
    {:else}
        <div class="grid">
            {#each filteredItems() as item}
                <div class="item-card" class:inactive={!item.is_active}>
                    <div class="item-header">
                        <div class="badge" style={`background:${categoryPalette[item.category] || '#334155'}`}
                            >{item.category}</div
                        >
                        <div class={`status ${item.is_active ? 'active' : 'inactive'}`}>
                            {item.is_active ? 'Active' : 'Inactive'}
                        </div>
                    </div>

                    <div class="item-body">
                        <h3>{item.name}</h3>
                        <p>{item.description || 'No description provided.'}</p>
                    </div>

                    {#if item.image_url}
                        <div class="preview">
                            <img src={item.image_url} alt={`Preview of ${item.name}`} loading="lazy" />
                        </div>
                    {/if}

                    <div class="item-meta">
                        <div>
                            <span class="label">Price</span>
                            <span class="value">{formatCurrency(item.price)}</span>
                        </div>
                        <div>
                            <span class="label">Stock</span>
                            <span class="value">{item.stock_quantity}</span>
                        </div>
                        <div>
                            <span class="label">Updated</span>
                            <span class="value">{formatDate(item.updated_at)}</span>
                        </div>
                    </div>

                    <div class="actions">
                        <button class="ghost" on:click={() => openEditModal(item)}>Edit</button>
                        <button class="danger" on:click={() => openDeleteModal(item)}>Delete</button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

{#if showFormModal}
    <button
        type="button"
        class="modal-overlay"
        aria-label="Close modal"
        on:click={handleOverlayClick}
        on:keydown={handleOverlayKeydown}
    >
        <div class="modal" role="dialog" aria-modal="true">
            <h2>{isEditing ? 'Edit Item' : 'Add Item'}</h2>
            <p class="modal-subtitle">
                {isEditing
                    ? 'Update the selected catalogue item.'
                    : 'Provide the details for the new inventory item.'}
            </p>

            <div class="form-grid">
                <label>
                    Name*
                    <input type="text" bind:value={formData.name} placeholder="Item name" required />
                </label>
                <label>
                    Category
                    <input type="text" bind:value={formData.category} placeholder="e.g. Seating" />
                </label>
                <label>
                    Price*
                    <input type="number" step="0.01" min="0" bind:value={formData.price} placeholder="250" />
                </label>
                <label>
                    Stock Quantity*
                    <input type="number" min="0" bind:value={formData.stock_quantity} placeholder="10" />
                </label>
                <label class="full">
                    Description
                    <textarea rows="3" bind:value={formData.description} placeholder="Enter a concise description"></textarea>
                </label>
                <label class="full">
                    Image URL
                    <input type="url" bind:value={formData.image_url} placeholder="https://" />
                </label>
                <label class="toggle">
                    <input type="checkbox" bind:checked={formData.is_active} />
                    <span>Show item in storefront</span>
                </label>
            </div>

            <div class="modal-actions">
                <button class="ghost" type="button" on:click={closeModals}>Cancel</button>
                <button class="primary" type="button" on:click|preventDefault={handleSave} disabled={processing}>
                    {processing ? 'Saving...' : isEditing ? 'Save Changes' : 'Add Item'}
                </button>
            </div>
        </div>
    </button>
{/if}

{#if showDeleteModal && selectedItem}
    <button
        type="button"
        class="modal-overlay"
        aria-label="Close confirmation"
        on:click={handleOverlayClick}
        on:keydown={handleOverlayKeydown}
    >
        <div class="modal" role="dialog" aria-modal="true">
            <h2>Delete Item?</h2>
            <p class="modal-subtitle warning">
                This will permanently remove <strong>{selectedItem.name}</strong> from the catalogue.
            </p>
            <div class="modal-actions">
                <button class="ghost" type="button" on:click={closeModals}>Cancel</button>
                <button class="danger" type="button" on:click={handleDelete} disabled={processing}>
                    {processing ? 'Deleting...' : 'Delete'}
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
        color: #f8fafc;
    }

    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    header h1 {
        margin: 0;
        font-size: 2rem;
    }

    .subtitle {
        margin: 0.35rem 0 0;
        color: rgba(248, 250, 252, 0.65);
    }

    .toolbar {
        display: flex;
        flex-wrap: wrap;
        gap: 1rem;
        justify-content: space-between;
        margin-bottom: 1.5rem;
    }

    .search input {
        min-width: 260px;
        padding: 0.75rem 1rem;
        border-radius: 12px;
        border: 1px solid rgba(148, 163, 184, 0.3);
        background: rgba(15, 23, 42, 0.65);
        color: inherit;
        transition: border 0.2s ease;
    }

    .search input:focus {
        outline: none;
        border-color: #50c9c3;
        box-shadow: 0 0 0 3px rgba(80, 201, 195, 0.15);
    }

    .filters {
        display: flex;
        align-items: center;
        gap: 1rem;
    }

    .filters label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        color: rgba(248, 250, 252, 0.75);
    }

    select {
        padding: 0.5rem 0.75rem;
        border-radius: 10px;
        border: 1px solid rgba(148, 163, 184, 0.4);
        background: rgba(15, 23, 42, 0.65);
        color: inherit;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
    }

    .item-card {
        background: rgba(15, 23, 42, 0.85);
        border-radius: 20px;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
        border: 1px solid rgba(148, 163, 184, 0.08);
        box-shadow: 0 18px 45px rgba(15, 23, 42, 0.35);
        transition: transform 0.2s ease, border 0.2s ease;
    }

    .item-card:hover {
        transform: translateY(-4px);
        border-color: rgba(80, 201, 195, 0.3);
    }

    .item-card.inactive {
        opacity: 0.75;
    }

    .item-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }

    .badge {
        padding: 0.35rem 0.75rem;
        border-radius: 999px;
        font-size: 0.75rem;
        font-weight: 600;
        color: #0f172a;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .status {
        font-size: 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.4px;
        padding: 0.25rem 0.75rem;
        border-radius: 999px;
    }

    .status.active {
        background: rgba(80, 201, 195, 0.15);
        color: #34d399;
    }

    .status.inactive {
        background: rgba(239, 68, 68, 0.15);
        color: #f87171;
    }

    .item-body h3 {
        margin: 0;
        font-size: 1.2rem;
    }

    .item-body p {
        margin: 0;
        color: rgba(148, 163, 184, 0.9);
        font-size: 0.95rem;
    }

    .preview {
        border-radius: 16px;
        overflow: hidden;
        border: 1px solid rgba(148, 163, 184, 0.2);
        aspect-ratio: 4 / 3;
        background: rgba(30, 41, 59, 0.75);
    }

    .preview img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        display: block;
    }

    .item-meta {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 1rem;
        padding: 0.75rem 1rem;
        border-radius: 14px;
        background: rgba(30, 41, 59, 0.65);
        font-size: 0.85rem;
    }

    .label {
        display: block;
        color: rgba(148, 163, 184, 0.7);
        font-size: 0.75rem;
        margin-bottom: 0.35rem;
    }

    .value {
        font-weight: 600;
        color: #f8fafc;
    }

    .actions {
        display: flex;
        gap: 0.75rem;
        justify-content: flex-end;
    }

    button {
        border: none;
        cursor: pointer;
        font-weight: 600;
        border-radius: 12px;
        padding: 0.6rem 1.1rem;
        display: inline-flex;
        align-items: center;
        gap: 0.35rem;
        transition: transform 0.2s ease, box-shadow 0.2s ease, opacity 0.2s ease;
    }

    button:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .primary {
        background: linear-gradient(135deg, #50c9c3, #96deda);
        color: #0f172a;
        box-shadow: 0 12px 30px rgba(80, 201, 195, 0.35);
    }

    .primary:hover:not(:disabled) {
        transform: translateY(-2px);
    }

    .ghost {
        background: rgba(15, 23, 42, 0.7);
        color: rgba(248, 250, 252, 0.85);
        border: 1px solid rgba(148, 163, 184, 0.2);
    }

    .ghost:hover:not(:disabled) {
        border-color: rgba(148, 163, 184, 0.35);
    }

    .danger {
        background: rgba(239, 68, 68, 0.18);
        color: #f87171;
        border: 1px solid rgba(239, 68, 68, 0.35);
    }

    .danger:hover:not(:disabled) {
        background: rgba(239, 68, 68, 0.25);
    }

    .loading {
        text-align: center;
        padding: 4rem 0;
        color: rgba(248, 250, 252, 0.8);
    }

    .spinner {
        width: 48px;
        height: 48px;
        margin: 0 auto 1rem;
        border-radius: 50%;
        border: 4px solid rgba(148, 163, 184, 0.2);
        border-top-color: #50c9c3;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .empty-state {
        display: flex;
        justify-content: center;
        padding: 4rem 0;
    }

    .empty-card {
        background: rgba(15, 23, 42, 0.85);
        border-radius: 24px;
        padding: 3rem;
        text-align: center;
        max-width: 420px;
        box-shadow: 0 20px 45px rgba(15, 23, 42, 0.4);
    }

    .empty-card h2 {
        margin: 0.75rem 0;
    }

    .empty-card p {
        color: rgba(148, 163, 184, 0.85);
        margin-bottom: 1.5rem;
    }

    .emoji {
        font-size: 2.5rem;
    }

    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(2, 8, 23, 0.7);
        backdrop-filter: blur(10px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1.5rem;
    }

    .modal {
        background: #0f172a;
        border-radius: 24px;
        padding: 2rem;
        width: min(600px, 100%);
        color: #f8fafc;
        border: 1px solid rgba(148, 163, 184, 0.15);
        box-shadow: 0 25px 70px rgba(15, 23, 42, 0.55);
    }

    .modal h2 {
        margin: 0 0 0.5rem;
    }

    .modal-subtitle {
        margin: 0 0 1.5rem;
        color: rgba(148, 163, 184, 0.8);
    }

    .modal-subtitle.warning {
        color: #f87171;
    }

    .form-grid {
        display: grid;
        grid-template-columns: repeat(2, minmax(0, 1fr));
        gap: 1rem;
        margin-bottom: 1.5rem;
    }

    .form-grid label {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        font-size: 0.9rem;
        color: rgba(148, 163, 184, 0.85);
    }

    .form-grid input,
    .form-grid textarea {
        padding: 0.75rem 1rem;
        border-radius: 12px;
        border: 1px solid rgba(148, 163, 184, 0.25);
        background: rgba(15, 23, 42, 0.65);
        color: inherit;
        font-size: 0.95rem;
        resize: vertical;
    }

    .form-grid input:focus,
    .form-grid textarea:focus {
        outline: none;
        border-color: #50c9c3;
        box-shadow: 0 0 0 3px rgba(80, 201, 195, 0.15);
    }

    .form-grid label.full {
        grid-column: 1/-1;
    }

    .form-grid label.toggle {
        grid-column: 1/-1;
        flex-direction: row;
        align-items: center;
        gap: 0.75rem;
        background: rgba(15, 23, 42, 0.65);
        border-radius: 12px;
        padding: 0.75rem 1rem;
        border: 1px solid rgba(148, 163, 184, 0.15);
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.75rem;
    }

    .error-banner {
        margin-bottom: 1rem;
        padding: 1rem 1.25rem;
        border-radius: 14px;
        background: rgba(239, 68, 68, 0.18);
        border: 1px solid rgba(239, 68, 68, 0.35);
        color: #fecaca;
    }

    @media (max-width: 768px) {
        header {
            flex-direction: column;
            align-items: flex-start;
        }

        .toolbar {
            flex-direction: column;
            align-items: stretch;
        }

        .filters {
            flex-wrap: wrap;
        }

        .form-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
