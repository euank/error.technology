
all:
	CGO_ENABLED=0 go build -installsuffix cgo -a -ldflags "-s" -o server server.go
	docker build -f Dockerfile.server -t euank/error-tech-server .

push:
	docker push euank/error-tech-server
