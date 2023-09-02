####################################################################################################
## Builder
####################################################################################################
FROM rust:1.72.0-alpine3.18 AS builder

RUN apk add --no-cache musl-dev pkgconfig openssl openssl-dev
ENV OPENSSL_LIB_DIR=/usr/lib
ENV OPENSSL_INCLUDE_DIR=/usr/include

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
    cargo build --verbose  && \
    cp /zyreva/target/debug/server /zyreva/server

####################################################################################################
## Final image
####################################################################################################
FROM alpine:3.18.3

RUN apk add --no-cache openssl

# Import from builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

WORKDIR /zyreva

# Copy our build
COPY --from=builder /zyreva/server ./

# Use an unprivileged user.
USER zyreva:zyreva

CMD ["/zyreva/server"]