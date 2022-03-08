# Step 1 - compile code binary
FROM golang:1.17.8-alpine AS builder

LABEL maintainer="Charles Guertin <charlesguertin@live.ca>"

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT}

RUN apk add --no-cache --update ca-certificates make

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN make build

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
