FROM golang:1.13-stretch
#FROM i386/golang:1.13-stretch

# Multi-stage builds are not supported on docker version 1.12/1.13 - thus in openshift 3.X

WORKDIR /go/src/go-tooling-portal

COPY . /go/src/go-tooling-portal/

RUN CGO_ENABLED=0 go build -a -installsuffix nocgo -o /go/bin/go-tooling-portal ./cmd/go-tooling-portal && \
    rm -Rf /go/src/go-tooling-portal && \
    mkdir /app && \
    mv /go/bin/go-tooling-portal /app/go-tooling-portal

WORKDIR /app

COPY . /app/

RUN rm -Rf /app/.git* /app/cmd /app/deployments /app/go.* && \
    chown -R 10001:root /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app /etc/passwd && \
    chmod a+x-w /app/build/uid_entrypoint.sh
    
EXPOSE 8080

ENTRYPOINT [ "/app/build/uid_entrypoint.sh" ]

USER 10001

CMD [ "./go-tooling-portal" ]

