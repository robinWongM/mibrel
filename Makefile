bindgen-go:
	uniffi-bindgen-go crates/nixpacks-ffi/src/nixpacks.udl -o internal/gen

build-rust: $(shell find crates -type f)
	cargo build

build: build-rust bindgen-go
	LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-}:${PWD}/target/debug/" \
	CGO_LDFLAGS="-lnixpacks_ffi -L${PWD}/target/debug/ -lm -ldl" \
	CGO_ENABLED=1 \
	go build cmd/zyreva/zyreva.go

call:
	curl --header "Content-Type: application/json" \
	--data '{"url": "https://github.com/skyzh/prisma-edge-vercel"}' \
	http://localhost:8080/zyreva.v1.GitService/Clone