# Builder image
FROM --platform=linux/amd64 golang:1.24 as builder

WORKDIR /go/src/go-tooling-portal

COPY . /go/src/go-tooling-portal

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

RUN go mod download && \
    go build -o go-tooling-portal -ldflags "-X 'github.com/sfarosu/go-tooling-portal/internal/version.BuildDate=$(date '+%Y-%m-%d %H:%M:%S')'-X 'github.com/sfarosu/go-tooling-portal/internal/version.GitShortHash=$(git rev-parse --short HEAD)'"

# Run image
FROM --platform=linux/amd64 alpine:3.21

RUN apk update && apk add --no-cache --upgrade openssh-client openssl tzdata

WORKDIR /app

COPY ./web/ ./web
COPY ./scripts/ ./scripts
COPY --from=builder /go/src/go-tooling-portal/go-tooling-portal .

RUN chown -R 10001:root /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app /etc/passwd && \
    chmod -R a+x-w /app/scripts && \
    chmod a+x-w /app/go-tooling-portal

EXPOSE 8080

ENTRYPOINT [ "./scripts/uid_entrypoint.sh" ]

USER 10001

CMD [ "./go-tooling-portal" ]
