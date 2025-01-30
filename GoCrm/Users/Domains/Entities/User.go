package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uuid.UUID `gorm:"column:id"`
	FirstName      string    `gorm:"column:first_name"`
	LastName       string    `gorm:"column:last_name"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	CreatedAt      string    `gorm:"column:created_at"`
	UpdatedAt      string    `gorm:"column:updated_at"`
	Status         string    `gorm:"column:status"`
	CreatedBy      string    `gorm:"column:created_by"`
	ModifiedUserId string    `gorm:"column:modified_user_id"`
}
