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
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	db *sqlx.DB
}

func NewProfileHandler(db *sqlx.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/profile", h.GetProfile).Methods(http.MethodGet)
	r.HandleFunc("/profile", h.UpdateProfile).Methods(http.MethodPut)
	r.HandleFunc("/profile/delete", h.RequestDeleteProfile).Methods(http.MethodPost)
	r.HandleFunc("/profile/password", h.ChangePassword).Methods(http.MethodPost)
	r.HandleFunc("/customers/{id}/profile", h.GetCustomerProfile).Methods(http.MethodGet)
}

// ClientProfileResponse combines customer and profile data
type ClientProfileResponse struct {
	CustomerID      int64   `json:"customer_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	Address         string  `json:"address"`
	City            string  `json:"city"`
	LoyaltyPoints   int     `json:"loyalty_points"`
	DateOfBirth     *string `json:"date_of_birth,omitempty"`
	Gender          string  `json:"gender,omitempty"`
	NewsletterOptIn bool    `json:"newsletter_opt_in"`
	ProfileImageUrl string  `json:"profile_image_url,omitempty"`
	AccountStatus   string  `json:"account_status"`
}

// GetProfile returns the current user's profile
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get customer ID from auth context
	customerID := r.Header.Get("X-Customer-ID")
	if customerID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	custID, _ := strconv.ParseInt(customerID, 10, 64)

	var profile ClientProfileResponse
	query := `
		SELECT 
			c.id as customer_id,
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.address,
			c.city,
			c.loyalty_points,
			cp.date_of_birth,
			cp.gender,
			COALESCE(cp.newsletter_opt_in, false) as newsletter_opt_in,
			cp.profile_image_url,
			COALESCE(cp.account_status, 'Active') as account_status
		FROM customers c
		LEFT JOIN client_profiles cp ON c.id = cp.customer_id
		WHERE c.id = $1
	`

	err := h.db.Get(&profile, query, custID)
	if err != nil {
		log.Error().Err(err).Int64("customer_id", custID).Msg("Failed to fetch profile")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Profile not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// GetCustomerProfile returns a specific customer's profile (for admin)
func (h *ProfileHandler) GetCustomerProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	custID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid customer ID"})
		return
	}

	var profile ClientProfileResponse
	query := `
		SELECT 
			c.id as customer_id,
			c.first_name,
			c.last_name,
			c.email,
			c.phone,
			c.address,
			c.city,
			c.loyalty_points,
			cp.date_of_birth,
			cp.gender,
			COALESCE(cp.newsletter_opt_in, false) as newsletter_opt_in,
			cp.profile_image_url,
			COALESCE(cp.account_status, 'Active') as account_status
		FROM customers c
		LEFT JOIN client_profiles cp ON c.id = cp.customer_id
		WHERE c.id = $1
	`

	err = h.db.Get(&profile, query, custID)
	if err != nil {
		log.Error().Err(err).Int64("customer_id", custID).Msg("Failed to fetch customer profile")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Customer not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get("X-Customer-ID")
	if customerID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	custID, _ := strconv.ParseInt(customerID, 10, 64)

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	tx, err := h.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update customer table
	if req.Phone != "" || req.Address != "" || req.City != "" {
		updates := []string{}
		args := []interface{}{}
		argIndex := 1

		if req.Phone != "" {
			updates = append(updates, "phone = $"+strconv.Itoa(argIndex))
			args = append(args, req.Phone)
			argIndex++
		}
		if req.Address != "" {
			updates = append(updates, "address = $"+strconv.Itoa(argIndex))
			args = append(args, req.Address)
			argIndex++
		}
		if req.City != "" {
			updates = append(updates, "city = $"+strconv.Itoa(argIndex))
			args = append(args, req.City)
			argIndex++
		}

		updates = append(updates, "updated_at = NOW()")
		args = append(args, custID)

		query := "UPDATE customers SET " + joinStrings(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex)
		_, err = tx.Exec(query, args...)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("Failed to update customer")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Upsert client profile
	profileQuery := `
		INSERT INTO client_profiles (customer_id, date_of_birth, gender, newsletter_opt_in, profile_image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (customer_id) DO UPDATE SET
			date_of_birth = COALESCE(EXCLUDED.date_of_birth, client_profiles.date_of_birth),
			gender = COALESCE(EXCLUDED.gender, client_profiles.gender),
			newsletter_opt_in = COALESCE(EXCLUDED.newsletter_opt_in, client_profiles.newsletter_opt_in),
			profile_image_url = COALESCE(EXCLUDED.profile_image_url, client_profiles.profile_image_url),
			updated_at = NOW()
	`

	var dob *string
	if req.DateOfBirth != "" {
		dob = &req.DateOfBirth
	}

	_, err = tx.Exec(profileQuery, custID, dob, req.Gender, req.NewsletterOptIn, req.ProfileImageUrl)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to update client profile")
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
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

// RequestDeleteProfile initiates account deletion process
func (h *ProfileHandler) RequestDeleteProfile(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get("X-Customer-ID")
	if customerID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	custID, _ := strconv.ParseInt(customerID, 10, 64)

	var req models.DeleteProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Verify password (get user linked to customer)
	var userID int64
	var passwordHash string
	err := h.db.Get(&userID, "SELECT user_id FROM client_profiles WHERE customer_id = $1", custID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Profile not found"})
		return
	}

	err = h.db.Get(&passwordHash, "SELECT password_hash FROM users WHERE id = $1", userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid password"})
		return
	}

	// Check for pending orders
	var pendingOrders int
	h.db.Get(&pendingOrders, "SELECT COUNT(*) FROM orders WHERE customer_id = $1 AND status IN ('Pending', 'Processing')", custID)
	if pendingOrders > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete account with pending orders"})
		return
	}

	// Mark for deletion (30 day grace period)
	now := time.Now()
	_, err = h.db.Exec(`
		UPDATE client_profiles 
		SET account_status = 'PendingDeletion', deletion_requested_at = $1, updated_at = NOW()
		WHERE customer_id = $2
	`, now, custID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to request deletion")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Deactivate user
	h.db.Exec("UPDATE users SET is_active = false WHERE id = $1", userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Account scheduled for deletion. You have 30 days to cancel this request.",
		"status":  "PendingDeletion",
	})
}

// ChangePassword changes the user's password
func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	customerID := r.Header.Get("X-Customer-ID")
	if customerID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	custID, _ := strconv.ParseInt(customerID, 10, 64)

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Get user ID
	var userID int64
	err := h.db.Get(&userID, "SELECT user_id FROM client_profiles WHERE customer_id = $1", custID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Profile not found"})
		return
	}

	// Get current password hash
	var passwordHash string
	err = h.db.Get(&passwordHash, "SELECT password_hash FROM users WHERE id = $1", userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.CurrentPassword))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update password
	_, err = h.db.Exec("UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2", string(newHash), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}
