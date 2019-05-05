default: install

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
	# produce binaries for the common platforms, both amd64 and i386.
	# can't handle windows since it looks like it might need the fuse lib.
	gox -os="linux darwin" -arch="amd64 386" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...
	# special build for linux arm systems, both 32 and 64 bit.
	# can't build arm64 yet. some code seems to allude to it, but the gox tool still reports an error.
	gox -os="linux" -arch="arm" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...
	GOOS=linux GOARCH=arm64 go build -o bin/gitfs_linux_arm64
