// services/user-service/internal/interface/handler/user_handler.go
package handler

import (
	"net/http"
	"time"

	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/domain"
	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/usecase"
	"github.com/gin-gonic/gin"
)

// リクエストの形式を定義
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// ログインリクエストの形式を定義
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// レスポンスの形式を定義
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// エラーレスポンスの形式
type ErrorResponse struct {
	Message string `json:"message"`
}

// ハンドラー構造体
type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

// ハンドラーの作成
func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: uc,
	}
}

// ユーザー作成のハンドラー
func (h *UserHandler) CreateUser(c *gin.Context) {
	// 1. リクエストのバリデーション
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request format",
		})
		return
	}

	// 2. ユースケースの入力データを作成
	input := usecase.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// 3. ユースケースの実行
	output, err := h.userUseCase.CreateUser(c.Request.Context(), input)
	if err != nil {
		// エラーの種類に応じて適切なステータスコードを返す
		switch err {
		case domain.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Message: "Email already exists",
			})
		case domain.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Invalid email format",
			})
		case domain.ErrWeakPassword:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: "Password does not meet security requirements",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: "Internal server error",
			})
		}
		return
	}

	// 4. レスポンスの作成と返却
	response := UserResponse{
		ID:        output.ID,
		Email:     output.Email,
		Name:      output.Name,
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}
	c.JSON(http.StatusCreated, response)
}

// プロフィール取得ハンドラー
func (h *UserHandler) GetProfile(c *gin.Context) {
	// コンテキストから認証済みユーザーのIDを取得
	userID := c.GetString("userID") // authMiddlewareでセットされることを想定
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	// ユーザー情報の取得
	user, err := h.userUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to get user profile",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// プロフィール更新ハンドラー
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request format",
		})
		return
	}

	// ユーザー情報の更新
	user, err := h.userUseCase.UpdateUserProfile(c.Request.Context(), userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// ログインハンドラー
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Invalid request format",
		})
		return
	}

	output, err := h.userUseCase.AuthenticateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Internal server error"

		if err == domain.ErrInvalidCredentials {
			status = http.StatusUnauthorized
			message = "Invalid credentials"
		}

		c.JSON(status, ErrorResponse{
			Message: message,
		})
		return
	}

	// JWTトークンの生成（詳細は後ほど実装）
	token := "dummy-jwt-token"

	response := UserResponse{
		ID:        output.ID,
		Email:     output.Email,
		Name:      output.Name,
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  response,
	})
}
