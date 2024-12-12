package gorm

import (
	"context"

	"github.com/todennus/shared/enumdef"
	"github.com/todennus/shared/errordef"
	"github.com/todennus/shared/xcontext"
	"github.com/todennus/user-service/domain"
	"github.com/todennus/user-service/infras/database/model"
	"github.com/xybor-x/snowflake"
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
	return errordef.ConvertGormError(xcontext.DB(ctx, repo.db).Create(&model).Error)
}

func (repo *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	model := model.UserModel{}
	if err := xcontext.DB(ctx, repo.db).Take(&model, "username=?", username).Error; err != nil {
		return nil, errordef.ConvertGormError(err)
	}

	return model.To()
}

func (repo *UserRepository) GetByID(ctx context.Context, userID snowflake.ID) (*domain.User, error) {
	model := model.UserModel{}
	if err := xcontext.DB(ctx, repo.db).Take(&model, "id=?", userID).Error; err != nil {
		return nil, errordef.ConvertGormError(err)
	}

	return model.To()
}

func (repo *UserRepository) GetAvatarByID(ctx context.Context, userID snowflake.ID) (snowflake.ID, error) {
	model := model.UserModel{}
	if err := xcontext.DB(ctx, repo.db).Select("avatar").Take(&model, "id=?", userID).Error; err != nil {
		return 0, errordef.ConvertGormError(err)
	}

	return snowflake.ParseInt64(model.Avatar), nil
}

func (repo *UserRepository) UpdateAvatarByID(ctx context.Context, userID, avatar snowflake.ID) error {
	return errordef.ConvertGormError(
		xcontext.DB(ctx, repo.db).Model(&model.UserModel{}).
			Where("id=?", userID).
			Update("avatar", avatar).Error,
	)
}

func (repo *UserRepository) CountByRole(ctx context.Context, role enumdef.UserRole) (int64, error) {
	var n int64
	err := xcontext.DB(ctx, repo.db).
		Model(&model.UserModel{}).
		Where("role=?", role).
		Count(&n).Error
	return n, errordef.ConvertGormError(err)
}
