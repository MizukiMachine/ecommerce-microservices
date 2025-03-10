// services/user-service/cmd/main.go
package main

import (
	"log"
	"os"
	"time"

	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/infrastructure/auth"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/infrastructure/database"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/infrastructure/middleware"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/infrastructure/persistence"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/interface/handler"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// 2. データベース接続の設定
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "user_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// 3. データベース接続
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 4. リポジトリの初期化
	userRepo := persistence.NewUserRepository(db)

	// 5. JWTサービスの初期化
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpiration, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	jwtService := auth.NewJWTService(jwtSecret, jwtExpiration)

	// 6. ユースケースの初期化
	userUseCase := usecase.NewUserUseCase(userRepo, jwtService)

	// 7. ハンドラーの初期化
	userHandler := handler.NewUserHandler(userUseCase)

	// 8. Ginルーターの設定
	router := gin.Default()

	// 9. 認証ミドルウェアの初期化
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// 10. 基本ミドルウェアの設定
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// ルーティングの設定
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			// 認証不要のエンドポイント
			users.POST("/register", userHandler.CreateUser)
			users.POST("/login", userHandler.Login)

			// 認証が必要なエンドポイント
			auth := users.Use(authMiddleware.AuthRequired())
			{
				auth.GET("/profile", userHandler.GetProfile)
				auth.PUT("/profile", userHandler.UpdateProfile)
				auth.POST("/refresh-token", authMiddleware.RefreshToken())
			}
		}
	}

	// 10. サーバーの起動
	port := getEnv("PORT", "8080")
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 環境変数を取得する関数
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 認証ミドルウェア
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT認証の実装
		c.Next()
	}
}
