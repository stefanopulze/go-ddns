# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -a -ldflags '-w -s' -o /dist/go-ddns ./cmd/server

FROM gcr.io/distroless/static-debian13:nonroot
ENV TZ=Europe/Rome

WORKDIR /app

COPY --from=builder ./dist/go-ddns /app/go-ddns

EXPOSE 8080

ENTRYPOINT ["/app/go-ddns"]