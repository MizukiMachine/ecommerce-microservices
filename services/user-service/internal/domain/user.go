package domain

import (
	"context"
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidEmail = errors.New("invalid email format")
	ErrWeakPassword = errors.New("password does not meet security requirements")
)

// ユーザー情報が正しいかチェックする関数
func (u *User) Validate() error {
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}

	if !isStrongPassword(u.Password) {
		return ErrWeakPassword
	}

	return nil
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}

// ユーザー情報を扱うための操作を定義
type UserRepository interface {
	// 新しいユーザーを作成
	Create(ctx context.Context, user *User) error
	// IDからユーザーを検索
	FindByID(ctx context.Context, id string) (*User, error)
	// メールアドレスからユーザーを検索
	FindByEmail(ctx context.Context, email string) (*User, error)
	// ユーザー情報を更新
	Update(ctx context.Context, user *User) error
	// ユーザーを削除
	Delete(ctx context.Context, id string) error
}
