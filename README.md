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

First, create a new directory for your Maru project configuration and tell Maru how to find and build your code. You can name the directory anything you like:
```
mkdir myproject ; cd myproject
maru init
```

Now build your project into a Docker container:
```
maru build
```

You can run the current version of your Docker container:
```
maru run [args to app]
```
This will output the Docker command used to run the container, which you can then use in scripts or pipelines to integrate your app into a larger workflow.

You can also run the Docker container using Singularity (e.g. on an HPC cluster):
```
maru singularity run [args to app]
```

## Documentation

For more details, see the [docs/UserManual.md](User Manual).

For developers, there are some notes available about [docs/Development.md](building and releasing Maru)

