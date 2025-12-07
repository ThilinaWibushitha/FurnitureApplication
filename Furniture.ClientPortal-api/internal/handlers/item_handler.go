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
	err := h.db.Select(&items, "SELECT * FROM items WHERE status != 'Discontinued' ORDER BY id DESC LIMIT 50")
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

	query := `INSERT INTO items (sku, name, description, department_id, category_id, main_image_url, status, base_price, base_cost, default_discount_percent, tax_class_id, is_published_online, online_sort_order, created_by, updated_by, created_at, updated_at) VALUES
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW(),NOW()) RETURNING id`

	var id int64
	err := h.db.QueryRow(query, item.SKU, item.Name, item.Description, item.DepartmentID, item.CategoryID, item.MainImageURL, item.Status, item.BasePrice, item.BaseCost, item.DefaultDiscountPercent, item.TaxClassID, item.IsPublishedOnline, item.OnlineSortOrder, item.CreatedBy, item.UpdatedBy).Scan(&id)
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

	query := `UPDATE items SET sku=$1, name=$2, description=$3, department_id=$4, category_id=$5, main_image_url=$6, status=$7, base_price=$8, base_cost=$9, default_discount_percent=$10, tax_class_id=$11, is_published_online=$12, online_sort_order=$13, updated_by=$14, updated_at=NOW() WHERE id=$15`

	_, err = h.db.Exec(query, item.SKU, item.Name, item.Description, item.DepartmentID, item.CategoryID, item.MainImageURL, item.Status, item.BasePrice, item.BaseCost, item.DefaultDiscountPercent, item.TaxClassID, item.IsPublishedOnline, item.OnlineSortOrder, item.UpdatedBy, id)
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

	query := `UPDATE items SET status='Discontinued', updated_at=NOW() WHERE id=$1`
	_, err = h.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete item")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
