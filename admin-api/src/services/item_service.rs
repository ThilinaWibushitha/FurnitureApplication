use deadpool_postgres::Pool;
use crate::models::{Item, CreateItemRequest, UpdateItemRequest};

pub async fn get_all_items(pool: &Pool) -> Result<Vec<Item>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let rows = client.query("SELECT id, name, description, price::FLOAT8, stock_quantity, category, image_url, is_active, created_at, updated_at FROM items ORDER BY created_at DESC", &[]).await?;
    
    Ok(rows.iter().map(|row| Item {
        id: row.get("id"),
        name: row.get("name"),
        description: row.get("description"),
        price: row.get("price"),
        stock_quantity: row.get("stock_quantity"),
        category: row.get("category"),
        image_url: row.get("image_url"),
        is_active: row.get("is_active"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    }).collect())
}

pub async fn get_item_by_id(pool: &Pool, id: i64) -> Result<Option<Item>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    let row = client.query_opt("SELECT id, name, description, price::FLOAT8, stock_quantity, category, image_url, is_active, created_at, updated_at FROM items WHERE id = $1", &[&id]).await?;
    
    Ok(row.map(|r| Item {
        id: r.get("id"),
        name: r.get("name"),
        description: r.get("description"),
        price: r.get("price"),
        stock_quantity: r.get("stock_quantity"),
        category: r.get("category"),
        image_url: r.get("image_url"),
        is_active: r.get("is_active"),
        created_at: r.get("created_at"),
        updated_at: r.get("updated_at"),
    }))
}

pub async fn create_item(pool: &Pool, req: &CreateItemRequest) -> Result<Item, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    log::info!("Creating item: name={}, desc={:?}", req.name, req.description);
    let row = client.query_one(
        "INSERT INTO items (name, description, price, stock_quantity, category, image_url) 
         VALUES ($1, $2::TEXT, $3::FLOAT8, $4, $5, $6::TEXT) 
         RETURNING id, name, description, price::FLOAT8, stock_quantity, category, image_url, is_active, created_at, updated_at",
        &[&req.name, &req.description.clone().unwrap_or_default(), &req.price, &req.stock_quantity, &req.category, &req.image_url.clone().unwrap_or_default()]
    ).await?;
    
    Ok(Item {
        id: row.get("id"),
        name: row.get("name"),
        description: row.get("description"),
        price: row.get("price"),
        stock_quantity: row.get("stock_quantity"),
        category: row.get("category"),
        image_url: row.get("image_url"),
        is_active: row.get("is_active"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}

pub async fn update_item(pool: &Pool, id: i64, req: &UpdateItemRequest) -> Result<Item, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_one(
        "UPDATE items SET 
            name = COALESCE($1, name),
            description = COALESCE($2::TEXT, description),
            price = COALESCE($3, price),
            stock_quantity = COALESCE($4, stock_quantity),
            category = COALESCE($5, category),
            image_url = COALESCE($6::TEXT, image_url),
            is_active = COALESCE($7, is_active),
            is_active = COALESCE($7, is_active),
            updated_at = NOW()
         WHERE id = $8 
         RETURNING id, name, description, price::FLOAT8, stock_quantity, category, image_url, is_active, created_at, updated_at",
        &[&req.name, &req.description.clone().unwrap_or_default(), &req.price, &req.stock_quantity, &req.category, &req.image_url.clone().unwrap_or_default(), &req.is_active, &id]
    ).await?;
    
    Ok(Item {
        id: row.get("id"),
        name: row.get("name"),
        description: row.get("description"),
        price: row.get("price"),
        stock_quantity: row.get("stock_quantity"),
        category: row.get("category"),
        image_url: row.get("image_url"),
        is_active: row.get("is_active"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}

pub async fn delete_item(pool: &Pool, id: i64) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    client.execute("DELETE FROM items WHERE id = $1", &[&id]).await?;
    Ok(())
}
