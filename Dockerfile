FROM golang AS builder
WORKDIR /build
COPY go.mod go.sum /build/

COPY . /build
WORKDIR /build/cmd/ssh-check
ENV CGO_ENABLED=0
RUN go build -v

FROM scratch
COPY --from=builder /build/cmd/ssh-check/ssh-check /app/ssh-check
ENTRYPOINT ["/app/ssh-check"]