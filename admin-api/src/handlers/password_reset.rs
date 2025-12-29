use actix_web::{web, HttpResponse, Responder};

async fn request_password_reset() -> impl Responder {
    HttpResponse::Ok().finish()
}

async fn reset_password() -> impl Responder {
    HttpResponse::Ok().finish()
}


pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/password-reset")
            .route("/request", web::post().to(request_password_reset))
            .route("/reset", web::post().to(reset_password))
    );
}