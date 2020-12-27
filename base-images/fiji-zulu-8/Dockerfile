#
# Fiji + Zulu Java 8 Dockerfile
#
FROM janeliascicomp/common-tools:1.0.2 as fiji-builder
ARG FIJI_VERSION=20201104-1356
WORKDIR /app/fiji
RUN curl -sL https://downloads.imagej.net/fiji/archive/$FIJI_VERSION/fiji-nojre.tar.gz | tar -xz -C /app/fiji
RUN chmod +x /app/fiji/Fiji.app/ImageJ-linux64 \
    && chown -R root:root /app/fiji \
    && ln -s /app/fiji/Fiji.app/ImageJ-linux64 /app/fiji/fiji

FROM janeliascicomp/zulu-jre:8u275b01

# Define maintainer
LABEL maintainer="Janelia Scientific Computing"

# Add fiji to the PATH
ENV PATH $PATH:/app/fiji/Fiji.app

# Install an archived snapshot of Fiji
RUN apt-get update \
    && apt-get install --no-install-recommends -y \
        libxcb1 libxcb-shm0 libxcb-shape0 libxcb-xfixes0 libasound2 
COPY --from=fiji-builder /app/fiji /app/fiji

# Define entrypoint
COPY entrypoint.sh /app/fiji/
ENTRYPOINT ["/app/fiji/entrypoint.sh"]

# Run once to create Java preferences
COPY demo.py /app/fiji/
RUN /app/fiji/entrypoint.sh --headless --ij2 --console --run /app/fiji/demo.py 'name="test"'

