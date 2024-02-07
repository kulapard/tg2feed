FROM golang:1.22 as builder

ARG REVISION

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

# Statically compile our app for use in a distroless container
RUN cd app && go build -o /build/tg2feed -ldflags "-X main.revision=${REVISION} -s -w"

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static

COPY --from=builder /build/tg2feed /app

ENTRYPOINT ["/app"]