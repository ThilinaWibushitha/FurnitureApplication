FROM rust:1.77.2-slim-bullseye

WORKDIR /usr/src/admin-api

COPY ./admin-api/Cargo.toml ./admin-api/Cargo.lock ./

RUN mkdir ./src && \
    echo 'fn main() {println!("Dummy main for dependency caching")}' > ./src/main.rs && \
    cargo build --release && \
    rm -rf ./src

COPY ./admin-api/src ./src
COPY ./admin-api/database ./database

RUN touch .env

RUN cargo build --release

EXPOSE 8080

CMD ["./target/release/admin-api"]

