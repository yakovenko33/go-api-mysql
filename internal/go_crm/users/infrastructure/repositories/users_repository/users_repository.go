package users_repository

import (
	"gorm.io/gorm"

	users_entities "go-api-docker/internal/go_crm/users/domains/entities"
)

type UsersRepositoryInterface interface {
	FindUserByEmail(email string) (*users_entities.User, error)
}

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRrepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (m UsersRepository) FindUserByEmail(email string) (*users_entities.User, error) {
	var user users_entities.User
	result := m.db.Where("email = ?", email).Find(&user)

	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
