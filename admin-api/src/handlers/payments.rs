use actix_web::{web, HttpResponse, get, post, put, delete};
use deadpool_postgres::Pool;
use crate::models::payment::{CreatePaymentRequest, UpdatePaymentRequest, ProcessRefundRequest};
use crate::services::payment_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(list_payments)
       .service(get_payment)
       .service(create_payment)
       .service(update_payment)
       .service(get_order_payments)
       .service(list_cancellations)
       .service(process_refund);
}

#[get("/payments")]
async fn list_payments(pool: web::Data<Pool>, query: web::Query<PaymentFilter>) -> HttpResponse {
    match payment_service::list_payments(&pool, &query).await {
        Ok(payments) => HttpResponse::Ok().json(payments),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[derive(Debug, serde::Deserialize)]
pub struct PaymentFilter {
    pub status: Option<String>,
    pub customer_id: Option<i64>,
    pub order_id: Option<i64>,
}

#[get("/payments/{id}")]
async fn get_payment(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let payment_id = path.into_inner();
    match payment_service::get_payment(&pool, payment_id).await {
        Ok(Some(payment)) => HttpResponse::Ok().json(payment),
        Ok(None) => HttpResponse::NotFound().json(serde_json::json!({"error": "Payment not found"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/orders/{id}/payments")]
async fn get_order_payments(pool: web::Data<Pool>, path: web::Path<i64>) -> HttpResponse {
    let order_id = path.into_inner();
    match payment_service::get_order_payments(&pool, order_id).await {
        Ok(payments) => HttpResponse::Ok().json(payments),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[post("/payments")]
async fn create_payment(pool: web::Data<Pool>, body: web::Json<CreatePaymentRequest>) -> HttpResponse {
    match payment_service::create_payment(&pool, &body).await {
        Ok(payment) => HttpResponse::Created().json(payment),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[put("/payments/{id}")]
async fn update_payment(pool: web::Data<Pool>, path: web::Path<i64>, body: web::Json<UpdatePaymentRequest>) -> HttpResponse {
    let payment_id = path.into_inner();
    match payment_service::update_payment(&pool, payment_id, &body).await {
        Ok(payment) => HttpResponse::Ok().json(payment),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/cancellations")]
async fn list_cancellations(pool: web::Data<Pool>) -> HttpResponse {
    match payment_service::list_cancellations(&pool).await {
        Ok(cancellations) => HttpResponse::Ok().json(cancellations),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[post("/cancellations/{id}/process")]
async fn process_refund(pool: web::Data<Pool>, path: web::Path<i64>, body: web::Json<ProcessRefundRequest>) -> HttpResponse {
    let cancellation_id = path.into_inner();
    match payment_service::process_refund(&pool, cancellation_id, &body).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Refund processed successfully"})),
        Err(e) => HttpResponse::BadRequest().json(serde_json::json!({"error": e.to_string()}))
    }
}
