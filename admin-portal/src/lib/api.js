const API_BASE = '/api';

function getAuthHeaders() {
    const token = typeof localStorage !== 'undefined' ? localStorage.getItem('token') : null;
    return token ? { 'Authorization': `Bearer ${token}` } : {};
}

async function fetchApi(endpoint, options = {}) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        headers: {
            'Content-Type': 'application/json',
            ...getAuthHeaders(),
            ...options.headers
        },
        ...options
    });

    if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
}

// Auth API
export async function login(email, password) {
    return fetchApi('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    });
}

export async function forgotPassword(email) {
    return fetchApi('/auth/forgot-password', {
        method: 'POST',
        body: JSON.stringify({ email })
    });
}

export async function resetPassword(token, newPassword) {
    return fetchApi('/auth/reset-password', {
        method: 'POST',
        body: JSON.stringify({ token, new_password: newPassword })
    });
}

export async function createAdmin(data) {
    return fetchApi('/auth/create-admin', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

// Users API
export async function getAdmins() {
    return fetchApi('/users/admins');
}

export async function getPasswordChangeRequests() {
    return fetchApi('/users/password-requests');
}

export async function resolvePasswordChangeRequest(requestId, approved) {
    return fetchApi('/users/password-requests/resolve', {
        method: 'POST',
        body: JSON.stringify({ request_id: requestId, approved })
    });
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

// Orders API
export async function getOrders(filter = {}) {
    const params = new URLSearchParams();
    if (filter.status) params.append('status', filter.status);
    if (filter.search) params.append('search', filter.search);
    const query = params.toString() ? `?${params.toString()}` : '';
    return fetchApi(`/orders${query}`);
}

export async function getOrder(id) {
    return fetchApi(`/orders/${id}`);
}

export async function rejectPayment(paymentId, data) {
    return fetchApi(`/payments/${paymentId}/reject`, {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

// Email API
export async function sendClientEmail(clientId, data) {
    return fetchApi(`/mail/client/${clientId}`, {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

export async function sendBulkEmail(data) {
    return fetchApi('/mail/bulk', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}
