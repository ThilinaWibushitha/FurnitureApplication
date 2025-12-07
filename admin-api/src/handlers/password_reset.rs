use actix_web::{web, HttpResponse, post};
use deadpool_postgres::Pool;
use crate::models::ResetPasswordRequest;
use crate::services::auth_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(reset_password);
}

#[post("/auth/reset-password")]
async fn reset_password(pool: web::Data<Pool>, body: web::Json<ResetPasswordRequest>) -> HttpResponse {
    match auth_service::reset_password(&pool, &body.token, &body.new_password).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({
            "message": "Password reset successfully"
        })),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({
            "error": e.to_string()
        }))
    }
}
