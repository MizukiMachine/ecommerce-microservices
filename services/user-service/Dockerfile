# ビルドステージ
FROM golang:1.21-alpine AS builder

# 必要なビルドツールのインストール
RUN apk add --no-cache git

# 作業ディレクトリの設定
WORKDIR /app

# モジュール依存関係のコピーと解決
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .
