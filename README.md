# Maru

Maru is a command-line interface for quickly and easily containerizing scientific applications. 

## Install

TBD

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
go generate
```

## Testing

TBD

