############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src
COPY . .
# Build the binary.
RUN go build -o /go/bin/spec-validator cmd/main.go
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/spec-validator /spec-validator
# Run the e2e-spec-validator binary.
ENTRYPOINT ["/spec-validator"]