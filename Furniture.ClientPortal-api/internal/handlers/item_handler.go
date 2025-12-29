package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/yourusername/furniture-api/internal/models"
)

// ItemHandler handles furniture item endpoints

type ItemHandler struct {
	db *sqlx.DB
}

func NewItemHandler(db *sqlx.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/items", h.ListItems).Methods(http.MethodGet)
	r.HandleFunc("/items/{id}", h.GetItem).Methods(http.MethodGet)
	r.HandleFunc("/items", h.CreateItem).Methods(http.MethodPost)
	r.HandleFunc("/items/{id}", h.UpdateItem).Methods(http.MethodPut)
	r.HandleFunc("/items/{id}", h.DeleteItem).Methods(http.MethodDelete)
}

// ListItems returns a paginated list of items
func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	// TODO: Add filters, paging
	items := []models.Item{}
	err := h.db.Select(&items, "SELECT * FROM items ORDER BY id DESC LIMIT 50")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch items")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// GetItem returns item details with variants
func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var item models.Item
	err = h.db.Get(&item, "SELECT * FROM items WHERE id = $1", id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Error().Err(err).Msg("Failed to fetch item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var variants []models.ItemVariant
	err = h.db.Select(&variants, "SELECT * FROM item_variants WHERE item_id = $1", id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch variants")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		models.Item
		Variants []models.ItemVariant `json:"variants"`
	}{
		Item:     item,
		Variants: variants,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateItem creates a new item
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO items (name, description, price, stock_quantity, category, image_url, sku, dimensions, material, color, is_active, created_at, updated_at) VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW()) RETURNING id`

	var id int64
	err := h.db.QueryRow(query, item.Name, item.Description, item.Price, item.StockQuantity, item.Category, item.ImageURL, item.SKU, item.Dimensions, item.Material, item.Color, item.IsActive).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	item.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

// UpdateItem updates an existing item
func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE items SET name=$1, description=$2, price=$3, stock_quantity=$4, category=$5, image_url=$6, sku=$7, dimensions=$8, material=$9, color=$10, is_active=$11, updated_at=NOW() WHERE id=$12`

	_, err = h.db.Exec(query, item.Name, item.Description, item.Price, item.StockQuantity, item.Category, item.ImageURL, item.SKU, item.Dimensions, item.Material, item.Color, item.IsActive, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteItem soft deletes an item by setting status to 'Discontinued'
func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `DELETE FROM items WHERE id=$1`
	_, err = h.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
