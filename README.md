# Maru

Maru is a command-line interface for quickly and easily containerizing scientific applications. 

## Prerequisites

Maru runs on both Linux and MacOS systems. You need to have [Docker installed](https://docs.docker.com/get-docker/) to use Maru.

## Installation

### Linux
This installs the `maru` binary into /usr/local/bin:
```
curl -sL https://github.com/JaneliaSciComp/maru/releases/latest/download/maru_linux_x86_64.tar.gz | tar -xz -C /usr/local/bin
```

### MacOS
This installs the `maru` binary into /usr/local/bin:
```
sudo curl -sL https://github.com/JaneliaSciComp/maru/releases/latest/download/maru_macos_x86_64.tar.gz | tar -xz -C /usr/local/bin
```

You can also download the [latest release](https://github.com/JaneliaSciComp/maru/releases/latest) and copy it to anywhere in your `$PATH`.

## Usage

Maru assumes that your project is available in a git repository, and it checks out and builds your code while 
building the container. 

To initialize a new Maru project in the current directory:
```
maru init
```

Build the Docker container for the Maru project in the current directory:
```
maru build
```

Run the Docker container for the Maru project in the current directory:
```
maru run [args to containerized program]
```

Change the git tag that will be used to during the next `maru build`:
```
maru set repo_tag <new tag>
```

Change the version tag that will be used to tag your built container:
```
maru set version <new version>
```

## Building

To compile and install Maru into your standard Go bin directory:
```
go install
```

Any time the templates change, the serialization needs to be updated as follows:
```
go generate ./...
```

## Testing

TBD

## Releasing

New releases are built and deployed to GitHub using GoReleaser. You first need to install GoReleaser and configure your GitHub token as per the [quickstart instructions](https://goreleaser.com/quick-start/). Then you can tag and release a new version as follows:

```
git tag -a 0.1.0 -m "Release 0.1.0"
git push origin 0.1.0
goreleaser --rm-dist
```

To test a SNAPSHOT release without tagging the code:
```
goreleaser --snapshot --skip-publish --rm-dist
```

