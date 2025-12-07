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

type RatingHandler struct {
	db *sqlx.DB
}

func NewRatingHandler(db *sqlx.DB) *RatingHandler {
	return &RatingHandler{db: db}
}

func (h *RatingHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/items/{id}/ratings", h.GetItemRatings).Methods(http.MethodGet)
	r.HandleFunc("/items/{id}/ratings", h.CreateRating).Methods(http.MethodPost)
	r.HandleFunc("/ratings/{id}", h.UpdateRating).Methods(http.MethodPut)
	r.HandleFunc("/ratings/{id}", h.DeleteRating).Methods(http.MethodDelete)
	r.HandleFunc("/customers/{id}/ratings", h.GetCustomerRatings).Methods(http.MethodGet)
	r.HandleFunc("/items/{id}/rating-summary", h.GetRatingSummary).Methods(http.MethodGet)
}

// RatingSummary represents aggregated rating data
type RatingSummary struct {
	ItemID        int64   `json:"item_id"`
	TotalRatings  int     `json:"total_ratings"`
	AverageRating float64 `json:"average_rating"`
	Rating5Count  int     `json:"rating_5_count"`
	Rating4Count  int     `json:"rating_4_count"`
	Rating3Count  int     `json:"rating_3_count"`
	Rating2Count  int     `json:"rating_2_count"`
	Rating1Count  int     `json:"rating_1_count"`
}

// GetItemRatings returns all approved ratings for an item
func (h *RatingHandler) GetItemRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	ratings := []models.ItemRating{}
	query := `
		SELECT r.*, c.first_name, c.last_name
		FROM item_ratings r
		LEFT JOIN customers c ON r.customer_id = c.id
		WHERE r.item_id = $1 AND r.is_approved = true
		ORDER BY r.created_at DESC
	`
	err = h.db.Select(&ratings, query, itemID)
	if err != nil {
		log.Error().Err(err).Int64("item_id", itemID).Msg("Failed to fetch ratings")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch ratings"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratings)
}

// GetRatingSummary returns rating statistics for an item
func (h *RatingHandler) GetRatingSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	var summary RatingSummary
	summary.ItemID = itemID

	query := `
		SELECT 
			COUNT(*) as total_ratings,
			COALESCE(AVG(rating)::numeric(3,2), 0) as average_rating,
			COUNT(*) FILTER (WHERE rating = 5) as rating_5_count,
			COUNT(*) FILTER (WHERE rating = 4) as rating_4_count,
			COUNT(*) FILTER (WHERE rating = 3) as rating_3_count,
			COUNT(*) FILTER (WHERE rating = 2) as rating_2_count,
			COUNT(*) FILTER (WHERE rating = 1) as rating_1_count
		FROM item_ratings
		WHERE item_id = $1 AND is_approved = true
	`

	err = h.db.QueryRow(query, itemID).Scan(
		&summary.TotalRatings,
		&summary.AverageRating,
		&summary.Rating5Count,
		&summary.Rating4Count,
		&summary.Rating3Count,
		&summary.Rating2Count,
		&summary.Rating1Count,
	)
	if err != nil {
		log.Error().Err(err).Int64("item_id", itemID).Msg("Failed to fetch rating summary")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch rating summary"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// CreateRating creates a new item rating
func (h *RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	itemID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid item ID"})
		return
	}

	var req models.CreateRatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate rating value
	if req.Rating < 1 || req.Rating > 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rating must be between 1 and 5"})
		return
	}

	// Get customer ID from auth context (simplified - should use JWT token)
	customerID := r.Header.Get("X-Customer-ID")
	if customerID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Customer ID required"})
		return
	}

	custID, _ := strconv.ParseInt(customerID, 10, 64)

	// Check if customer already rated this item
	var existingRating int
	err = h.db.Get(&existingRating, "SELECT COUNT(*) FROM item_ratings WHERE item_id = $1 AND customer_id = $2", itemID, custID)
	if err == nil && existingRating > 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "You have already rated this item"})
		return
	}

	// Check if this is a verified purchase
	var isVerified bool
	h.db.Get(&isVerified, `
		SELECT EXISTS(
			SELECT 1 FROM orders o
			JOIN order_items oi ON o.id = oi.order_id
			JOIN item_variants iv ON oi.item_variant_id = iv.id
			WHERE iv.item_id = $1 AND o.customer_id = $2 AND o.status = 'Completed'
		)
	`, itemID, custID)

	query := `
		INSERT INTO item_ratings (item_id, customer_id, order_id, rating, review_text, is_verified_purchase, is_approved, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, false, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var rating models.ItemRating
	rating.ItemID = itemID
	rating.CustomerID = custID
	rating.OrderID = req.OrderID
	rating.Rating = req.Rating
	rating.ReviewText = req.ReviewText
	rating.IsVerifiedPurchase = isVerified
	rating.IsApproved = false

	err = h.db.QueryRow(query, itemID, custID, req.OrderID, req.Rating, req.ReviewText, isVerified).
		Scan(&rating.ID, &rating.CreatedAt, &rating.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create rating")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create rating"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

// UpdateRating updates an existing rating
func (h *RatingHandler) UpdateRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid rating ID"})
		return
	}

	var req struct {
		Rating     int    `json:"rating"`
		ReviewText string `json:"review_text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate rating
	if req.Rating < 1 || req.Rating > 5 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rating must be between 1 and 5"})
		return
	}

	query := `
		UPDATE item_ratings 
		SET rating = $1, review_text = $2, is_approved = false, updated_at = NOW()
		WHERE id = $3
	`

	result, err := h.db.Exec(query, req.Rating, req.ReviewText, id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to update rating")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update rating"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rating not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Rating updated successfully"})
}

// DeleteRating deletes a rating
func (h *RatingHandler) DeleteRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid rating ID"})
		return
	}

	result, err := h.db.Exec("DELETE FROM item_ratings WHERE id = $1", id)
	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("Failed to delete rating")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete rating"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rating not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCustomerRatings returns all ratings by a customer
func (h *RatingHandler) GetCustomerRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	customerID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid customer ID"})
		return
	}

	ratings := []models.ItemRating{}
	query := `
		SELECT r.*, i.name as item_name
		FROM item_ratings r
		LEFT JOIN items i ON r.item_id = i.id
		WHERE r.customer_id = $1
		ORDER BY r.created_at DESC
	`
	err = h.db.Select(&ratings, query, customerID)
	if err != nil {
		log.Error().Err(err).Int64("customer_id", customerID).Msg("Failed to fetch customer ratings")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch ratings"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratings)
}
