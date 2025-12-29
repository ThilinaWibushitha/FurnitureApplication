package models

import "time"

// Department represents a furniture department
// Corresponds to departments table

type Department struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Color     string    `db:"color" json:"color"`
	Icon      string    `db:"icon" json:"icon"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Category represents a furniture category
// Corresponds to categories table

type Category struct {
	ID                int64     `db:"id" json:"id"`
	DepartmentID      int64     `db:"department_id" json:"department_id"`
	Name              string    `db:"name" json:"name"`
	ParentCategoryID  *int64    `db:"parent_category_id" json:"parent_category_id,omitempty"`
	SortOrder         int       `db:"sort_order" json:"sort_order"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	IsPublishedOnline bool      `db:"is_published_online" json:"is_published_online"`
	OnlineSortOrder   int       `db:"online_sort_order" json:"online_sort_order"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// Item represents a furniture item
// Corresponds to items table

type Item struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Description   *string   `db:"description" json:"description,omitempty"`
	Price         float64   `db:"price" json:"price"`
	StockQuantity int       `db:"stock_quantity" json:"stock_quantity"`
	Category      *string   `db:"category" json:"category,omitempty"`
	ImageURL      *string   `db:"image_url" json:"image_url,omitempty"`
	SKU           *string   `db:"sku" json:"sku,omitempty"`
	Dimensions    *string   `db:"dimensions" json:"dimensions,omitempty"`
	Material      *string   `db:"material" json:"material,omitempty"`
	Color         *string   `db:"color" json:"color,omitempty"`
	IsActive      *bool     `db:"is_active" json:"is_active,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// ItemVariant represents a variant of a furniture item
// Corresponds to item_variants table

type ItemVariant struct {
	ID              int64     `db:"id" json:"id"`
	ItemID          int64     `db:"item_id" json:"item_id"`
	VariantCode     string    `db:"variant_code" json:"variant_code"`
	Barcode         string    `db:"barcode" json:"barcode"`
	QRCode          string    `db:"qr_code" json:"qr_code"`
	Color           string    `db:"color" json:"color"`
	Fabric          string    `db:"fabric" json:"fabric"`
	Size            string    `db:"size" json:"size"`
	Material        string    `db:"material" json:"material"`
	Price           *float64  `db:"price" json:"price,omitempty"`
	Cost            *float64  `db:"cost" json:"cost,omitempty"`
	DiscountPercent *float64  `db:"discount_percent" json:"discount_percent,omitempty"`
	TaxClassID      *int64    `db:"tax_class_id" json:"tax_class_id,omitempty"`
	Status          string    `db:"status" json:"status"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Branch represents a store branch
// Corresponds to branches table

type Branch struct {
	ID       int64  `db:"id" json:"id"`
	Code     string `db:"code" json:"code"`
	Name     string `db:"name" json:"name"`
	Address  string `db:"address" json:"address"`
	Phone    string `db:"phone" json:"phone"`
	IsActive bool   `db:"is_active" json:"is_active"`
}

// InventoryLevel represents inventory levels for item variants at branches
// Corresponds to inventory_levels table

type InventoryLevel struct {
	ID               int64     `db:"id" json:"id"`
	BranchID         int64     `db:"branch_id" json:"branch_id"`
	ItemVariantID    int64     `db:"item_variant_id" json:"item_variant_id"`
	QuantityOnHand   float64   `db:"quantity_on_hand" json:"quantity_on_hand"`
	QuantityReserved float64   `db:"quantity_reserved" json:"quantity_reserved"`
	ReorderLevel     float64   `db:"reorder_level" json:"reorder_level"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
