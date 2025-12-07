use actix_web::{web, HttpResponse, post};
use serde::{Deserialize, Serialize};
use crate::services::mail_service;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(send_mail);
}

#[derive(Debug, Deserialize, Serialize)]
pub struct SendMailRequest {
    pub to: String,
    pub subject: String,
    pub body: String,
}

#[post("/mail/send")]
async fn send_mail(body: web::Json<SendMailRequest>) -> HttpResponse {
    match mail_service::send_email(&body.to, &body.subject, &body.body).await {
        Ok(_) => HttpResponse::Ok().json(serde_json::json!({"message": "Email sent successfully"})),
        Err(e) => HttpResponse::InternalServerError().json(serde_json::json!({"error": e.to_string()}))
    }
}
