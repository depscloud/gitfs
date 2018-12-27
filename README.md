## gitfs

_Better repository management for the everyday developer._

When starting a job at any company, you often spend a fair amount setting up your development environment.
This involves installing dependent tooling, cloning repositories, amongst many other things.
`gitfs` was developed to ease some of that process by removing the need to manually clone down every repository.
Instead, it creates a virtual directory structure around git urls, making it easy to navigate sources for projects.

## Support

Support on this library is limited right now, but open and under for active development.
Below, you will find some of the current limitations / capabilities for the project.

[![Build Status](https://travis-ci.org/mjpitz/gitfs.svg?branch=master)](https://travis-ci.org/mjpitz/gitfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/mjpitz/gitfs)](https://goreportcard.com/report/github.com/mjpitz/gitfs)

### Limited Remote Implementations

Currently, it supports a generic endpoint, but are open to supporting new remotes for various integrations.
We are planning to support Github, Gitlab, and Bitbucket in the short term.
PR's for each of these remotes based on their [config](pkg/config/config.proto) definitions are welcome.

- [ ] [gh-3: Add remote implementation for gitlab](https://github.com/mjpitz/gitfs/issues/3)
- [ ] [gh-2: Add remote implementation for bitbucket](https://github.com/mjpitz/gitfs/issues/2)
- [x] [gh-1: Add remote implementation for github](https://github.com/mjpitz/gitfs/issues/1)

### In Memory Storage

For the proof of concept, this project leverages an in memory file system to store the cloned repositories.
This leverages the billy library, making it really easy to swap the underlying git filesystem for an alternative implementation.
Support for on persistent stores will be added for better long term support.

- [ ] [gh-4: Add support for both memfs and osfs](https://github.com/mjpitz/gitfs/issues/4)

### Shallow Clones

It's important to note that when navigating to a repository, we perform a shallow clone on the fly.
The shallow clone helps keep your footprint to a minimum while keeping developers waiting to a minimum.
We do plan on adding support for full clones in a background process upon request.

- [ ] [gh-5: Support non-shallow clones](https://github.com/mjpitz/gitfs/issues/5)

## Getting Started

Getting started with this tool can be a little bit of a hassle for now.
It will get better over time as development continues on the project.

### Prerequisites

**libfuse**

You'll need to install FUSE on your system in order to use this project.
For directions on how to install, see the [libfuse](https://github.com/libfuse/libfuse) github page.

```
apt-get install libfuse-dev   # linux
brew cask install osxfuse     # osx
```

**Golang 11**

In addition to the libfuse module, you will also need [go 11](https://golang.org/doc/install) installed.

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
  
**Protocol Buffers**

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
 
### Building from Source

This project was written using Golang 11's modules.
As a result, we require that you use the modules component for installation.
The 5 commands below should install a working copy to your `${GOPATH}/bin` path.

```bash
mkdir -p ${GOPATH}/src/github.com/mjpitz
cd ${GOPATH}/src/github.com/mjpitz
git clone git@github.come:mjpitz/gitfs.git
cd gitfs
GO11MODULES=on go install
```

### Configuring gitfs

Once the binary has been installed, you will need to configure `gitfs`.
By default, `gitfs` looks for a configuration file at `${HOME}/.gitfs/config.yml`.
The location of the config file can be changed using the `-config` parameter.
Below, I have provided a snippet on the configuration required for using github as a remote. 

```yaml
# cat ${HOME}/.gitfs/config.yml
mount: "${HOME}/Development/code"
accounts:
  - github:
      oauth2:
        token: <oauth_token>
      users:
      - <username>
```

### Running gitfs

Once your configuration has been created and values have been replaced, you can run `gitfs`.
Note that running `gitfs` will spawn the filesystem service as a blocking process in your active terminal window.
You can run this as a systemctl daemon, but I find having access to the logs more readily at this time is a bit nicer.

To run, simply execute `gitfs`.
When run successfully, you should see an output message similar to the one below.

```
[mjpitz@mjpitz ~/Development/go/src/github.com/mjpitz/gitfs master]
$ gitfs
INFO[0000] [main] configured mount point: /Users/mjpitz/Development/code 
INFO[0000] [main] fetching repositories                 
INFO[0000] [remotes.github] processing organizations for user: mjpitz 
INFO[0000] [remotes.github] processing repositories for user: mjpitz 
INFO[0001] [main] parsing repositories into a directory structure 
INFO[0001] [main] attempting to mount /Users/mjpitz/Development/code 
INFO[0001] [main] now serving file system  
``` 

In a separate window, navigate to your configured `mount` point.
For me, I have this set `${HOME}/Development/code`.
From this directory, you should be able to navigate your projects on Github.

```
[mjpitz@mjpitz ~/Development/code]
$ ls
github.com

[mjpitz@mjpitz ~/Development/code]
$ cd github.com/

[mjpitz@mjpitz ~/Development/code/github.com]
$ ls
mjpitz

[mjpitz@mjpitz ~/Development/code/github.com]
$ cd mjpitz/

[mjpitz@mjpitz ~/Development/code/github.com/mjpitz]
$ ls
OpenGrok	   grpc-java	    laas	rpi
consul-api	   grpc.github.io   mjpitz.com	seo-portal
docker-clickhouse  grpcsh	    mp		serverless-plugin-simulate
docker-utils	   hbase-docker     okhttp	simple-daemon-node
dotfiles	   idea-framework   proctor	spring-config-repo
generator-idea	   java-gitlab-api  proctorjs
gitfs		   jgrapht	    proto2-3

[mjpitz@mjpitz ~/Development/code/github.com/mjpitz]
$ cd grpcsh

[mjpitz@mjpitz ~/Development/code/github.com/mjpitz/grpcsh master]
$ ls
README.md  bin	package-lock.json  package.json  yarn.lock
```
