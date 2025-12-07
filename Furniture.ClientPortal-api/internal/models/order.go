package models

import "time"

// Order represents a sales order
type Order struct {
	ID             int64     `db:"id" json:"id"`
	OrderNumber    string    `db:"order_number" json:"order_number"`
	CustomerID     *int64    `db:"customer_id" json:"customer_id"`
	UserID         *int64    `db:"user_id" json:"user_id"`
	BranchID       int64     `db:"branch_id" json:"branch_id"`
	Status         string    `db:"status" json:"status"`
	TotalAmount    float64   `db:"total_amount" json:"total_amount"`
	TaxAmount      float64   `db:"tax_amount" json:"tax_amount"`
	DiscountAmount float64   `db:"discount_amount" json:"discount_amount"`
	PaymentMethod  string    `db:"payment_method" json:"payment_method"`
	Notes          string    `db:"notes" json:"notes"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID            int64     `db:"id" json:"id"`
	OrderID       int64     `db:"order_id" json:"order_id"`
	ItemVariantID *int64    `db:"item_variant_id" json:"item_variant_id"`
	Quantity      int       `db:"quantity" json:"quantity"`
	UnitPrice     float64   `db:"unit_price" json:"unit_price"`
	Subtotal      float64   `db:"subtotal" json:"subtotal"`
	Discount      float64   `db:"discount" json:"discount"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
