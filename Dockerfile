# Step 1 - compile code binary
FROM golang:1.17.8-alpine AS builder

LABEL maintainer="Charles Guertin <charlesguertin@live.ca>"

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT} \
    BUILD_DATE=${BUILD_DATE} \
    VERSION=${VERSION} \
    GIT_COMMIT=${GIT_COMMIT}

RUN apk add --no-cache --update ca-certificates make

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o ./k3supdater . \
    -X github.com/cguertin14/k3supdater/cmd.BuildDate=${BUILD_DATE} \
    -X github.com/cguertin14/k3supdater/cmd.GitCommit=${GIT_COMMIT} \
    -X github.com/cguertin14/k3supdater/cmd.Version=${VERSION}

# Add user & group
RUN addgroup -S updater-group && \
    adduser -S updater-user -G updater-group


# Step 2 - import necessary files to run program.
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/k3supdater /k3supdater
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER updater-user

ENTRYPOINT ["/k3supdater"]
