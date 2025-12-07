use deadpool_postgres::Pool;
use crate::models::UserResponse;

pub async fn get_all_users(pool: &Pool) -> Result<Vec<UserResponse>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query(
        "SELECT id, username, email, full_name, account_type, is_active, is_verified FROM users ORDER BY created_at DESC",
        &[]
    ).await?;

    Ok(rows.iter().map(|row| UserResponse {
        id: row.get("id"),
        username: row.get("username"),
        email: row.get("email"),
        full_name: row.get("full_name"),
        account_type: row.get("account_type"),
        is_active: row.get("is_active"),
        is_verified: row.get("is_verified"),
    }).collect())
}

pub async fn get_pending_users(pool: &Pool) -> Result<Vec<UserResponse>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query(
        "SELECT id, username, email, full_name, account_type, is_active, is_verified 
         FROM users WHERE is_verified = false AND account_type = 'client' ORDER BY created_at DESC",
        &[]
    ).await?;

    Ok(rows.iter().map(|row| UserResponse {
        id: row.get("id"),
        username: row.get("username"),
        email: row.get("email"),
        full_name: row.get("full_name"),
        account_type: row.get("account_type"),
        is_active: row.get("is_active"),
        is_verified: row.get("is_verified"),
    }).collect())
}

pub async fn get_user_by_id(pool: &Pool, user_id: i64) -> Result<Option<UserResponse>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let row = client.query_opt(
        "SELECT id, username, email, full_name, account_type, is_active, is_verified FROM users WHERE id = $1",
        &[&user_id]
    ).await?;

    Ok(row.map(|r| UserResponse {
        id: r.get("id"),
        username: r.get("username"),
        email: r.get("email"),
        full_name: r.get("full_name"),
        account_type: r.get("account_type"),
        is_active: r.get("is_active"),
        is_verified: r.get("is_verified"),
    }))
}

pub async fn approve_user(pool: &Pool, user_id: i64, approved: bool) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    client.execute(
        "UPDATE users SET is_verified = $1, updated_at = NOW() WHERE id = $2",
        &[&approved, &user_id]
    ).await?;
    Ok(())
}

pub async fn update_profile(pool: &Pool, user_id: i64, data: &serde_json::Value) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    if let Some(full_name) = data.get("full_name").and_then(|v| v.as_str()) {
        client.execute("UPDATE users SET full_name = $1, updated_at = NOW() WHERE id = $2", &[&full_name, &user_id]).await?;
    }
    Ok(())
}

pub async fn deactivate_user(pool: &Pool, user_id: i64) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    client.execute("UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1", &[&user_id]).await?;
    Ok(())
}
