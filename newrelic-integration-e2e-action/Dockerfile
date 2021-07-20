ARG module=validator
############################
# STEP 1 build binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src
COPY . .
RUN go build -o /go/bin/${module} cmd/${module}/main.go
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/${module} /${module}
# Run the binary.
ENTRYPOINT ["/${module}"]