use actix_web::{web, HttpResponse, post};
use deadpool_postgres::Pool;
use crate::models::{LoginRequest, LoginResponse, UserResponse, CreateUserRequest};
use crate::services::auth_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(login)
       .service(create_admin);
}

#[post("/auth/login")]
async fn login(pool: web::Data<Pool>, body: web::Json<LoginRequest>) -> HttpResponse {
    match auth_service::authenticate(&pool, &body).await {
        Ok(response) => HttpResponse::Ok().json(response),
        Err(e) => HttpResponse::Unauthorized().json(serde_json::json!({
            "error": e.to_string()
        }))
    }
}

#[post("/auth/create-admin")]
async fn create_admin(pool: web::Data<Pool>, body: web::Json<CreateUserRequest>) -> HttpResponse {
    let mut request = body.into_inner();
    request.account_type = "admin".to_string();
    
    match auth_service::create_user(&pool, &request).await {
        Ok(user) => HttpResponse::Created().json(user),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))
    }
}
