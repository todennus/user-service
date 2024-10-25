start-rest:
	go run ./cmd/main.go rest

start-grpc:
	go run ./cmd/main.go grpc

docker-build:
	docker build -t todennus/user-service -f ./build/package/Dockerfile .
