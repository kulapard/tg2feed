FROM golang:1.21 as builder

ARG REVISION

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

# Statically compile our app for use in a distroless container
RUN cd app && go build -o /build/tg2rss -ldflags "-X main.revision=${REVISION} -s -w"

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static

COPY --from=builder /build/tg2rss /app

ENTRYPOINT ["/app"]