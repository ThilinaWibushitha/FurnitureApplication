use actix_cors::Cors;
use actix_web::{web, App, HttpServer, middleware};
use dotenv::dotenv;
use std::env;

mod config;
mod handlers;
mod models;
mod services;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();
    env_logger::init();

    let host = env::var("SERVER_HOST").unwrap_or_else(|_| "127.0.0.1".to_string());
    let port = env::var("SERVER_PORT").unwrap_or_else(|_| "5200".to_string());
    let addr = format!("{}:{}", host, port);

    let pool = config::db::create_pool().await;

    println!("🚀 Admin API running at http://{}", addr);

    HttpServer::new(move || {
        let cors = Cors::default()
            .allow_any_origin()
            .allow_any_method()
            .allow_any_header();

        App::new()
            .app_data(web::Data::new(pool.clone()))
            .wrap(cors)
            .wrap(middleware::Logger::default())
            .configure(handlers::config)
    })
    .bind(&addr)?
    .run()
    .await
}
