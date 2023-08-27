run:
	go run cmd/zyreva/zyreva.go

call:
	buf curl \
	--schema buf.build/robinwongm/zyreva \
	--data '{"url": "https://github.com/skyzh/prisma-edge-vercel"}' \
	http://localhost:8080/zyreva.v1.GitService/Clone