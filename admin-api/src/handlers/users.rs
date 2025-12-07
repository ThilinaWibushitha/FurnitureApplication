use actix_web::{web, HttpResponse, get, post, put, delete};
use deadpool_postgres::Pool;
use crate::models::{ApproveUserRequest, ResolvePasswordChangeRequest};
use crate::services::user_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(get_all_users)
       .service(get_admins)
       .service(get_pending_users)
       .service(get_password_requests)
       .service(get_user)
       .service(approve_user)
       .service(resolve_password_request)
       .service(update_user_profile)
       .service(deactivate_user);
}

#[get("/users")]
async fn get_all_users(pool: web::Data<Pool>) -> HttpResponse {
    match user_service::get_all_users(&pool).await {
        Ok(users) => HttpResponse::Ok().json(users),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/users/admins")]
async fn get_admins(pool: web::Data<Pool>) -> HttpResponse {
    match user_service::get_admins(&pool).await {
        Ok(users) => HttpResponse::Ok().json(users),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/users/pending")]
async fn get_pending_users(pool: web::Data<Pool>) -> HttpResponse {
    match user_service::get_pending_users(&pool).await {
        Ok(users) => HttpResponse::Ok().json(users),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/users/password-requests")]
async fn get_password_requests(pool: web::Data<Pool>) -> HttpResponse {
    match user_service::get_password_change_requests(&pool).await {
        Ok(requests) => HttpResponse::Ok().json(requests),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/users/{id}")]
async fn get_user(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let user_id = path.into_inner();
    match user_service::get_user_by_id(&pool, user_id).await {
        Ok(Some(user)) => HttpResponse::Ok().json(user),
        Ok(None) => HttpResponse::NotFound().json(serde_json::json!({"error": "User not found"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[post("/users/approve")]
async fn approve_user(pool: web::Data<Pool>, body: web::Json<ApproveUserRequest>) -> HttpResponse {
    match user_service::approve_user(&pool, body.user_id, body.approved).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "User updated"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[post("/users/password-requests/resolve")]
async fn resolve_password_request(pool: web::Data<Pool>, body: web::Json<ResolvePasswordChangeRequest>) -> HttpResponse {
    match user_service::resolve_password_change_request(&pool, body.request_id, body.approved).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Request resolved"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[put("/users/{id}/profile")]
async fn update_user_profile(pool: web::Data<Pool>, path: web::Path<i64>, body: web::Json<serde_json::Value>) -> HttpResponse {
    let user_id = path.into_inner();
    match user_service::update_profile(&pool, user_id, &body).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Profile updated"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[delete("/users/{id}")]
async fn deactivate_user(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let user_id = path.into_inner();
    match user_service::deactivate_user(&pool, user_id).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "User deactivated"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}
