Contributing
===

### prerequisites

* [Opening a Pull Request](https://help.github.com/articles/about-pull-requests/)
* Review [Code of Conduct](CODE_OF_CONDUCT.md)

### as a technical writer

Like many projects, technical documentation is often needed.
Instead of restricting contributions strictly to engineers, we want to encourage contributions from all walks of life.
Documentation is open to contributions from the public.
They can be found under the `docs` folder and hosted on github here: https://mjpitz.github.io/gitfs
Changes to documentation do not require tickets in advanced, but do require pull requests.

### as a developer / engineer

Before diving into committing and developing on the code, let's cover some basic prerequisites specific to developers.
As a general user or technical writer, you shouldn't need to worry about installing these dependencies. 

**Prerequisite: Golang 11**

In addition to the libfuse module from the README, you will also need [go 11](https://golang.org/doc/install) installed.

```
wget -O go1.11.tar.gz https://dl.google.com/go/go1.11.linux-amd64.tar.gz   # linux
wget -O go1.11.tar.gz https://dl.google.com/go/go1.11.darwin-amd64.tar.gz  # osx

tar -xvf go1.11.tar.gz
sudo mv go /usr/local
rm go1.11.tar.gz

# add these to your environment
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```
  
**Prerequisite: Protocol Buffers**

Lastly, you will need the [protocol buffers compiler](https://developers.google.com/protocol-buffers/docs/downloads) should you plan on doing any development on this project.

```
wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip   #linux
wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-osx-x86_64.zip     #osx

unzip -d protoc protoc.zip
sudo mv protoc /usr/local

# add this to your environment 
export PATH=/usr/local/protoc/bin:$PATH
```

We also use the [gogo](https://github.com/gogo/protobuf) generation for protobuf, so you will need the compiler plugin installed.

```
go get github.com/gogo/protobuf/protoc-gen-gogo
```

**Committing Code**

All work items should have a corresponding issue under the github issue manager.
The first commit associated tickets must be prefixed with the github issue number as such: `gh-{issue}: <message>`.
All following commits to the branch do not require that prefix as they will be squashed before merge.
Your first commit will be used as the commit message, unless project owners desire a more clear message.

**Dependency Management**

This project uses go modules for managing it's dependencies.
When working on this project, always make sure you run `go mod tidy` before pushing.
This will ensure that we do not pull in more dependencies then needed for compilation. 

```
go get -u <dependency>
go mod tidy
```

**Testing**

```
make test
```

**Building / Installing**

`make install` can be used to compile a development version for the current platform.
The binary will be placed in your `$GOPATH/bin` directory.

`make deploy` can be used to cross compile the application for different platforms using `gox`.

**Issue Reporting**

If you believe you've found an issue, first check existing issues for a similar issue.
If no issue has been opened, then please submit a ticket.
After a ticket is opened, you are welcome to open a PR for review.
All development on this issue must follow the guidelines above.

### as a project owner

As a project owner, you will be responsible for managing releases of the application.
To make life easy, the project uses `npm` and `package.json` for version management.
This is because managing versions using this method is trivial.

```
npm version patch|minor|major
git push
git push --tags
```

A passing build is required for publishing.
There are currently no hooks that prevent this, so it is your responsibility to check the build. 
