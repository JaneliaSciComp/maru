# Staged build using builder container
FROM janeliascicomp/builder:1.2.0 as builder
ARG APP_TAG=master

# Checkout and build the code
WORKDIR /tmp/app
RUN git clone --branch $APP_TAG --depth 1 https://github.com/PreibischLab/BigStitcher.git . \
    && mvn -Pfatjar package

# Find the built jar, based on the version in the pom file
RUN xq -r '.project.artifactId+"-"+.project.version+".jar"' pom.xml > filename \
    && mv /tmp/app/target/`cat filename` app.jar

# Create final image
FROM openjdk:8

COPY --from=builder /tmp/app/app.jar /app/app.jar

RUN echo "#!/bin/bash" >> /entrypoint.sh \
    && echo 'java -cp /app/app.jar net.preibisch.mvrecon.fiji.plugin.resave.Resave_N5 "$@"' >> /entrypoint.sh \
    && chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh"]
