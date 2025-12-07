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

pub async fn get_admins(pool: &Pool) -> Result<Vec<UserResponse>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query(
        "SELECT id, username, email, full_name, account_type, is_active, is_verified 
         FROM users WHERE account_type IN ('admin', 'main_admin') ORDER BY created_at DESC",
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

use crate::models::PasswordChangeRequest;

pub async fn get_password_change_requests(pool: &Pool) -> Result<Vec<PasswordChangeRequest>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query(
        "SELECT r.id, r.user_id, u.username, u.email, r.new_password_hash, r.status, r.created_at 
         FROM password_change_requests r
         JOIN users u ON r.user_id = u.id
         WHERE r.status = 'Pending'
         ORDER BY r.created_at ASC",
        &[]
    ).await?;

    Ok(rows.iter().map(|row| PasswordChangeRequest {
        id: row.get("id"),
        user_id: row.get("user_id"),
        username: row.get("username"),
        email: row.get("email"),
        new_password_hash: row.get("new_password_hash"),
        status: row.get("status"),
        created_at: row.get("created_at"),
    }).collect())
}

pub async fn resolve_password_change_request(pool: &Pool, request_id: i64, approved: bool) -> Result<(), Box<dyn std::error::Error>> {
    let mut client = pool.get().await?;
    let transaction = client.transaction().await?;
    
    if approved {
        // Get the new hash
        let row = transaction.query_opt(
            "SELECT user_id, new_password_hash FROM password_change_requests WHERE id = $1",
            &[&request_id]
        ).await?;
        
        if let Some(r) = row {
             let user_id: i64 = r.get("user_id");
             let new_hash: String = r.get("new_password_hash");
             
             // Update user password
             transaction.execute(
                 "UPDATE users SET password_hash = $1 WHERE id = $2",
                 &[&new_hash, &user_id]
             ).await?;
             
             // Mark request approved
             transaction.execute(
                 "UPDATE password_change_requests SET status = 'Approved', updated_at = NOW() WHERE id = $1",
                 &[&request_id]
             ).await?;
        } else {
            return Err("Request not found".into());
        }
    } else {
        // Mark rejected
        transaction.execute(
            "UPDATE password_change_requests SET status = 'Rejected', updated_at = NOW() WHERE id = $1",
            &[&request_id]
        ).await?;
    }
    
    transaction.commit().await?;
    Ok(())
}
