# ---- Build Stage ----
FROM golang:1.24 AS builder

ARG VERSION=development

WORKDIR /app

# Copy the code
COPY . .

# Install dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -X 'main.Version=${VERSION}'" -o /hyaline ./cmd/hyaline.go

# ---- Final Stage ----
FROM scratch

LABEL org.opencontainers.image.source=https://github.com/appgardenstudios/hyaline

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /hyaline /hyaline

ENTRYPOINT ["/hyaline"]

CMD ["--help"]