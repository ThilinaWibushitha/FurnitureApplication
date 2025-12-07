use deadpool_postgres::Pool;
use crate::models::{LoginRequest, LoginResponse, UserResponse, CreateUserRequest};
use bcrypt::{hash, verify, DEFAULT_COST};
use jsonwebtoken::{encode, Header, EncodingKey};
use serde::{Deserialize, Serialize};
use std::env;

#[derive(Debug, Serialize, Deserialize)]
struct Claims {
    sub: i64,
    email: String,
    account_type: String,
    exp: usize,
}

pub async fn authenticate(pool: &Pool, request: &LoginRequest) -> Result<LoginResponse, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_opt(
        "SELECT id, username, email, password_hash, full_name, account_type, is_active, is_verified 
         FROM users WHERE email = $1 AND account_type = 'admin'",
        &[&request.email]
    ).await?;

    let row = row.ok_or("Invalid credentials")?;
    
    let password_hash: String = row.get("password_hash");
    if !verify(&request.password, &password_hash)? {
        return Err("Invalid credentials".into());
    }

    let user = UserResponse {
        id: row.get("id"),
        username: row.get("username"),
        email: row.get("email"),
        full_name: row.get("full_name"),
        account_type: row.get("account_type"),
        is_active: row.get("is_active"),
        is_verified: row.get("is_verified"),
    };

    let secret = env::var("JWT_SECRET").unwrap_or_else(|_| "secret".to_string());
    let claims = Claims {
        sub: user.id,
        email: user.email.clone(),
        account_type: user.account_type.clone(),
        exp: (chrono::Utc::now() + chrono::Duration::hours(24)).timestamp() as usize,
    };

    let token = encode(&Header::default(), &claims, &EncodingKey::from_secret(secret.as_bytes()))?;

    Ok(LoginResponse { token, user })
}

pub async fn create_user(pool: &Pool, request: &CreateUserRequest) -> Result<UserResponse, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let password_hash = hash(&request.password, DEFAULT_COST)?;

    let row = client.query_one(
        "INSERT INTO users (username, email, password_hash, full_name, account_type, is_verified) 
         VALUES ($1, $2, $3, $4, $5, true) RETURNING id, username, email, full_name, account_type, is_active, is_verified",
        &[&request.username, &request.email, &password_hash, &request.full_name, &request.account_type]
    ).await?;

    Ok(UserResponse {
        id: row.get("id"),
        username: row.get("username"),
        email: row.get("email"),
        full_name: row.get("full_name"),
        account_type: row.get("account_type"),
        is_active: row.get("is_active"),
        is_verified: row.get("is_verified"),
    })
}
