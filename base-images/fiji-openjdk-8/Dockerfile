#
# Fiji + Java 8 Dockerfile
#

# Pull base JDK-8 image.
FROM openjdk:8

# Define maintainer.
LABEL maintainer="Janelia Scientific Computing"

# Install FFMPEG dependencies for H5J
RUN apt-get update && apt-get install -y libxcb1 libxcb-shm0 libxcb-shape0 libxcb-xfixes0 libasound2

# Install the latest Fiji
WORKDIR /opt/fiji
RUN curl -sL https://downloads.imagej.net/fiji/latest/fiji-nojre.zip | jar xv \
    && chmod +x /opt/fiji/Fiji.app/ImageJ-linux64

# Add fiji to the PATH
ENV PATH $PATH:/opt/fiji/Fiji.app

# Define entrypoint.
COPY entrypoint.sh /opt/fiji
ENTRYPOINT ["/opt/fiji/entrypoint.sh"]

# Run once to create Java preferences.
COPY demo.py /opt/fiji/
RUN /opt/fiji/entrypoint.sh --headless --ij2 --console --run /opt/fiji/demo.py 'name="test"'

