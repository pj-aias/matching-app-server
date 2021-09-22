# Build verifier
FROM rust:latest AS build_verifier

RUN cd /root && git clone https://github.com/pj-aias/aias-verifier.git

WORKDIR /root/aias-verifier

#RUN rustup target add x86_64-unknown-linux-gnu
RUN cargo fetch
RUN cargo build --release


# Main API Server
FROM golang:latest

WORKDIR /app

COPY --from=build_verifier /root/aias-verifier/target/release/aias-verifier /usr/local/bin/aias-verifier

CMD go run .