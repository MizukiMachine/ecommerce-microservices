package persistence

import (
	"context"
	"time"

	"github.com/MizukiMachine/ecommerce-microservices/services/user-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// データベースのテーブル構造
type UserModel struct {
	ID        string `gorm:"primaryKey;type:uuid"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// リポジトリの構造体
type userRepository struct {
	db *gorm.DB
}

// リポジトリを作成する関数
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	// テーブルの自動作成（本番環境では手動マイグレーションを推奨）
	db.AutoMigrate(&UserModel{})

	return &userRepository{
		db: db,
	}
}

// ドメインモデルをDBモデルに変換
func toModel(user *domain.User) *UserModel {
	return &UserModel{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// DBモデルをドメインモデルに変換
func toDomain(model *UserModel) *domain.User {
	return &domain.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Name:      model.Name,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// ユーザーの作成
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	// UUIDの生成
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	model := toModel(user)
	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		// メールアドレスの重複エラーをドメインエラーに変換
		if isDuplicateKeyError(result.Error) {
			return domain.ErrEmailAlreadyExists
		}
		return result.Error
	}

	// 生成されたIDを元のユーザーオブジェクトに反映
	user.ID = model.ID
	return nil
}

// メールアドレスでユーザーを検索
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toDomain(&model), nil
}

// IDでユーザーを検索
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel
	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return toDomain(&model), nil
}

// ユーザー情報の更新
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	model := toModel(user)
	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return domain.ErrEmailAlreadyExists
		}
		return result.Error
	}
	return nil
}

// ユーザーの削除
func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	return result.Error
}

// データベースのエラーを判定するヘルパー関数
func isDuplicateKeyError(err error) bool {
	// PostgreSQLの一意性制約違反のエラーコードをチェック
	// 実際のコードはデータベースに応じて適切に実装する必要があります
	return err.Error() == "ERROR: duplicate key value violates unique constraint"
}
