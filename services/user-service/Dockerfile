# services/user-service/Dockerfile

# ビルドステージ
FROM golang:1.21-alpine AS builder

# 必要なビルドツールのインストール
RUN apk add --no-cache git

WORKDIR /app

# モジュール依存関係のコピーと解決
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# アプリケーションのビルド
# CGO_ENABLED=0: 静的リンクを有効にする
# -ldflags="-w -s": バイナリサイズの削減
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o userservice ./cmd/main.go

# 実行ステージ
FROM alpine:3.19

# セキュリティ更新とCA証明書のインストール
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# 非root ユーザーの作成
RUN adduser -D -u 1000 appuser

# 作業ディレクトリの設定
WORKDIR /app

# ビルドしたバイナリのコピー
COPY --from=builder /app/userservice .
COPY --from=builder /app/.env.example .env

# 実行ユーザーの変更
USER appuser

EXPOSE 8080

# ヘルスチェックの設定
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# アプリケーションの実行
CMD ["./userservice"]
