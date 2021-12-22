# Builder image
FROM golang:1.17 as builder

WORKDIR /go/src/go-tooling-portal

COPY . /go/src/go-tooling-portal

RUN go mod download && \
    CGO_ENABLED=0 go build -o go-tooling-portal

# Run image
FROM alpine:3.13

RUN apk update && apk add --no-cache --upgrade openssh-client openssl

WORKDIR /app

COPY ./web/ ./web
COPY ./build/ ./build
COPY --from=builder /go/src/go-tooling-portal/go-tooling-portal .

RUN chown -R 10001:root /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app /etc/passwd && \
    chmod -R a+x-w /app/build && \
    chmod a+x-w /app/go-tooling-portal

EXPOSE 8080

ENTRYPOINT [ "./build/uid_entrypoint.sh" ]

USER 10001

CMD [ "./go-tooling-portal" ]
