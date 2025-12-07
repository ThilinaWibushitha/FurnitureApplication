use deadpool_postgres::Pool;
use crate::models::payment::{Payment, CreatePaymentRequest, UpdatePaymentRequest, OrderCancellation, ProcessRefundRequest};
use tokio_postgres::types::ToSql;
use crate::handlers::payments::PaymentFilter;

pub async fn list_payments(pool: &Pool, filter: &PaymentFilter) -> Result<Vec<Payment>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let mut query = String::from("SELECT id, order_id, customer_id, amount::FLOAT8, payment_method, payment_status, transaction_ref, gateway_response, processed_by, processed_at, notes, created_at, updated_at FROM payments WHERE 1=1");
    let mut param_values: Vec<Box<dyn ToSql + Sync>> = Vec::new();

    if let Some(ref status) = filter.status {
        query.push_str(&format!(" AND payment_status = ${}", param_values.len() + 1));
        param_values.push(Box::new(status.clone()));
    }

    if let Some(customer_id) = filter.customer_id {
        query.push_str(&format!(" AND customer_id = ${}", param_values.len() + 1));
        param_values.push(Box::new(customer_id));
    }

    if let Some(order_id) = filter.order_id {
        query.push_str(&format!(" AND order_id = ${}", param_values.len() + 1));
        param_values.push(Box::new(order_id));
    }

    query.push_str(" ORDER BY created_at DESC LIMIT 100");

    let params: Vec<&(dyn ToSql + Sync)> = param_values.iter().map(|p| p.as_ref()).collect();
    let rows = client.query(&query, &params[..]).await?;
    
    let payments: Vec<Payment> = rows.iter().map(|row| Payment {
        id: row.get("id"),
        order_id: row.get("order_id"),
        customer_id: row.get("customer_id"),
        amount: row.get("amount"),
        payment_method: row.get("payment_method"),
        payment_status: row.get("payment_status"),
        transaction_ref: row.get("transaction_ref"),
        gateway_response: row.get("gateway_response"),
        processed_by: row.get("processed_by"),
        processed_at: row.get("processed_at"),
        notes: row.get("notes"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    }).collect();
    
    Ok(payments)
}

pub async fn get_payment(pool: &Pool, payment_id: i64) -> Result<Option<Payment>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_opt(
        "SELECT id, order_id, customer_id, amount::FLOAT8, payment_method, payment_status, transaction_ref, gateway_response, processed_by, processed_at, notes, created_at, updated_at FROM payments WHERE id = $1",
        &[&payment_id]
    ).await?;
    
    Ok(row.map(|r| Payment {
        id: r.get("id"),
        order_id: r.get("order_id"),
        customer_id: r.get("customer_id"),
        amount: r.get("amount"),
        payment_method: r.get("payment_method"),
        payment_status: r.get("payment_status"),
        transaction_ref: r.get("transaction_ref"),
        gateway_response: r.get("gateway_response"),
        processed_by: r.get("processed_by"),
        processed_at: r.get("processed_at"),
        notes: r.get("notes"),
        created_at: r.get("created_at"),
        updated_at: r.get("updated_at"),
    }))
}

pub async fn get_order_payments(pool: &Pool, order_id: i64) -> Result<Vec<Payment>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let rows = client.query(
        "SELECT id, order_id, customer_id, amount::FLOAT8, payment_method, payment_status, transaction_ref, gateway_response, processed_by, processed_at, notes, created_at, updated_at FROM payments WHERE order_id = $1 ORDER BY created_at DESC",
        &[&order_id]
    ).await?;
    
    let payments: Vec<Payment> = rows.iter().map(|row| Payment {
        id: row.get("id"),
        order_id: row.get("order_id"),
        customer_id: row.get("customer_id"),
        amount: row.get("amount"),
        payment_method: row.get("payment_method"),
        payment_status: row.get("payment_status"),
        transaction_ref: row.get("transaction_ref"),
        gateway_response: row.get("gateway_response"),
        processed_by: row.get("processed_by"),
        processed_at: row.get("processed_at"),
        notes: row.get("notes"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    }).collect();
    
    Ok(payments)
}

pub async fn create_payment(pool: &Pool, req: &CreatePaymentRequest) -> Result<Payment, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    // Get customer ID from order
    let customer_row = client.query_opt(
        "SELECT customer_id FROM orders WHERE id = $1",
        &[&req.order_id]
    ).await?;
    
    let customer_id: Option<i64> = customer_row.and_then(|r| r.get("customer_id"));
    
    let row = client.query_one(
        r#"
        INSERT INTO payments (order_id, customer_id, amount, payment_method, payment_status, notes, created_at, updated_at)
        VALUES ($1, $2, $3, $4, 'Pending', $5, NOW(), NOW())
        RETURNING id, order_id, customer_id, amount::FLOAT8, payment_method, payment_status, transaction_ref, gateway_response, processed_by, processed_at, notes, created_at, updated_at
        "#,
        &[&req.order_id, &customer_id, &req.amount, &req.payment_method, &req.notes]
    ).await?;
    
    Ok(Payment {
        id: row.get("id"),
        order_id: row.get("order_id"),
        customer_id: row.get("customer_id"),
        amount: row.get("amount"),
        payment_method: row.get("payment_method"),
        payment_status: row.get("payment_status"),
        transaction_ref: row.get("transaction_ref"),
        gateway_response: row.get("gateway_response"),
        processed_by: row.get("processed_by"),
        processed_at: row.get("processed_at"),
        notes: row.get("notes"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}

pub async fn update_payment(pool: &Pool, payment_id: i64, req: &UpdatePaymentRequest) -> Result<Payment, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let mut updates = vec!["updated_at = NOW()".to_string()];
    let mut param_values: Vec<Box<dyn ToSql + Sync>> = Vec::new();

    if let Some(ref status) = req.payment_status {
        let idx = param_values.len() + 1;
        updates.push(format!("payment_status = ${}", idx));
        param_values.push(Box::new(status.clone()));

        if status == "Completed" {
            updates.push("processed_at = NOW()".to_string());
        }
    }

    if let Some(ref transaction_ref) = req.transaction_ref {
        let idx = param_values.len() + 1;
        updates.push(format!("transaction_ref = ${}", idx));
        param_values.push(Box::new(transaction_ref.clone()));
    }

    if let Some(ref notes) = req.notes {
        let idx = param_values.len() + 1;
        updates.push(format!("notes = ${}", idx));
        param_values.push(Box::new(notes.clone()));
    }

    param_values.push(Box::new(payment_id));
    let id_idx = param_values.len();

    let query = format!(
        "UPDATE payments SET {} WHERE id = ${} 
         RETURNING id, order_id, customer_id, amount::FLOAT8, payment_method, payment_status, transaction_ref, gateway_response, processed_by, processed_at, notes, created_at, updated_at",
        updates.join(", "),
        id_idx
    );

    let params: Vec<&(dyn ToSql + Sync)> = param_values
        .iter()
        .map(|p| p.as_ref() as &(dyn ToSql + Sync))
        .collect();
    
    let row = client.query_one(&query, &params[..]).await?;
    
    Ok(Payment {
        id: row.get("id"),
        order_id: row.get("order_id"),
        customer_id: row.get("customer_id"),
        amount: row.get("amount"),
        payment_method: row.get("payment_method"),
        payment_status: row.get("payment_status"),
        transaction_ref: row.get("transaction_ref"),
        gateway_response: row.get("gateway_response"),
        processed_by: row.get("processed_by"),
        processed_at: row.get("processed_at"),
        notes: row.get("notes"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    })
}

pub async fn list_cancellations(pool: &Pool) -> Result<Vec<OrderCancellation>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let rows = client.query(
        "SELECT id, order_id, payment_id, requested_by, reason, cancellation_type, refund_amount::FLOAT8, refund_status, processed_by, processed_at, notes, created_at, updated_at FROM order_cancellations WHERE refund_status = 'Pending' ORDER BY created_at DESC",
        &[]
    ).await?;
    
    let cancellations: Vec<OrderCancellation> = rows.iter().map(|row| OrderCancellation {
        id: row.get("id"),
        order_id: row.get("order_id"),
        payment_id: row.get("payment_id"),
        requested_by: row.get("requested_by"),
        reason: row.get("reason"),
        cancellation_type: row.get("cancellation_type"),
        refund_amount: row.get("refund_amount"),
        refund_status: row.get("refund_status"),
        processed_by: row.get("processed_by"),
        processed_at: row.get("processed_at"),
        notes: row.get("notes"),
        created_at: row.get("created_at"),
        updated_at: row.get("updated_at"),
    }).collect();
    
    Ok(cancellations)
}

pub async fn process_refund(pool: &Pool, cancellation_id: i64, req: &ProcessRefundRequest) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    client.execute(
        r#"
        UPDATE order_cancellations 
        SET refund_status = $1, notes = COALESCE($2, notes), processed_at = NOW(), updated_at = NOW()
        WHERE id = $3
        "#,
        &[&req.status, &req.notes, &cancellation_id]
    ).await?;
    
    // If completed, update the payment status to Refunded
    if req.status == "Completed" {
        let cancellation = client.query_opt(
            "SELECT payment_id FROM order_cancellations WHERE id = $1",
            &[&cancellation_id]
        ).await?;
        
        if let Some(row) = cancellation {
            if let Some(payment_id) = row.get::<_, Option<i64>>("payment_id") {
                client.execute(
                    "UPDATE payments SET payment_status = 'Refunded', updated_at = NOW() WHERE id = $1",
                    &[&payment_id]
                ).await?;
            }
        }
    }
    
    Ok(())
}
