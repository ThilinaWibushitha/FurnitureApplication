use deadpool_postgres::{Config, Runtime};
use tokio_postgres::NoTls;
use std::env;
use std::fs;
use bcrypt::{hash, DEFAULT_COST};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv::dotenv().ok();
    println!("🔌 Connecting to database...");

    let mut cfg = Config::new();
    let database_url = env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://postgres:1881@localhost:5432/furniture_db".to_string());
    
    // Parse DATABASE_URL manually/simply for config
    // (Assuming standard format for simplicity or relying on deadpool's url parsing if it supported it directly, 
    // but deadpool-postgres Config is struct based. Let's use the same logic as db.rs if we want to be safe, 
    // or just use tokio_postgres::connect directly for a script)
    
    // Actually, to avoid circular dependencies or complex config parsing again, 
    // let's just use tokio_postgres directly for this script.
    
    let (client, connection) = tokio_postgres::connect(&database_url, NoTls).await?;

    // The connection object performs the actual communication with the database,
    // so spawn it off to run in the background.
    tokio::spawn(async move {
        if let Err(e) = connection.await {
            eprintln!("connection error: {}", e);
        }
    });

    println!("✅ Connected.");

    // Read Schema
    println!("📜 Applying Schema...");
    let schema_sql = fs::read_to_string("database/schema.sql")?;
    client.batch_execute(&schema_sql).await?;
    println!("✅ Schema applied.");

    // Read Seed
    println!("🌱 Seeding Data...");
    let seed_sql = fs::read_to_string("database/seed.sql")?;
    client.batch_execute(&seed_sql).await?;
    println!("✅ Data seeded.");

    // Create Admin User
    // Create Main Admin User
    println!("👤 Creating Main Admin User (main@admin.com / 1234)...");
    let password_hash = hash("1234", DEFAULT_COST)?;
    
    let rows_updated = client.execute(
        "INSERT INTO users (username, email, password_hash, full_name, account_type, is_active, is_verified) 
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         ON CONFLICT (username) DO UPDATE SET password_hash = $3, email = $2, account_type = $5",
        &[&"main_admin", &"main@admin.com", &password_hash, &"Main Administrator", &"main_admin", &true, &true]
    ).await?;
    
    if rows_updated > 0 {
        println!("✅ Main Admin user created/updated.");
    } else {
        println!("ℹ️ Main Admin user already exists.");
    }
    
    // Create Standard Admin User for testing
    println!("👤 Creating Standard Admin User (admin@admin.com / 1234)...");
    let admin_rows = client.execute(
        "INSERT INTO users (username, email, password_hash, full_name, account_type, is_active, is_verified) 
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         ON CONFLICT (username) DO UPDATE SET password_hash = $3, email = $2, account_type = $5",
        &[&"admin", &"admin@admin.com", &password_hash, &"Standard Admin", &"admin", &true, &true]
    ).await?;

    if admin_rows > 0 {
        println!("✅ Standard Admin user created/updated.");
    }

    println!("✨ Database setup complete!");
    Ok(())
}
