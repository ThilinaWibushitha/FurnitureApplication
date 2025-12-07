use deadpool_postgres::Pool;
use crate::models::payment::{ClientWithProfile, UpdateClientProfileRequest};
use crate::handlers::clients::ClientFilter;

pub async fn list_clients(pool: &Pool, filter: &ClientFilter) -> Result<Vec<ClientWithProfile>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let mut query = String::from(r#"
        SELECT 
            c.id as customer_id,
            c.first_name,
            c.last_name,
            c.email,
            c.phone,
            c.address,
            c.city,
            c.loyalty_points,
            cp.date_of_birth,
            cp.gender,
            COALESCE(cp.newsletter_opt_in, false) as newsletter_opt_in,
            cp.profile_image_url,
            COALESCE(cp.account_status, 'Active') as account_status
        FROM customers c
        LEFT JOIN client_profiles cp ON c.id = cp.customer_id
        WHERE 1=1
    "#);
    
    if let Some(ref status) = filter.status {
        query.push_str(&format!(" AND COALESCE(cp.account_status, 'Active') = '{}'", status));
    }
    
    if let Some(ref search) = filter.search {
        query.push_str(&format!(
            " AND (c.first_name ILIKE '%{}%' OR c.last_name ILIKE '%{}%' OR c.email ILIKE '%{}%')",
            search, search, search
        ));
    }
    
    query.push_str(" ORDER BY c.created_at DESC LIMIT 100");
    
    let rows = client.query(&query, &[]).await?;
    
    let clients: Vec<ClientWithProfile> = rows.iter().map(|row| ClientWithProfile {
        customer_id: row.get("customer_id"),
        first_name: row.get("first_name"),
        last_name: row.get("last_name"),
        email: row.get("email"),
        phone: row.get("phone"),
        address: row.get("address"),
        city: row.get("city"),
        loyalty_points: row.get("loyalty_points"),
        date_of_birth: row.get("date_of_birth"),
        gender: row.get("gender"),
        newsletter_opt_in: row.get("newsletter_opt_in"),
        profile_image_url: row.get("profile_image_url"),
        account_status: row.get("account_status"),
    }).collect();
    
    Ok(clients)
}

pub async fn get_client(pool: &Pool, customer_id: i64) -> Result<Option<ClientWithProfile>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let row = client.query_opt(
        r#"
        SELECT 
            c.id as customer_id,
            c.first_name,
            c.last_name,
            c.email,
            c.phone,
            c.address,
            c.city,
            c.loyalty_points,
            cp.date_of_birth,
            cp.gender,
            COALESCE(cp.newsletter_opt_in, false) as newsletter_opt_in,
            cp.profile_image_url,
            COALESCE(cp.account_status, 'Active') as account_status
        FROM customers c
        LEFT JOIN client_profiles cp ON c.id = cp.customer_id
        WHERE c.id = $1
        "#,
        &[&customer_id]
    ).await?;
    
    Ok(row.map(|r| ClientWithProfile {
        customer_id: r.get("customer_id"),
        first_name: r.get("first_name"),
        last_name: r.get("last_name"),
        email: r.get("email"),
        phone: r.get("phone"),
        address: r.get("address"),
        city: r.get("city"),
        loyalty_points: r.get("loyalty_points"),
        date_of_birth: r.get("date_of_birth"),
        gender: r.get("gender"),
        newsletter_opt_in: r.get("newsletter_opt_in"),
        profile_image_url: r.get("profile_image_url"),
        account_status: r.get("account_status"),
    }))
}

pub async fn update_client_profile(pool: &Pool, customer_id: i64, req: &UpdateClientProfileRequest) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    // Update customer table
    let mut customer_updates = vec![];
    if req.phone.is_some() { customer_updates.push("phone = COALESCE($1, phone)"); }
    if req.address.is_some() { customer_updates.push("address = COALESCE($2, address)"); }
    if req.city.is_some() { customer_updates.push("city = COALESCE($3, city)"); }
    
    if !customer_updates.is_empty() {
        customer_updates.push("updated_at = NOW()");
        let query = format!(
            "UPDATE customers SET {} WHERE id = $4",
            customer_updates.join(", ")
        );
        client.execute(&query, &[&req.phone, &req.address, &req.city, &customer_id]).await?;
    }
    
    // Upsert client profile
    client.execute(
        r#"
        INSERT INTO client_profiles (customer_id, date_of_birth, gender, newsletter_opt_in, profile_image_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        ON CONFLICT (customer_id) DO UPDATE SET
            date_of_birth = COALESCE(EXCLUDED.date_of_birth, client_profiles.date_of_birth),
            gender = COALESCE(EXCLUDED.gender, client_profiles.gender),
            newsletter_opt_in = COALESCE(EXCLUDED.newsletter_opt_in, client_profiles.newsletter_opt_in),
            profile_image_url = COALESCE(EXCLUDED.profile_image_url, client_profiles.profile_image_url),
            updated_at = NOW()
        "#,
        &[&customer_id, &req.date_of_birth, &req.gender, &req.newsletter_opt_in, &req.profile_image_url]
    ).await?;
    
    Ok(())
}

pub async fn get_pending_deletions(pool: &Pool) -> Result<Vec<ClientWithProfile>, Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    let rows = client.query(
        r#"
        SELECT 
            c.id as customer_id,
            c.first_name,
            c.last_name,
            c.email,
            c.phone,
            c.address,
            c.city,
            c.loyalty_points,
            cp.date_of_birth,
            cp.gender,
            COALESCE(cp.newsletter_opt_in, false) as newsletter_opt_in,
            cp.profile_image_url,
            cp.account_status
        FROM customers c
        JOIN client_profiles cp ON c.id = cp.customer_id
        WHERE cp.account_status = 'PendingDeletion'
        ORDER BY cp.deletion_requested_at ASC
        "#,
        &[]
    ).await?;
    
    let clients: Vec<ClientWithProfile> = rows.iter().map(|row| ClientWithProfile {
        customer_id: row.get("customer_id"),
        first_name: row.get("first_name"),
        last_name: row.get("last_name"),
        email: row.get("email"),
        phone: row.get("phone"),
        address: row.get("address"),
        city: row.get("city"),
        loyalty_points: row.get("loyalty_points"),
        date_of_birth: row.get("date_of_birth"),
        gender: row.get("gender"),
        newsletter_opt_in: row.get("newsletter_opt_in"),
        profile_image_url: row.get("profile_image_url"),
        account_status: row.get("account_status"),
    }).collect();
    
    Ok(clients)
}

pub async fn cancel_deletion(pool: &Pool, customer_id: i64) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    client.execute(
        r#"
        UPDATE client_profiles 
        SET account_status = 'Active', deletion_requested_at = NULL, updated_at = NOW()
        WHERE customer_id = $1
        "#,
        &[&customer_id]
    ).await?;
    
    // Reactivate user if linked
    client.execute(
        r#"
        UPDATE users SET is_active = true 
        WHERE id = (SELECT user_id FROM client_profiles WHERE customer_id = $1)
        "#,
        &[&customer_id]
    ).await?;
    
    Ok(())
}

pub async fn confirm_deletion(pool: &Pool, customer_id: i64) -> Result<(), Box<dyn std::error::Error>> {
    let client = pool.get().await?;
    
    // Mark as deleted (soft delete)
    client.execute(
        r#"
        UPDATE client_profiles 
        SET account_status = 'Deleted', deleted_at = NOW(), updated_at = NOW()
        WHERE customer_id = $1
        "#,
        &[&customer_id]
    ).await?;
    
    // Deactivate user permanently
    client.execute(
        r#"
        UPDATE users SET is_active = false 
        WHERE id = (SELECT user_id FROM client_profiles WHERE customer_id = $1)
        "#,
        &[&customer_id]
    ).await?;
    
    // Anonymize customer data for GDPR compliance
    client.execute(
        r#"
        UPDATE customers 
        SET first_name = 'Deleted', last_name = 'User', email = CONCAT('deleted_', id, '@deleted.com'),
            phone = NULL, address = NULL, updated_at = NOW()
        WHERE id = $1
        "#,
        &[&customer_id]
    ).await?;
    
    Ok(())
}
