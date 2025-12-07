use actix_web::{web, HttpResponse, get, put};
use log::error;
use deadpool_postgres::Pool;
use crate::models::payment::{UpdateClientProfileRequest};
use crate::services::client_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(list_clients)
       .service(get_client)
       .service(update_client_profile)
       .service(get_pending_deletions)
       .service(cancel_deletion)
       .service(confirm_deletion);
}

#[derive(Debug, serde::Deserialize)]
pub struct ClientFilter {
    pub status: Option<String>,
    pub search: Option<String>,
}

#[get("/clients")]
async fn list_clients(pool: web::Data<Pool>, query: web::Query<ClientFilter>) -> HttpResponse {
    match client_service::list_clients(&pool, &query).await {
        Ok(clients) => HttpResponse::Ok().json(clients),
        Err(e) => {
            error!("Error listing clients: {:?}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
        }
    }
}

#[get("/clients/{id}")]
async fn get_client(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let customer_id = path.into_inner();
    match client_service::get_client(&pool, customer_id).await {
        Ok(Some(client)) => HttpResponse::Ok().json(client),
        Ok(None) => HttpResponse::NotFound().json(serde_json::json!({"error": "Client not found"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[put("/clients/{id}")]
async fn update_client_profile(
    pool: web::Data<Pool>, 
    path: web::Path<i64>, 
    body: web::Json<UpdateClientProfileRequest>
) -> HttpResponse {
    let customer_id = path.into_inner();
    match client_service::update_client_profile(&pool, customer_id, &body).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Client profile updated"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/clients/pending-deletions")]
async fn get_pending_deletions(pool: web::Data<Pool>) -> HttpResponse {
    match client_service::get_pending_deletions(&pool).await {
        Ok(clients) => HttpResponse::Ok().json(clients),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[actix_web::post("/clients/{id}/cancel-deletion")]
async fn cancel_deletion(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let customer_id = path.into_inner();
    match client_service::cancel_deletion(&pool, customer_id).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Deletion cancelled, account restored"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[actix_web::post("/clients/{id}/confirm-deletion")]
async fn confirm_deletion(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let customer_id = path.into_inner();
    match client_service::confirm_deletion(&pool, customer_id).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Account permanently deleted"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}
