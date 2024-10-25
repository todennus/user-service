package gorm

import (
	"context"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/infras/database/model"
	"github.com/todennus/x/enum"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := model.NewUser(user)
	return errordef.ConvertGormError(repo.db.WithContext(ctx).Create(&model).Error)
}

func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	model := model.UserModel{}
	if err := repo.db.WithContext(ctx).Take(&model, "username=?", username).Error; err != nil {
		return nil, errordef.ConvertGormError(err)
	}

	return model.To()
}

func (repo *UserRepository) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	model := model.UserModel{}
	if err := repo.db.WithContext(ctx).Take(&model, "id=?", userID).Error; err != nil {
		return nil, errordef.ConvertGormError(err)
	}

	return model.To()
}

func (repo *UserRepository) CountByRole(ctx context.Context, role enum.Enum[enumdef.UserRole]) (int64, error) {
	var n int64
	err := repo.db.WithContext(ctx).Model(&model.UserModel{}).Where("role=?", role.String()).Count(&n).Error
	return n, errordef.ConvertGormError(err)
}
