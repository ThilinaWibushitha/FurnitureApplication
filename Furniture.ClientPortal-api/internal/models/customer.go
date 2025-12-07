package models

import "time"

// Customer represents a customer
type Customer struct {
	ID            int64     `db:"id" json:"id"`
	FirstName     string    `db:"first_name" json:"first_name"`
	LastName      string    `db:"last_name" json:"last_name"`
	Email         string    `db:"email" json:"email"`
	Phone         string    `db:"phone" json:"phone"`
	Address       string    `db:"address" json:"address"`
	City          string    `db:"city" json:"city"`
	LoyaltyPoints int       `db:"loyalty_points" json:"loyalty_points"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// Supplier represents a supplier
type Supplier struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	ContactName string    `db:"contact_name" json:"contact_name"`
	Email       string    `db:"email" json:"email"`
	Phone       string    `db:"phone" json:"phone"`
	Address     string    `db:"address" json:"address"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
