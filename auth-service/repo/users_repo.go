package repo

import (
	"auth_service/entity"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetByEmail(email string) (*entity.User, error)
}

type UsersRepoImpl struct {
	Db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UsersRepoImpl{
		Db: db,
	}
}

func (repo *UsersRepoImpl) GetByEmail(email string) (*entity.User, error) {
	model := new(entity.User)
	if err := repo.Db.Where("email = ?", email).First(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}
