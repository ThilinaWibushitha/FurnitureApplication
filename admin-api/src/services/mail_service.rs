use lettre::{Message, SmtpTransport, Transport};
use lettre::transport::smtp::authentication::Credentials;
use std::env;

pub async fn send_email(to: &str, subject: &str, body: &str) -> Result<(), Box<dyn std::error::Error>> {
    let mail_server = env::var("MAIL_SERVER").unwrap_or_else(|_| "smtp.gmail.com".to_string());
    let mail_user = env::var("MAIL_USER")?;
    let mail_password = env::var("MAIL_PASSWORD")?;

    let email = Message::builder()
        .from(mail_user.parse()?)
        .to(to.parse()?)
        .subject(subject)
        .body(body.to_string())?;

    let creds = Credentials::new(mail_user, mail_password);

    let mailer = SmtpTransport::relay(&mail_server)?
        .credentials(creds)
        .build();

    mailer.send(&email)?;
    Ok(())
}
