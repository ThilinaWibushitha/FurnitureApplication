use chrono::{DateTime, Utc, NaiveDate};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Payment {
    pub id: i64,
    pub order_id: i64,
    pub customer_id: Option<i64>,
    pub amount: f64,
    pub payment_method: String,
    pub payment_status: String,
    pub transaction_ref: Option<String>,
    pub gateway_response: Option<serde_json::Value>,
    pub processed_by: Option<i64>,
    pub processed_at: Option<DateTime<Utc>>,
    pub notes: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreatePaymentRequest {
    pub order_id: i64,
    pub amount: f64,
    pub payment_method: String,
    pub notes: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdatePaymentRequest {
    pub payment_status: Option<String>,
    pub transaction_ref: Option<String>,
    pub notes: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ClientProfile {
    pub id: i64,
    pub customer_id: i64,
    pub user_id: Option<i64>,
    pub date_of_birth: Option<NaiveDate>,
    pub gender: Option<String>,
    pub preferences: Option<serde_json::Value>,
    pub newsletter_opt_in: bool,
    pub profile_image_url: Option<String>,
    pub account_status: String,
    pub deletion_requested_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateClientProfileRequest {
    pub phone: Option<String>,
    pub address: Option<String>,
    pub city: Option<String>,
    pub date_of_birth: Option<String>,
    pub gender: Option<String>,
    pub newsletter_opt_in: Option<bool>,
    pub profile_image_url: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ClientWithProfile {
    pub customer_id: i64,
    pub first_name: String,
    pub last_name: Option<String>,
    pub email: Option<String>,
    pub phone: Option<String>,
    pub address: Option<String>,
    pub city: Option<String>,
    pub loyalty_points: i32,
    pub date_of_birth: Option<NaiveDate>,
    pub gender: Option<String>,
    pub newsletter_opt_in: bool,
    pub profile_image_url: Option<String>,
    pub account_status: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct OrderCancellation {
    pub id: i64,
    pub order_id: i64,
    pub payment_id: Option<i64>,
    pub requested_by: Option<i64>,
    pub reason: String,
    pub cancellation_type: String,
    pub refund_amount: Option<f64>,
    pub refund_status: String,
    pub processed_by: Option<i64>,
    pub processed_at: Option<DateTime<Utc>>,
    pub notes: Option<String>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ProcessRefundRequest {
    pub cancellation_id: i64,
    pub status: String, // "Completed" or "Rejected"
    pub notes: Option<String>,
}
