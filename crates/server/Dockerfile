####################################################################################################
## Builder
####################################################################################################
FROM rust:1.72.0-slim-buster AS builder

RUN rustup target add x86_64-unknown-linux-musl
RUN apt update && apt install -y musl-tools musl-dev
RUN update-ca-certificates

# Create appuser
ENV USER=zyreva
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /zyreva

COPY ./Cargo* .
COPY ./crates ./crates

RUN --mount=type=cache,target=/usr/local/cargo/registry \
    --mount=type=cache,target=/zyreva/target \
    cargo build --target x86_64-unknown-linux-musl --release && \
    cp /zyreva/target/x86_64-unknown-linux-musl/release/server /zyreva/server

####################################################################################################
## Final image
####################################################################################################
FROM alpine:3.18.3

# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /zyreva

# Copy our build
COPY --from=builder /zyreva/server ./

# Use an unprivileged user.
USER zyreva:zyreva

CMD ["/zyreva/server"]