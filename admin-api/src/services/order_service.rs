use deadpool_postgres::Pool;
use crate::models::{Order, UpdateOrderRequest};

pub async fn get_all_orders(pool: &Pool) -> Result<Vec<Order>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query("SELECT id, user_id, total_amount::FLOAT8, status, payment_status, created_at, updated_at FROM orders ORDER BY created_at DESC", &[]).await?;
    
    Ok(rows.iter().map(|row| Order {
        id: row.get("id"),
        user_id: row.get("user_id"),
        total_amount: row.get("total_amount"),
        status: row.get("status"),
        payment_status: row.get("payment_status"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    }).collect())
}

pub async fn get_order_by_id(pool: &Pool, id: i64) -> Result<Option<Order>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let row = client.query_opt("SELECT id, user_id, total_amount::FLOAT8, status, payment_status, created_at, updated_at FROM orders WHERE id = $1", &[&id]).await?;
    
    Ok(row.map(|r| Order {
        id: r.get("id"),
        user_id: r.get("user_id"),
        total_amount: r.get("total_amount"),
        status: r.get("status"),
        payment_status: r.get("payment_status"),
        created_at: r.get("created_at"),
        updated_at: r.get("updated_at"),
    }))
}

pub async fn update_order(pool: &Pool, id: i64, req: &UpdateOrderRequest) -> Result<Order, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_one(
        "UPDATE orders SET 
            status = COALESCE($1, status),
            payment_status = COALESCE($2, payment_status),
            updated_at = NOW()
         WHERE id = $3 
         RETURNING id, user_id, total_amount::FLOAT8, status, payment_status, created_at, updated_at",
        &[&req.status, &req.payment_status, &id]
    ).await?;
    
    Ok(Order {
        id: row.get("id"),
        user_id: row.get("user_id"),
        total_amount: row.get("total_amount"),
        status: row.get("status"),
        payment_status: row.get("payment_status"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}
