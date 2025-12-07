package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/furniture-api/internal/models"
)

type OrderHandler struct {
	db *sqlx.DB
}

func NewOrderHandler(db *sqlx.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

func (h *OrderHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/orders", h.ListOrders).Methods(http.MethodGet)
	r.HandleFunc("/orders/{id}", h.GetOrder).Methods(http.MethodGet)
	r.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders := []models.Order{}
	err := h.db.Select(&orders, "SELECT * FROM orders ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch orders")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order models.Order
	err = h.db.Get(&order, "SELECT * FROM orders WHERE id = $1", id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var items []models.OrderItem
	err = h.db.Select(&items, "SELECT * FROM order_items WHERE order_id = $1", id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch order items")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		models.Order
		Items []models.OrderItem `json:"items"`
	}{
		Order: order,
		Items: items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		models.Order
		Items []models.OrderItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert Order
	queryOrder := `INSERT INTO orders (order_number, customer_id, user_id, branch_id, status, total_amount, tax_amount, discount_amount, payment_method, notes, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW()) RETURNING id`
	var orderID int64
	err = tx.QueryRow(queryOrder, req.OrderNumber, req.CustomerID, req.UserID, req.BranchID, req.Status, req.TotalAmount, req.TaxAmount, req.DiscountAmount, req.PaymentMethod, req.Notes).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to insert order")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert Items
	queryItem := `INSERT INTO order_items (order_id, item_variant_id, quantity, unit_price, subtotal, discount, created_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())`
	for _, item := range req.Items {
		_, err = tx.Exec(queryItem, orderID, item.ItemVariantID, item.Quantity, item.UnitPrice, item.Subtotal, item.Discount)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("Failed to insert order item")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Order.ID = orderID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
