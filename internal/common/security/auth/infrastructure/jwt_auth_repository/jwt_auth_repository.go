package jwt_auth_repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JwtAuthRepositoryInterface interface {
	AddTokens(entity *Tokens) (string, error)
	//DeleteTokens() (string, error)
	RefreshTokenExist(refreshToken string) (*Tokens, error)
}

type Tokens struct {
	gorm.Model
	ID           uuid.UUID `gorm:"column:id"`
	AccessToken  string    `gorm:"column:access_token"`
	RefreshToken string    `gorm:"column:refresh_token"`
	UserAgent    string    `gorm:"column:user_agent"`
	ExpiresIn    time.Time `gorm:"type:datetime;column:expires_in"`
	CreatedAt    time.Time `gorm:"type:datetime;column:created_at"`
}

type JwtAuthRepository struct {
	db *gorm.DB
}

func NewJwtAuthRepository(db *gorm.DB) *JwtAuthRepository {
	return &JwtAuthRepository{
		db: db,
	}
}

func (m JwtAuthRepository) AddTokens(entity *Tokens) (string, error) {
	entity.ID = uuid.New()

	result := m.db.Create(entity)
	if result.Error != nil {
		return "", result.Error
	}

	return entity.ID.String(), nil
}

func (m JwtAuthRepository) RefreshTokenExist(refreshToken string) (*Tokens, error) {
	var tokens Tokens
	result := m.db.Where("refresh_token = ?", refreshToken).Find(&tokens)

	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &tokens, nil
}
