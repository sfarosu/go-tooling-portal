# Builder image
FROM golang:1.27 AS builder

WORKDIR /go/src/go-tooling-portal

COPY . /go/src/go-tooling-portal

ENV GO111MODULE=on \
    CGO_ENABLED=0

RUN go mod download && \
    go build -o go-tooling-portal -ldflags "-X 'github.com/sfarosu/go-tooling-portal/internal/version.BuildDate=$(date '+%Y-%m-%d %H:%M:%S')'-X 'github.com/sfarosu/go-tooling-portal/internal/version.GitShortHash=$(git rev-parse --short HEAD)'"

# Run image
FROM alpine:3.24

RUN apk add --no-cache tzdata

WORKDIR /app

COPY ./web/ ./web
COPY ./scripts/ ./scripts
COPY --from=builder /go/src/go-tooling-portal/go-tooling-portal .

RUN chown -R 10001:root /app && \
    chgrp -R 0 /app && \
    chmod -R g=u /app /etc/passwd && \
    chmod -R a+x-w /app/scripts && \
    chmod a+x-w /app/go-tooling-portal

EXPOSE 8100

ENTRYPOINT [ "./scripts/uid_entrypoint.sh" ]

USER 10001

CMD [ "./go-tooling-portal" ]
