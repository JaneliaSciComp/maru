# Maru

Maru is an opinionated command-line interface for quickly and easily containerizing scientific applications, enabling reproducible results and low-friction collaboration. 

Maru provides a CLI wizard for easily creating containers of various flavors (Python, Java, MATLAB, Fiji, etc.) which follow best practices and are optimized for performance. It also makes it easy to keep your container versioned and provides convenience commands so that you can focus on your algorithms instead of learning arcane details about containers. 

## Get Maru

Maru depends on [Docker](https://docs.docker.com/get-docker/) and/or [Singularity](https://github.com/hpcng/singularity).

The following one-liners install the `maru` binary into /usr/local/bin:

### Linux
```
sudo curl -sL https://data.janelia.org/maru-linux | tar -xz -C /usr/local/bin
```

### MacOS
```
curl -sL https://data.janelia.org/maru-macos | tar -xz -C /usr/local/bin
```

You can also download the [latest release](https://github.com/JaneliaSciComp/maru/releases/latest) and manually copy it to anywhere in your `$PATH`.

## Usage

Maru assumes that your project is committed to a remote git repository. During the container build, Maru will clone your repository and run any build commands you specify.

First, initialize the project configuration:
```
maru init
```
This is an interactive process that will ask you questions about where to find your project and how to build it.

Now you can build your project into a Docker container:
```
maru build
```

From your project directory, you can always run the current version of your Docker container:
```
maru run [args to app]
```
This will also output the Docker command used to run the container, which you can then use in scripts or pipelines to integrate your app into a larger workflow.

You can also run the Docker container using Singularity (e.g. on an HPC cluster):
```
maru singularity run [args to app]
```

## Documentation

Maru has lots of features for building, releasing, and distributing your containers. For more details, see the [User Manual](docs/UserManual.md).

For developers, there are some notes available about [building and releasing Maru](docs/Development.md).

