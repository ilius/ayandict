FROM golang:1.24-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" \
    -o /out/ayandict-web ./cmd/ayandict-web

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=build /out/ayandict-web /usr/local/bin/ayandict-web

EXPOSE 8357

ENTRYPOINT ["/usr/local/bin/ayandict-web"]
CMD ["--expose"]
