# Builder stage
FROM golang AS builder

ENV GO111MODULE=on

ARG version
ARG build

ENV VERSION=$version
ENV BUILD=$build
ENV OUTBINDIR=/app

WORKDIR /app

# Enable Docker cache for go modules.
COPY go.mod .
COPY go.sum .
RUN GOPATH=/app/.go-pkg:/app GOBIN=/app/.go-bin go mod download

COPY Makefile .
RUN make install

COPY . .

RUN make compile-linux

# Runner stage
FROM scratch
# If we need the config.toml in the future that'll need to be copied over too.
COPY --from=builder /app/hbaas-server /app/

# Document that the service listens on port 8080.
EXPOSE 8000
EXPOSE 4443

# Run the outyet command by default when the container starts.
CMD ["/app/hbaas-server", "run-server"]
