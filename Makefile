default: install

# moved out of deps to decrease build time
build-deps:
	GO111MODULE=off go get -u oss.indeed.com/go/go-groups
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	GO111MODULE=off go get -u github.com/gogo/protobuf/protoc-gen-gogo
	GO111MODULE=off go get -u github.com/mitchellh/gox

fmt:
	go-groups -w .
	gofmt -s -w .

deps:
	go get -v ./...

test:
	go vet ./...
	golint -set_exit_status ./...
	go test -v -race ./...

install:
	go install

deploy:
	mkdir -p bin
	gox -os="darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...
	gox -os="linux" -arch="amd64 386 arm arm64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...
