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

// InventoryHandler handles branches and inventory

type InventoryHandler struct {
	db *sqlx.DB
}

func NewInventoryHandler(db *sqlx.DB) *InventoryHandler {
	return &InventoryHandler{db: db}
}

func (h *InventoryHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/branches", h.ListBranches).Methods(http.MethodGet)
	r.HandleFunc("/branches", h.CreateBranch).Methods(http.MethodPost)
	r.HandleFunc("/branches/{id}", h.UpdateBranch).Methods(http.MethodPut)
	r.HandleFunc("/branches/{id}", h.DeleteBranch).Methods(http.MethodDelete)

	r.HandleFunc("/inventory/levels", h.ListInventoryLevels).Methods(http.MethodGet)
	r.HandleFunc("/inventory/adjust", h.AdjustInventory).Methods(http.MethodPost)
}

func (h *InventoryHandler) ListBranches(w http.ResponseWriter, r *http.Request) {
	branches := []models.Branch{}
	err := h.db.Select(&branches, "SELECT * FROM branches WHERE is_active = TRUE ORDER BY name")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch branches")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(branches)
}

func (h *InventoryHandler) CreateBranch(w http.ResponseWriter, r *http.Request) {
	var b models.Branch
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO branches (code, name, address, phone, is_active) VALUES ($1,$2,$3,$4,TRUE) RETURNING id`
	var id int64
	err := h.db.QueryRow(query, b.Code, b.Name, b.Address, b.Phone).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create branch")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func (h *InventoryHandler) UpdateBranch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var b models.Branch
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE branches SET code=$1, name=$2, address=$3, phone=$4, is_active=$5 WHERE id=$6`
	_, err = h.db.Exec(query, b.Code, b.Name, b.Address, b.Phone, b.IsActive, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update branch")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) DeleteBranch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE branches SET is_active=FALSE WHERE id=$1`
	_, err = h.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete branch")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) ListInventoryLevels(w http.ResponseWriter, r *http.Request) {
	levels := []models.InventoryLevel{}
	// TODO: Add filters
	err := h.db.Select(&levels, "SELECT * FROM inventory_levels")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch inventory levels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(levels)
}

func (h *InventoryHandler) AdjustInventory(w http.ResponseWriter, r *http.Request) {
	var adj struct {
		BranchID       int64   `json:"branch_id"`
		ItemVariantID  int64   `json:"item_variant_id"`
		ChangeQuantity float64 `json:"change_quantity"`
		MovementType   string  `json:"movement_type"`
		Notes          string  `json:"notes"`
		UnitCost       float64 `json:"unit_cost"`
	}

	if err := json.NewDecoder(r.Body).Decode(&adj); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Transaction to update inventory_levels and insert inventory_movements
	tx, err := h.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update inventory_levels
	queryUpdate := `UPDATE inventory_levels SET quantity_on_hand = quantity_on_hand + $1 WHERE branch_id = $2 AND item_variant_id = $3`
	res, err := tx.Exec(queryUpdate, adj.ChangeQuantity, adj.BranchID, adj.ItemVariantID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to update inventory_levels")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to get rows affected")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If no row updated, insert new row
	if rowsAffected == 0 {
		queryInsert := `INSERT INTO inventory_levels (branch_id, item_variant_id, quantity_on_hand, quantity_reserved, reorder_level, created_at, updated_at) VALUES ($1,$2,$3,0,0,NOW(),NOW())`
		_, err = tx.Exec(queryInsert, adj.BranchID, adj.ItemVariantID, adj.ChangeQuantity)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("Failed to insert inventory_levels")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Insert inventory_movement
	queryMovement := `INSERT INTO inventory_movements (branch_id, item_variant_id, change_quantity, movement_type, notes, unit_cost, created_at) VALUES ($1,$2,$3,$4,$5,$6,NOW())`
	_, err = tx.Exec(queryMovement, adj.BranchID, adj.ItemVariantID, adj.ChangeQuantity, adj.MovementType, adj.Notes, adj.UnitCost)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to insert inventory_movement")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
