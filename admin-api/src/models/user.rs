use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct User {
    pub id: i64,
    pub username: String,
    pub email: String,
    pub password_hash: String,
    pub full_name: String,
    pub account_type: String, // "admin" or "client"
    pub is_active: bool,
    pub is_verified: bool,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
    pub last_login_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct CreateUserRequest {
    pub username: String,
    pub email: String,
    pub password: String,
    pub full_name: String,
    pub account_type: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoginRequest {
    pub email: String,
    pub password: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoginResponse {
    pub token: String,
    pub user: UserResponse,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UserResponse {
    pub id: i64,
    pub username: String,
    pub email: String,
    pub full_name: String,
    pub account_type: String,
    pub is_active: bool,
    pub is_verified: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ApproveUserRequest {
    pub user_id: i64,
    pub approved: bool,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct PasswordChangeRequest {
    pub id: i64,
    pub user_id: i64,
    pub username: String,
    pub email: String,
    pub new_password_hash: String,
    pub status: String,
    pub created_at: DateTime<Utc>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ResolvePasswordChangeRequest {
    pub request_id: i64,
    pub approved: bool,
}
