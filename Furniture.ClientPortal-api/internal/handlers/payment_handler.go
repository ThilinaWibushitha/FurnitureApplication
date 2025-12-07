package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/furniture-api/internal/models"
)

type PaymentHandler struct {
	db *sqlx.DB
}

func NewPaymentHandler(db *sqlx.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

func (h *PaymentHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/payments", h.ListPayments).Methods(http.MethodGet)
	r.HandleFunc("/payments/{id}", h.GetPayment).Methods(http.MethodGet)
	r.HandleFunc("/payments", h.CreatePayment).Methods(http.MethodPost)
	r.HandleFunc("/payments/{id}", h.UpdatePayment).Methods(http.MethodPut)
	r.HandleFunc("/payments/{id}/cancel", h.CancelPayment).Methods(http.MethodPost)
	r.HandleFunc("/orders/{id}/payments", h.GetOrderPayments).Methods(http.MethodGet)
}

// ListPayments returns all payments with optional filters
func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	customerID := r.URL.Query().Get("customer_id")

	query := "SELECT * FROM payments WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += " AND payment_status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if customerID != "" {
		query += " AND customer_id = $" + strconv.Itoa(argIndex)
		args = append(args, customerID)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT 100"

	payments := []models.Payment{}
	err := h.db.Select(&payments, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch payments")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch payments"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// GetPayment returns a single payment by ID
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payment ID"})
		return
	}

	var payment models.Payment
	err = h.db.Get(&payment, "SELECT * FROM payments WHERE id = $1", id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to fetch payment")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// GetOrderPayments returns all payments for a specific order
func (h *PaymentHandler) GetOrderPayments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid order ID"})
		return
	}

	payments := []models.Payment{}
	err = h.db.Select(&payments, "SELECT * FROM payments WHERE order_id = $1 ORDER BY created_at DESC", orderID)
	if err != nil {
		log.Error().Err(err).Int64("order_id", orderID).Msg("Failed to fetch order payments")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch payments"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate order exists
	var orderExists bool
	err := h.db.Get(&orderExists, "SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)", req.OrderID)
	if err != nil || !orderExists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	// Get customer ID from order
	var customerID *int64
	h.db.Get(&customerID, "SELECT customer_id FROM orders WHERE id = $1", req.OrderID)

	query := `
		INSERT INTO payments (order_id, customer_id, amount, payment_method, payment_status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'Pending', $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var payment models.Payment
	payment.OrderID = req.OrderID
	payment.CustomerID = customerID
	payment.Amount = req.Amount
	payment.PaymentMethod = req.PaymentMethod
	payment.PaymentStatus = "Pending"
	payment.Notes = req.Notes

	err = h.db.QueryRow(query, req.OrderID, customerID, req.Amount, req.PaymentMethod, req.Notes).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create payment")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create payment"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// UpdatePayment updates a payment status
func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payment ID"})
		return
	}

	var req models.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.PaymentStatus != "" {
		updates = append(updates, "payment_status = $"+strconv.Itoa(argIndex))
		args = append(args, req.PaymentStatus)
		argIndex++

		if req.PaymentStatus == "Completed" {
			now := time.Now()
			updates = append(updates, "processed_at = $"+strconv.Itoa(argIndex))
			args = append(args, now)
			argIndex++
		}
	}

	if req.TransactionRef != "" {
		updates = append(updates, "transaction_ref = $"+strconv.Itoa(argIndex))
		args = append(args, req.TransactionRef)
		argIndex++
	}

	if req.Notes != "" {
		updates = append(updates, "notes = $"+strconv.Itoa(argIndex))
		args = append(args, req.Notes)
		argIndex++
	}

	if len(updates) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No fields to update"})
		return
	}

	updates = append(updates, "updated_at = NOW()")
	args = append(args, id)

	query := "UPDATE payments SET " + joinStrings(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)

	result, err := h.db.Exec(query, args...)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to update payment")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update payment"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
		return
	}

	// Fetch updated payment
	var payment models.Payment
	h.db.Get(&payment, "SELECT * FROM payments WHERE id = $1", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// CancelPayment cancels a payment and initiates refund
func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payment ID"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Start transaction
	tx, err := h.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get current payment
	var payment models.Payment
	err = tx.Get(&payment, "SELECT * FROM payments WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
		return
	}

	if payment.PaymentStatus == "Cancelled" || payment.PaymentStatus == "Refunded" {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment already cancelled or refunded"})
		return
	}

	// Update payment status
	_, err = tx.Exec("UPDATE payments SET payment_status = 'Cancelled', updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to cancel payment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create cancellation record
	_, err = tx.Exec(`
		INSERT INTO order_cancellations (order_id, payment_id, reason, cancellation_type, refund_amount, refund_status, created_at, updated_at)
		VALUES ($1, $2, $3, 'PaymentOnly', $4, 'Pending', NOW(), NOW())
	`, payment.OrderID, id, req.Reason, payment.Amount)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to create cancellation record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment cancelled successfully", "status": "Cancelled"})
}

// Helper function to join strings
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
