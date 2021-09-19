# Build verifier
FROM --platform=linux/amd64 rust:latest AS build_verifier

RUN cd /root && git clone https://github.com/pj-aias/aias-verifier.git

WORKDIR /root/aias-verifier

#RUN rustup target add x86_64-unknown-linux-gnu
RUN cargo fetch
RUN cargo build


# Main API Server
FROM --platform=linux/amd64 golang:latest

WORKDIR /app

COPY --from=build_verifier /root/aias-verifier/target/release/aias-verifier ./

CMD go run .