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

type CustomerHandler struct {
	db *sqlx.DB
}

func NewCustomerHandler(db *sqlx.DB) *CustomerHandler {
	return &CustomerHandler{db: db}
}

func (h *CustomerHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/customers", h.ListCustomers).Methods(http.MethodGet)
	r.HandleFunc("/customers", h.CreateCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}", h.UpdateCustomer).Methods(http.MethodPut)
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	customers := []models.Customer{}
	err := h.db.Select(&customers, "SELECT * FROM customers ORDER BY last_name, first_name")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch customers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO customers (first_name, last_name, email, phone, address, city, loyalty_points, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,0,NOW(),NOW()) RETURNING id`
	var id int64
	err := h.db.QueryRow(query, c.FirstName, c.LastName, c.Email, c.Phone, c.Address, c.City).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create customer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE customers SET first_name=$1, last_name=$2, email=$3, phone=$4, address=$5, city=$6, updated_at=NOW() WHERE id=$7`
	_, err = h.db.Exec(query, c.FirstName, c.LastName, c.Email, c.Phone, c.Address, c.City, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update customer")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
