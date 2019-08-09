FROM golang:1.12-stretch 

# Multi-stage builds are not supported on docker version 1.12/1.13 - thus in openshift 3.X

WORKDIR /go/src/tooling-portal

COPY *.go /go/src/tooling-portal/

RUN CGO_ENABLED=0 go build -a -installsuffix nocgo -o /go/bin/tooling-portal . && \
    mkdir /app && \
    cp /go/bin/tooling-portal /app/tooling-portal

WORKDIR /app

COPY . /app/

RUN rm -Rf /app/*.go /app/Dockerfile /app/.git  && \
    rm -Rf /go/src/tooling-portal && \
    useradd -d /app -s /bin/bash tooling && \
    chown -R tooling:root /app && \
    chmod -R u+x /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app /etc/passwd /etc/group
    
EXPOSE 8080

USER tooling

CMD ["./tooling-portal"]




