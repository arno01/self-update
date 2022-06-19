FROM alpine:3.16

RUN apk update && \
    apk add jq procps coreutils bash ca-certificates runit

# custom init handles signals such as SIGTERM and propagates them down to
# services running inside this container.
COPY docker/my_init.sh /sbin/my_init

# runit starts your apps when detects `/etc/service/<your-app>/run`
# as well as it starts auto-updating script.
COPY docker/service /etc/service

# app gets built by the pipeline (github actions)
COPY app /usr/local/bin/app

# VERSION gets passed by the pipeline
# it is used by the auto-updating script to compare against the upstream version
ARG VERSION=dev
RUN echo $VERSION | tee /version
ENTRYPOINT ["/sbin/my_init"]
