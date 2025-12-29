package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yourusername/furniture-api/internal/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

// RegisterRoutes registers user-related routes
func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.CreateUser).Methods("POST")
	router.HandleFunc("/auth/login", h.Login).Methods("POST")
	router.HandleFunc("/users/{id}", h.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}/profile", h.UpdateProfile).Methods("PUT")
	router.HandleFunc("/users/username/{username}", h.GetUserByUsername).Methods("GET")
	router.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" || req.FullName == "" || req.AccountType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate account type
	if req.AccountType != "admin" && req.AccountType != "client" {
		http.Error(w, "Invalid account type", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert user
	var userID int64
	err = tx.QueryRow(`
		INSERT INTO users (username, email, password_hash, full_name, account_type, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, false, NOW(), NOW())
		RETURNING id
	`, req.Username, req.Email, string(hashedPassword), req.FullName, req.AccountType).Scan(&userID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Username or email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Insert profile if additional info provided
	if req.Phone != "" || req.Address != "" {
		_, err = tx.Exec(`
			INSERT INTO user_profiles (user_id, phone, address, city, state, postal_code, country, company_name, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		`, userID, req.Phone, req.Address, req.City, req.State, req.PostalCode, req.Country, req.CompanyName)

		if err != nil {
			log.Printf("Error creating user profile: %v", err)
			// Don't fail the whole operation if profile creation fails
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get created user with profile
	user, err := h.getUserByID(userID)
	if err != nil {
		log.Printf("Error retrieving created user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Login authenticates a user
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.User
	var passwordHash string
	err := h.db.QueryRow(`
		SELECT id, username, email, password_hash, full_name, account_type, is_active, is_verified, last_login_at, created_at, updated_at
		FROM users 
		WHERE username = $1 OR email = $1
	`, req.Username).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, &user.FullName,
		&user.AccountType, &user.IsActive, &user.IsVerified, &user.LastLogin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Check if user is active
	if !user.IsActive {
		http.Error(w, "Account is deactivated", http.StatusUnauthorized)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Update last login
	now := time.Now()
	_, err = h.db.Exec("UPDATE users SET last_login_at = $1 WHERE id = $2", now, user.ID)
	if err != nil {
		log.Printf("Error updating last login: %v", err)
	}

	user.LastLogin = &now

	// Get user profile
	profile, err := h.getUserProfile(user.ID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
	}

	// Create response structure matching frontend expectation
	authResponse := struct {
		Token string              `json:"token"`
		User  models.UserResponse `json:"user"`
	}{
		Token: "mock-jwt-token",
		User: models.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			FullName:    user.FullName,
			AccountType: user.AccountType,
			IsActive:    user.IsActive,
			IsVerified:  user.IsVerified,
			CreatedAt:   user.CreatedAt,
			LastLogin:   user.LastLogin,
			Profile:     profile,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.getUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByUsername retrieves a user by username
func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var user models.User
	err := h.db.QueryRow(`
		SELECT id, username, email, full_name, account_type, is_active, is_verified, last_login_at, created_at, updated_at
		FROM users 
		WHERE username = $1
	`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.AccountType, &user.IsActive, &user.IsVerified, &user.LastLogin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Get user profile
	profile, err := h.getUserProfile(user.ID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
	}

	// Create response
	response := models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		AccountType: user.AccountType,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		LastLogin:   user.LastLogin,
		Profile:     profile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var profile models.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update profile
	_, err := h.db.Exec(`
		UPDATE user_profiles 
		SET phone = $1, address = $2, city = $3, state = $4, postal_code = $5, country = $6, company_name = $7, updated_at = NOW()
		WHERE user_id = $8
	`, profile.Phone, profile.Address, profile.City, profile.State, profile.PostalCode, profile.Country, profile.CompanyName, userID)

	if err != nil {
		log.Printf("Error updating profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser deactivates a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	_, err := h.db.Exec("UPDATE users SET is_active = false WHERE id = $1", userID)
	if err != nil {
		log.Printf("Error deactivating user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Helper functions

func (h *UserHandler) getUserByID(userID interface{}) (*models.UserResponse, error) {
	var user models.User
	err := h.db.QueryRow(`
		SELECT id, username, email, full_name, account_type, is_active, is_verified, last_login_at, created_at, updated_at
		FROM users 
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.AccountType, &user.IsActive, &user.IsVerified, &user.LastLogin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get user profile
	profile, err := h.getUserProfile(user.ID)
	if err != nil {
		log.Printf("Error getting user profile: %v", err)
	}

	return &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		AccountType: user.AccountType,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		CreatedAt:   user.CreatedAt,
		LastLogin:   user.LastLogin,
		Profile:     profile,
	}, nil
}

func (h *UserHandler) getUserProfile(userID int64) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := h.db.QueryRow(`
		SELECT id, user_id, phone, address, city, state, postal_code, country, company_name, profile_image_url, preferences, created_at, updated_at
		FROM user_profiles 
		WHERE user_id = $1
	`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Phone, &profile.Address, &profile.City,
		&profile.State, &profile.PostalCode, &profile.Country, &profile.CompanyName,
		&profile.ProfileImageUrl, &profile.Preferences, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No profile is ok
		}
		return nil, err
	}

	return &profile, nil
}
