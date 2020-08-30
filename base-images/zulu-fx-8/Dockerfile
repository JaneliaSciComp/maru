FROM ubuntu:18.04

# eliminate input warnings from debconf
ARG DEBIAN_FRONTEND=noninteractive

# set DISPLAY to host.docker.internal:0 for running in Docker for Mac
ENV DISPLAY=:0 
# set timezone for GTK
ENV TZ=US/Eastern
# JDK paths
ENV JAVA_HOME=/usr/lib/jvm/zulu-8-fx-amd64
ENV PATH="$JAVA_HOME/bin:${PATH}"

RUN apt-get update \
    # UTF-8 by default
    && apt-get install --no-install-recommends -y locales curl \
    && echo "LC_ALL=en_US.UTF-8" >> /etc/environment \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "LANG=en_US.UTF-8" > /etc/locale.conf \
    && locale-gen en_US.UTF-8 \
    # Dependencies for X11 rendering
    && apt-get install --no-install-recommends -y libxext6 libxrender1 libxtst6 libxi6 \
    # Because https://support.azul.com/hc/en-us/articles/360034030692-Using-Fonts-with-OpenJDK-Zulu
    && apt-get install --no-install-recommends -y fonts-dejavu fontconfig \
    # GTK required by JavaFX
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && apt-get install --no-install-recommends -y gtk3.0 \
    # Install ZuluFX
    && mkdir -p /usr/lib/jvm/zulu-8-fx-amd64 \
	&& curl -sL https://cdn.azul.com/zulu/bin/zulu8.48.0.53-ca-fx-jdk8.0.265-linux_x64.tar.gz \
    | tar -xz --strip-components=1 -C /usr/lib/jvm/zulu-8-fx-amd64 \
    # Clean up this layer
    && apt-get purge -y curl \
    && rm -rf /var/lib/apt/lists/* 

# Add demo
WORKDIR /app
ADD HelloWorld.java .
RUN javac HelloWorld.java
# prime.order=sw forces software rendering so that this will run on Mac without OpenGL
CMD [ "java", "-Dprism.order=sw", "-cp", ".", "HelloWorld" ]

