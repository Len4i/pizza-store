# Build the binary
FROM golang:1.22.2 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/


# Test
RUN CGO_ENABLED=0  GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o pizza_store cmd/main.go

# Use distroless as minimal base image to package the  bnary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/pizza_store .
USER 65532:65532

ENTRYPOINT ["/pizza_store"]