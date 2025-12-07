package models

import "time"

// Payment represents a payment transaction
type Payment struct {
	ID              int64      `db:"id" json:"id"`
	OrderID         int64      `db:"order_id" json:"order_id"`
	CustomerID      *int64     `db:"customer_id" json:"customer_id,omitempty"`
	Amount          float64    `db:"amount" json:"amount"`
	PaymentMethod   string     `db:"payment_method" json:"payment_method"`
	PaymentStatus   string     `db:"payment_status" json:"payment_status"`
	TransactionRef  string     `db:"transaction_ref" json:"transaction_ref,omitempty"`
	GatewayResponse string     `db:"gateway_response" json:"gateway_response,omitempty"`
	ProcessedBy     *int64     `db:"processed_by" json:"processed_by,omitempty"`
	ProcessedAt     *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	Notes           string     `db:"notes" json:"notes,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

// CreatePaymentRequest is the request body for creating a payment
type CreatePaymentRequest struct {
	OrderID       int64   `json:"order_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Notes         string  `json:"notes,omitempty"`
}

// UpdatePaymentRequest is the request body for updating a payment
type UpdatePaymentRequest struct {
	PaymentStatus  string `json:"payment_status,omitempty"`
	TransactionRef string `json:"transaction_ref,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

// ItemRating represents a customer rating for an item
type ItemRating struct {
	ID                 int64      `db:"id" json:"id"`
	ItemID             int64      `db:"item_id" json:"item_id"`
	CustomerID         int64      `db:"customer_id" json:"customer_id"`
	OrderID            *int64     `db:"order_id" json:"order_id,omitempty"`
	Rating             int        `db:"rating" json:"rating"`
	ReviewText         string     `db:"review_text" json:"review_text,omitempty"`
	IsVerifiedPurchase bool       `db:"is_verified_purchase" json:"is_verified_purchase"`
	IsApproved         bool       `db:"is_approved" json:"is_approved"`
	ApprovedBy         *int64     `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt         *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
}

// CreateRatingRequest is the request body for creating a rating
type CreateRatingRequest struct {
	ItemID     int64  `json:"item_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	ReviewText string `json:"review_text,omitempty"`
	OrderID    *int64 `json:"order_id,omitempty"`
}

// ClientProfile represents extended client information
type ClientProfile struct {
	ID                  int64      `db:"id" json:"id"`
	CustomerID          int64      `db:"customer_id" json:"customer_id"`
	UserID              *int64     `db:"user_id" json:"user_id,omitempty"`
	DateOfBirth         *string    `db:"date_of_birth" json:"date_of_birth,omitempty"`
	Gender              string     `db:"gender" json:"gender,omitempty"`
	Preferences         string     `db:"preferences" json:"preferences,omitempty"`
	NewsletterOptIn     bool       `db:"newsletter_opt_in" json:"newsletter_opt_in"`
	ProfileImageUrl     string     `db:"profile_image_url" json:"profile_image_url,omitempty"`
	AccountStatus       string     `db:"account_status" json:"account_status"`
	DeletionRequestedAt *time.Time `db:"deletion_requested_at" json:"deletion_requested_at,omitempty"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
}

// UpdateProfileRequest is the request body for updating client profile
type UpdateProfileRequest struct {
	Phone           string `json:"phone,omitempty"`
	Address         string `json:"address,omitempty"`
	City            string `json:"city,omitempty"`
	DateOfBirth     string `json:"date_of_birth,omitempty"`
	Gender          string `json:"gender,omitempty"`
	NewsletterOptIn *bool  `json:"newsletter_opt_in,omitempty"`
	ProfileImageUrl string `json:"profile_image_url,omitempty"`
}

// OrderCancellation represents a cancellation request
type OrderCancellation struct {
	ID               int64      `db:"id" json:"id"`
	OrderID          int64      `db:"order_id" json:"order_id"`
	PaymentID        *int64     `db:"payment_id" json:"payment_id,omitempty"`
	RequestedBy      *int64     `db:"requested_by" json:"requested_by,omitempty"`
	Reason           string     `db:"reason" json:"reason"`
	CancellationType string     `db:"cancellation_type" json:"cancellation_type"`
	RefundAmount     *float64   `db:"refund_amount" json:"refund_amount,omitempty"`
	RefundStatus     string     `db:"refund_status" json:"refund_status"`
	ProcessedBy      *int64     `db:"processed_by" json:"processed_by,omitempty"`
	ProcessedAt      *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	Notes            string     `db:"notes" json:"notes,omitempty"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
}

// CancelOrderRequest is the request body for cancelling an order
type CancelOrderRequest struct {
	Reason           string   `json:"reason" binding:"required"`
	CancellationType string   `json:"cancellation_type" binding:"required"`
	RefundAmount     *float64 `json:"refund_amount,omitempty"`
}

// DeleteProfileRequest is the request body for account deletion
type DeleteProfileRequest struct {
	Reason   string `json:"reason,omitempty"`
	Password string `json:"password" binding:"required"`
}
