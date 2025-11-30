# Multi-stage build for backend API
# go.mod / toolchain に合わせて 1.24 系を使用
FROM golang:1.24.5-bullseye AS builder

WORKDIR /app

# 依存取得用に mod/sum だけ先にコピー
COPY go.mod go.sum ./
RUN go mod download

# 残りのソースをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/server ./cmd/api

# ---- runtime stage ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/server /app/server
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/server"]
