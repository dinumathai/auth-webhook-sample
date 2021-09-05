FROM golang:latest as builder

ENV RELEASE_VERSION=0.0.1

WORKDIR /go/src/github.com/dinumathai/auth-webhook-sample

COPY . /go/src/github.com/dinumathai/auth-webhook-sample

# tests with coverage
#RUN go test -v github.com/dinumathai/auth-webhook-sample/...

RUN CGO_ENABLED=0 go install -ldflags="-X main.Version=${RELEASE_VERSION}" github.com/dinumathai/auth-webhook-sample

# Create New Base Image.
FROM alpine:latest

WORKDIR /go/bin/
RUN mkdir config

COPY --from=builder /go/bin/auth-webhook-sample .
COPY --from=builder /go/src/github.com/dinumathai/auth-webhook-sample/swaggerui swaggerui/
COPY --from=builder /go/src/github.com/dinumathai/auth-webhook-sample/config/*.yaml config/

EXPOSE 8443
CMD ["./auth-webhook-sample"]
