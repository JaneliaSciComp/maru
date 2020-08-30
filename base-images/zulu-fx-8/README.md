# zulu-fx-8 Docker image

This image bundles the [ZuluFX build](https://www.azul.com/downloads/zulu-community/?version=java-8-lts&os=ubuntu&architecture=x86-64-bit&package=jdk-fx) which includes OpenJDK and OpenJFX. JavaFX programs can be run inside the container and displayed in X11. This base image can be used to create Docker containers for JavaFX apps.

## Linux Docker

The container should just work like any X11 application.

## Docker for Mac

To forward X11 from inside Docker container to a Mac host:

1. Install the [XQuartz](https://www.xquartz.org) X11 server
2. Launch XQuartz and open Preferences
3. In the Security tab, enable "Allow connections from network clients"
4. Restart XQuartz
5. In a Terminal, run `xhost + localhost` (this step must be rerun every time XQuartz is restarted)
6. Run the container with the `-e DISPLAY=host.docker.internal:0` option

