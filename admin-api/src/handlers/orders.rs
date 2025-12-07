use actix_web::{web, HttpResponse, get, put};
use deadpool_postgres::Pool;
use crate::models::UpdateOrderRequest;
use crate::services::order_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(get_all_orders)
       .service(get_order)
       .service(update_order);
}

#[get("/orders")]
async fn get_all_orders(pool: web::Data<Pool>) -> HttpResponse {
    match order_service::get_all_orders(&pool).await {
        Ok(orders) => HttpResponse::Ok().json(orders),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/orders/{id}")]
async fn get_order(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    match order_service::get_order_by_id(&pool, path.into_inner()).await {
        Ok(Some(order)) => HttpResponse::Ok().json(order),
        Ok(None) => HttpResponse::NotFound().json(serde_json::json!({"error": "Order not found"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[put("/orders/{id}")]
async fn update_order(pool: web::Data<Pool>, path: web::Path<i64>, body: web::Json<UpdateOrderRequest>) -> HttpResponse {
    match order_service::update_order(&pool, path.into_inner(), &body).await {
        Ok(order) => HttpResponse::Ok().json(order),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}
