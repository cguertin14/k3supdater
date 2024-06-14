# Step 1 - compile code binary
FROM golang:1.22.4-alpine AS builder

LABEL maintainer="Charles Guertin <charlesguertin@live.ca>"

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT=""
ARG BUILD_DATE
ARG VERSION
ARG GIT_COMMIT

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=${TARGETVARIANT} \
    BUILD_DATE=${BUILD_DATE} \
    VERSION=${VERSION} \
    GIT_COMMIT=${GIT_COMMIT}

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build \
    -ldflags "-X github.com/cguertin14/k3supdater/cmd.Version=${VERSION} -X github.com/cguertin14/k3supdater/cmd.BuildDate=${BUILD_DATE} -X github.com/cguertin14/k3supdater/cmd.GitCommit=${GIT_COMMIT}" \
    -o ./k3supdater .

# Step 2 - import necessary files to run program.
FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=builder /app/k3supdater /k3supdater
ENTRYPOINT ["/k3supdater"]
