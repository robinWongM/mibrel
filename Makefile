build:
	go build cmd/zyreva/zyreva.go

call:
	curl --header "Content-Type: application/json" \
	--data '{"url": "https://github.com/skyzh/prisma-edge-vercel"}' \
	http://localhost:8080/zyreva.v1.GitService/Clone