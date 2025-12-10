use actix_web::{web, HttpResponse, post};
use deadpool_postgres::Pool;
use crate::models::{LoginRequest, CreateUserRequest, RequestPasswordResetRequest, ResetPasswordRequest};
use crate::services::auth_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(login)
       .service(create_admin)
       .service(request_password_reset)
       .service(reset_password);
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

#[post("/auth/main-admin/request-password-reset")]
async fn request_password_reset(pool: web::Data<Pool>, body: web::Json<RequestPasswordResetRequest>) -> HttpResponse {
    match auth_service::request_password_reset(&pool, &body.email).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({
            "message": "Password reset link sent to your email"
        })),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))
    }
}

#[post("/auth/main-admin/reset-password")]
async fn reset_password(pool: web::Data<Pool>, body: web::Json<ResetPasswordRequest>) -> HttpResponse {
    match auth_service::reset_password(&pool, &body.token, &body.password).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({
            "message": "Password reset successfully"
        })),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({
            "error": e.to_string()
        }))
    }
}
