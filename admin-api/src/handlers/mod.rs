pub mod auth;
pub mod users;
pub mod items;
pub mod orders;
pub mod reports;
pub mod mail;
pub mod payments;
pub mod clients;
pub mod password_reset;

use actix_web::web;

pub fn config(cfg: &mut web::ServiceConfig) {
    cfg.service(
        web::scope("/api/v1")
            .configure(auth::config)
            .configure(users::config)
            .configure(items::config)
            .configure(orders::config)
            .configure(reports::config)
            .configure(mail::config)
            .configure(payments::config)
            .configure(clients::config)
            .configure(password_reset::config)
    );
}
