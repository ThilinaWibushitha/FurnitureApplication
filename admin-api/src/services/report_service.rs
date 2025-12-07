use deadpool_postgres::Pool;
use crate::models::SalesReport;

pub async fn get_daily_sales(pool: &Pool) -> Result<SalesReport, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_one(
        "SELECT 
            COALESCE(SUM(total_amount), 0)::FLOAT8 as total_sales,
            COUNT(*) as total_orders,
            COALESCE(SUM(oi.quantity), 0) as items_sold
         FROM orders o
         LEFT JOIN order_items oi ON o.id = oi.order_id
         WHERE DATE(o.created_at) = CURRENT_DATE",
        &[]
    ).await?;
    
    Ok(SalesReport {
        period: "daily".to_string(),
        total_sales: row.get::<_, f64>("total_sales"),
        total_orders: row.get::<_, i64>("total_orders"),
        items_sold: row.get::<_, i64>("items_sold"),
    })
}

pub async fn get_monthly_sales(pool: &Pool) -> Result<SalesReport, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_one(
        "SELECT 
            COALESCE(SUM(total_amount), 0)::FLOAT8 as total_sales,
            COUNT(*) as total_orders,
            COALESCE(SUM(oi.quantity), 0) as items_sold
         FROM orders o
         LEFT JOIN order_items oi ON o.id = oi.order_id
         WHERE DATE_TRUNC('month', o.created_at) = DATE_TRUNC('month', CURRENT_DATE)",
        &[]
    ).await?;
    
    Ok(SalesReport {
        period: "monthly".to_string(),
        total_sales: row.get::<_, f64>("total_sales"),
        total_orders: row.get::<_, i64>("total_orders"),
        items_sold: row.get::<_, i64>("items_sold"),
    })
}
