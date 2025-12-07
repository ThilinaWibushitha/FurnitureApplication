use actix_web::{web, HttpResponse, get, post, put, delete};
use deadpool_postgres::Pool;
use crate::models::{CreateItemRequest, UpdateItemRequest};
use crate::services::item_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(get_all_items)
       .service(get_item)
       .service(create_item)
       .service(update_item)
       .service(delete_item);
}

#[get("/items")]
async fn get_all_items(pool: web::Data<Pool>) -> HttpResponse {
    match item_service::get_all_items(&pool).await {
        Ok(items) => HttpResponse::Ok().json(items),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/items/{id}")]
async fn get_item(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    match item_service::get_item_by_id(&pool, path.into_inner()).await {
        Ok(Some(item)) => HttpResponse::Ok().json(item),
        Ok(None) => HttpResponse::NotFound().json(serde_json::json!({"error": "Item not found"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[post("/items")]
async fn create_item(pool: web::Data<Pool>, body: web::Json<CreateItemRequest>) -> HttpResponse {
    match item_service::create_item(&pool, &body).await {
        Ok(item) => HttpResponse::Created().json(item),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[put("/items/{id}")]
async fn update_item(pool: web::Data<Pool>, path: web::Path<i64>, body: web::Json<UpdateItemRequest>) -> HttpResponse {
    match item_service::update_item(&pool, path.into_inner(), &body).await {
        Ok(item) => HttpResponse::Ok().json(item),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[delete("/items/{id}")]
async fn delete_item(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    match item_service::delete_item(&pool, path.into_inner()).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Item deleted"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}
