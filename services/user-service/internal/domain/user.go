package domain

import (
	"context"
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password does not meet security requirements")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

// User エンティティ
type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ドメインのビジネスルール
func (u *User) Validate() error {
	// メールアドレスの検証
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}

	// パスワードの検証
	if !isStrongPassword(u.Password) {
		return ErrWeakPassword
	}

	return nil
}

// メールアドレスのバリデーション
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// パスワード強度のチェック
func isStrongPassword(password string) bool {
	// 最小8文字、大文字小文字数字を含む
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// UserRepository インターフェース
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
