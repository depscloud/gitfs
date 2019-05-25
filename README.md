## gitfs

_Better repository management for the everyday developer._

When starting a job at any company, you often spend a fair amount setting up your development environment.
This involves installing dependent tooling, cloning repositories, amongst many other things.
`gitfs` was developed to ease some of that process by removing the need to manually clone down every repository.
Instead, it creates a virtual directory structure around git urls, making it easy to navigate sources for projects.

## Support

Support on this library is limited right now, but open to and under active development.
Below, you will find some of the current limitations / capabilities for the project.

![GitHub](https://img.shields.io/github/license/deps-cloud/gitfs.svg)
[![Build Status](https://travis-ci.com/deps-cloud/gitfs.svg?branch=master)](https://travis-ci.com/deps-cloud/gitfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/deps-cloud/gitfs)](https://goreportcard.com/report/github.com/deps-cloud/gitfs)

### Limited Remote Implementations

Currently, gitfs supports a generic endpoint, but are open to supporting new remotes for various integrations.
We are planning to support Github, Gitlab, and Bitbucket in the short term.
PR's for each of these remotes based on their [config](pkg/config/config.proto) definitions are welcome.

- [x] [gh-3: Add remote implementation for gitlab](https://github.com/deps-cloud/gitfs/issues/3)
- [ ] [gh-2: Add remote implementation for bitbucket](https://github.com/deps-cloud/gitfs/issues/2)
- [x] [gh-1: Add remote implementation for github](https://github.com/deps-cloud/gitfs/issues/1)

## Getting Started

Getting started with this tool can be a little bit of a hassle for now.
It will get better over time as development continues on the project.

### Prerequisites

**libfuse**

You'll need to install FUSE on your system in order to use this project.
For directions on how to install, see the [libfuse](https://github.com/libfuse/libfuse) github page.

```
sudo apt-get install libfuse-dev   # Ubuntu/Debian
sudo pacman -S fuse                # Manjaro/Arch

brew cask install osxfuse          # OSX
```

### Releases

Releases for this project can be found here under github's releases integration.

https://github.com/deps-cloud/gitfs/releases

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
clone:
  repositoryRoot: "${HOME}/.gitfs/cache" # this is where your repos will be cloned
  depth: 0 # depth 0 means full clone
```

### Running gitfs

Once your configuration has been created and values have been replaced, you can run `gitfs start`.
Note that running `gitfs` will spawn the filesystem service as a blocking process in your active terminal window.
You can run this as a systemctl daemon, but I find having access to the logs more readily at this time is a bit nicer.

To run, simply execute `gitfs start`.
When run successfully, you should see an output message similar to the one below.

```
$ gitfs start
INFO[0000] [main] configured mount point: /Users/deps-cloud/Development/code 
INFO[0000] [main] fetching repositories                 
INFO[0000] [remotes.github] processing organizations for user: deps-cloud 
INFO[0000] [remotes.github] processing repositories for user: deps-cloud 
INFO[0001] [main] parsing repositories into a directory structure 
INFO[0001] [main] attempting to mount /Users/deps-cloud/Development/code 
INFO[0001] [main] now serving file system  
``` 

In a separate window, navigate to your configured `mount` point.
For me, I have this set `${HOME}/Development/code`.
From this directory, you should be able to navigate your projects on Github.

```
[deps-cloud@deps-cloud ~/Development/code]
$ ls
github.com

[deps-cloud@deps-cloud ~/Development/code]
$ cd github.com/

[deps-cloud@deps-cloud ~/Development/code/github.com]
$ ls
deps-cloud

[deps-cloud@deps-cloud ~/Development/code/github.com]
$ cd deps-cloud/

[deps-cloud@deps-cloud ~/Development/code/github.com/deps-cloud]
$ ls
OpenGrok	   grpc-java	    laas	rpi
consul-api	   grpc.github.io   deps-cloud.com	seo-portal
docker-clickhouse  grpcsh	    mp		serverless-plugin-simulate
docker-utils	   hbase-docker     okhttp	simple-daemon-node
dotfiles	   idea-framework   proctor	spring-config-repo
generator-idea	   java-gitlab-api  proctorjs
gitfs		   jgrapht	    proto2-3

[deps-cloud@deps-cloud ~/Development/code/github.com/deps-cloud]
$ cd grpcsh

[deps-cloud@deps-cloud ~/Development/code/github.com/deps-cloud/grpcsh master]
$ ls
README.md  bin	package-lock.json  package.json  yarn.lock
```

### Stopping gitfs

In order to stop the gitfs server, you will need to run the `stop` command.
This forces the FUSE volume to be unmounted.

```bash
$ gitfs stop
INFO[0000] Unmounting /Users/deps-cloud/Development/code     
INFO[0000] Successfully unmounted /Users/deps-cloud/Development/code 
```

Should you run this command and find that the file system did not unmount, you can try using `umount`.

```bash
# replace with your mount point
$ sudo umount -f <your_mount_point>
```
