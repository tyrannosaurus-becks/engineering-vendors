fmt:
	go fmt ./...

test:
	go test -v ./...

dev:
	cd cmd/vendors && go build && go install && cd -