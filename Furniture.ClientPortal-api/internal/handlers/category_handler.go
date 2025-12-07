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

// CategoryHandler handles departments and categories

type CategoryHandler struct {
	db *sqlx.DB
}

func NewCategoryHandler(db *sqlx.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/departments", h.ListDepartments).Methods(http.MethodGet)
	r.HandleFunc("/departments", h.CreateDepartment).Methods(http.MethodPost)
	r.HandleFunc("/departments/{id}", h.UpdateDepartment).Methods(http.MethodPut)
	r.HandleFunc("/departments/{id}", h.DeleteDepartment).Methods(http.MethodDelete)

	r.HandleFunc("/categories", h.ListCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories", h.CreateCategory).Methods(http.MethodPost)
	r.HandleFunc("/categories/{id}", h.UpdateCategory).Methods(http.MethodPut)
	r.HandleFunc("/categories/{id}", h.DeleteCategory).Methods(http.MethodDelete)
}

// Departments

func (h *CategoryHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	departments := []models.Department{}
	err := h.db.Select(&departments, "SELECT * FROM departments WHERE is_active = TRUE ORDER BY sort_order")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch departments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departments)
}

func (h *CategoryHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var d models.Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO departments (name, color, icon, sort_order, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,TRUE,NOW(),NOW()) RETURNING id`
	var id int64
	err := h.db.QueryRow(query, d.Name, d.Color, d.Icon, d.SortOrder).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create department")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	d.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(d)
}

func (h *CategoryHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d models.Department
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE departments SET name=$1, color=$2, icon=$3, sort_order=$4, is_active=$5, updated_at=NOW() WHERE id=$6`
	_, err = h.db.Exec(query, d.Name, d.Color, d.Icon, d.SortOrder, d.IsActive, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update department")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE departments SET is_active=FALSE, updated_at=NOW() WHERE id=$1`
	_, err = h.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete department")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Categories

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories := []models.Category{}
	err := h.db.Select(&categories, "SELECT * FROM categories WHERE is_active = TRUE ORDER BY sort_order")
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch categories")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `INSERT INTO categories (department_id, name, parent_category_id, sort_order, is_active, is_published_online, online_sort_order, created_at, updated_at) VALUES ($1,$2,$3,$4,TRUE,$5,$6,NOW(),NOW()) RETURNING id`
	var id int64
	err := h.db.QueryRow(query, c.DepartmentID, c.Name, c.ParentCategoryID, c.SortOrder, c.IsPublishedOnline, c.OnlineSortOrder).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create category")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c models.Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE categories SET department_id=$1, name=$2, parent_category_id=$3, sort_order=$4, is_active=$5, is_published_online=$6, online_sort_order=$7, updated_at=NOW() WHERE id=$8`
	_, err = h.db.Exec(query, c.DepartmentID, c.Name, c.ParentCategoryID, c.SortOrder, c.IsActive, c.IsPublishedOnline, c.OnlineSortOrder, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update category")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := `UPDATE categories SET is_active=FALSE, updated_at=NOW() WHERE id=$1`
	_, err = h.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete category")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
