# Microservices-based E-commerce Platform

eコマースプラットフォームの構築

### マイクロサービスアーキテクチャ
- サービス間通信: gRPC & REST
- イベント駆動アーキテクチャ

### ユーザーサービス（Go）
- JWT認証による安全なユーザー管理
- RBACによるきめ細かな権限制御
- OAuth2.0による外部サービス連携

### 言語とフレームワーク
- Go: ユーザー管理・商品管理サービス
  - gin-gonic/gin
  - gorm
- Rust: 注文・在庫管理サービス
  - actix-web
  - tokio
