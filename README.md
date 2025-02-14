# Microservices-based E-commerce Platform

eコマースプラットフォーム

- サービス間通信: gRPC & REST

### ユーザーサービス（Go）
- JWT認証
- RBAC
- OAuth2.0
  
### 言語とフレームワーク
- Go: ユーザー管理・商品管理サービス
  - gin-gonic/gin
  - gorm
- Rust: 注文・在庫管理サービス
  - actix-web
  - tokio
