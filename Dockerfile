FROM rust:1.83-slim-bullseye

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/admin-api

# Copy cargo files for dependency caching
COPY ./admin-api/Cargo.toml ./admin-api/Cargo.lock ./

# Build dependencies with dummy main
RUN mkdir ./src && \
    echo 'fn main() {println!("Dummy main for dependency caching")}' > ./src/main.rs && \
    cargo build --release && \
    rm -rf ./src

# Copy actual source code
COPY ./admin-api/src ./src
COPY ./admin-api/database ./database

RUN touch .env

# Clean and rebuild the actual application
RUN cargo clean && cargo build --release

EXPOSE 5200

CMD ["./target/release/admin-api"]

