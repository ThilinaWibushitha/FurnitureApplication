use deadpool_postgres::{Config, Pool, Runtime};
use std::env;
use tokio_postgres::NoTls;
use url::Url;

pub async fn create_pool() -> Pool {
    let mut cfg = Config::new();
    
    let database_url = env::var("DATABASE_URL")
        .unwrap_or_else(|_| "postgres://postgres:1881@localhost:5432/furniture_db".to_string());
    
    // Parse DATABASE_URL
    let url = Url::parse(&database_url).expect("Invalid DATABASE_URL");
    
    cfg.host = url.host_str().map(String::from);
    cfg.port = url.port();
    cfg.user = Some(url.username().to_string());
    cfg.password = url.password().map(String::from);
    cfg.dbname = Some(url.path().trim_start_matches('/').to_string());

    cfg.create_pool(Some(Runtime::Tokio1), NoTls)
        .expect("Failed to create database pool")
}
