# fiji-openjdk-8 Docker image

This is an unofficial Fiji image including OpenJDK 8, the FFMPEG library, and the latest Fiji at the time of the build. It was derived from the [official image Dockerfile](https://github.com/fiji/dockerfiles). 

## Why not just use the official Fiji images?

We are currently maintaining these images because:
* The official images are outdated and do not have the latest plugins we need.
* Certain plugins (e.g. H5J Loader) do not work in the official images due to lack of shared libraries.
* We need Docker containers that we can run with either Docker or Singularity, and with Nextflow.

These images should be considered **under development**. In the future, we may try to merge our changes back to the official repo.

