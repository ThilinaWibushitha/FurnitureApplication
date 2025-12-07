const API_BASE = 'http://localhost:5200/api';

async function fetchApi(endpoint, options = {}) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        headers: {
            'Content-Type': 'application/json',
            ...options.headers
        },
        ...options
    });

    if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
}

// Payments API
export async function getPayments(filter = {}) {
    const params = new URLSearchParams();
    if (filter.status) params.append('status', filter.status);
    if (filter.customerId) params.append('customer_id', filter.customerId);
    const query = params.toString() ? `?${params.toString()}` : '';
    return fetchApi(`/payments${query}`);
}

export async function getPayment(id) {
    return fetchApi(`/payments/${id}`);
}

export async function createPayment(data) {
    return fetchApi('/payments', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

export async function updatePayment(id, data) {
    return fetchApi(`/payments/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

// Cancellations API
export async function getCancellations() {
    return fetchApi('/cancellations');
}

export async function processRefund(id, data) {
    return fetchApi(`/cancellations/${id}/process`, {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

// Clients API
export async function getClients(filter = {}) {
    const params = new URLSearchParams();
    if (filter.status) params.append('status', filter.status);
    if (filter.search) params.append('search', filter.search);
    const query = params.toString() ? `?${params.toString()}` : '';
    return fetchApi(`/clients${query}`);
}

export async function getClient(id) {
    return fetchApi(`/clients/${id}`);
}

export async function updateClient(id, data) {
    return fetchApi(`/clients/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

export async function getPendingDeletions() {
    return fetchApi('/clients/pending-deletions');
}

export async function cancelDeletion(id) {
    return fetchApi(`/clients/${id}/cancel-deletion`, { method: 'POST' });
}

export async function confirmDeletion(id) {
    return fetchApi(`/clients/${id}/confirm-deletion`, { method: 'POST' });
}

// Items API
export async function getItems() {
    return fetchApi('/items');
}

export async function getItem(id) {
    return fetchApi(`/items/${id}`);
}

export async function createItem(data) {
    return fetchApi('/items', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

export async function updateItem(id, data) {
    return fetchApi(`/items/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

export async function deleteItem(id) {
    return fetchApi(`/items/${id}`, {
        method: 'DELETE'
    });
}

// Reports API
export async function getDailySales() {
    return fetchApi('/reports/daily');
}

export async function getMonthlySales() {
    return fetchApi('/reports/monthly');
}
