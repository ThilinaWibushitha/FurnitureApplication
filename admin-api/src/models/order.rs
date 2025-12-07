use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Order {
    pub id: i64,
    pub user_id: i64,
    pub total_amount: f64,
    pub status: String, // "pending", "confirmed", "shipped", "delivered", "cancelled"
    pub payment_status: String, // "pending", "paid", "refunded"
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OrderItem {
    pub id: i64,
    pub order_id: i64,
    pub item_id: i64,
    pub quantity: i32,
    pub unit_price: f64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateOrderRequest {
    pub status: Option<String>,
    pub payment_status: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SalesReport {
    pub period: String,
    pub total_sales: f64,
    pub total_orders: i64,
    pub items_sold: i64,
}
