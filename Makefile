start-rest:
	go run ./cmd/main.go rest

start-grpc:
	go run ./cmd/main.go grpc

start-swagger:
	go run ./cmd/main.go swagger

docker-build:
	docker build -t todennus/user-service -f ./build/package/Dockerfile .
