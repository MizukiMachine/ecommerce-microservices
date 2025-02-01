package usecase

import (
	"context"
	"time"

	domain "github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

// ユースケースの入力データ
type CreateUserInput struct {
	Email    string
	Password string
	Name     string
}

// ユースケースの出力データ
type UserOutput struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

// ユースケース構造体
type UserUseCase struct {
	userRepo domain.UserRepository
}

// ユースケースの作成
func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

// ユーザー作成のユースケース
func (uc *UserUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	// 1. パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. ドメインオブジェクトの作成
	user := &domain.User{
		Email:     input.Email,
		Password:  string(hashedPassword),
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 3. ドメインのバリデーション
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 4. メールアドレスの重複チェック
	existingUser, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// 5. ユーザーの保存
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 6. 出力データの作成
	return &UserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

// ユーザー認証のユースケース
// ユーザー情報取得
func (uc *UserUseCase) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return &UserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

// プロフィール更新
func (uc *UserUseCase) UpdateUserProfile(ctx context.Context, userID string, name string) (*UserOutput, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	user.Name = name
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &UserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (uc *UserUseCase) AuthenticateUser(ctx context.Context, email, password string) (*UserOutput, error) {
	// 1. メールアドレスでユーザーを検索
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// 2. パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// 3. 出力データの作成
	return &UserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
