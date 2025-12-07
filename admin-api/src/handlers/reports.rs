use actix_web::{web, HttpResponse, get};
use deadpool_postgres::Pool;
use crate::services::report_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(get_daily_report)
       .service(get_monthly_report);
}

#[get("/reports/daily")]
async fn get_daily_report(pool: web::Data<Pool>) -> HttpResponse {
    match report_service::get_daily_sales(&pool).await {
        Ok(report) => HttpResponse::Ok().json(report),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}

#[get("/reports/monthly")]
async fn get_monthly_report(pool: web::Data<Pool>) -> HttpResponse {
    match report_service::get_monthly_sales(&pool).await {
        Ok(report) => HttpResponse::Ok().json(report),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}
