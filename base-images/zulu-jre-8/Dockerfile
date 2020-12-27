# Download Zulu JRE
FROM janeliascicomp/common-tools:1.0.2 as builder
WORKDIR /opt/zulu-8-jre-amd64
RUN curl -sL https://cdn.azul.com/zulu/bin/zulu8.50.0.51-ca-jre8.0.275-linux_x64.tar.gz \
    | tar -xz --strip-components=1 -C /opt/zulu-8-jre-amd64

# Create final image
FROM ubuntu:18.04

# eliminate input warnings from debconf
ARG DEBIAN_FRONTEND=noninteractive

# JDK paths
ENV JAVA_HOME=/usr/lib/jvm/zulu-8-jre-amd64
ENV PATH="$JAVA_HOME/bin:${PATH}"

RUN apt-get update \
    # Set UTF-8 by default, this is needed to run with Singularity
    && apt-get install --no-install-recommends -y locales \
    && echo "LC_ALL=en_US.UTF-8" >> /etc/environment \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "LANG=en_US.UTF-8" > /etc/locale.conf \
    && locale-gen en_US.UTF-8 \
    # Clean up this layer
    && rm -rf /var/lib/apt/lists/* 

# Copy in the JRE
COPY --from=builder /opt/zulu-8-jre-amd64 /usr/lib/jvm/zulu-8-jre-amd64

