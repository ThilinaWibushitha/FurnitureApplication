package models

import (
	"time"
)

// User represents a system user
type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	FullName     string     `db:"full_name" json:"full_name"`
	AccountType  string     `db:"account_type" json:"account_type"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	IsVerified   bool       `db:"is_verified" json:"is_verified"`
	LastLogin    *time.Time `db:"last_login_at" json:"last_login_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// UserProfile represents extended user information
type UserProfile struct {
	ID              int64     `db:"id" json:"id"`
	UserID          int64     `db:"user_id" json:"user_id"`
	Phone           string    `db:"phone" json:"phone"`
	Address         string    `db:"address" json:"address"`
	City            string    `db:"city" json:"city"`
	State           string    `db:"state" json:"state"`
	PostalCode      string    `db:"postal_code" json:"postal_code"`
	Country         string    `db:"country" json:"country"`
	CompanyName     string    `db:"company_name" json:"company_name"`
	ProfileImageUrl string    `db:"profile_image_url" json:"profile_image_url"`
	Preferences     string    `db:"preferences" json:"preferences"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUserRequest represents user registration request
type CreateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	AccountType string `json:"account_type"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	CompanyName string `json:"company_name"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
	ID          int64        `json:"id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	FullName    string       `json:"full_name"`
	AccountType string       `json:"account_type"`
	IsActive    bool         `json:"is_active"`
	IsVerified  bool         `json:"is_verified"`
	CreatedAt   time.Time    `json:"created_at"`
	LastLogin   *time.Time   `json:"last_login_at"`
	Profile     *UserProfile `json:"profile,omitempty"`
}
