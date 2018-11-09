default: gitfs

fmt:
	go fmt

test:
	go test

install:
	go install

gitfs: install
	gitfs

