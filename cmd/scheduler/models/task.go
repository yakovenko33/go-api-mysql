package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CronTask struct {
	gorm.Model
	ID         uuid.UUID `gorm:"column:id"`
	Name       string    `yaml:"name"`
	RoutingKey string    `yaml:"routing_key"`
	Queue      string    `yaml:"queue"`
	Cron       string    `yaml:"cron"`
	RetryTTL   int       `yaml:"retry_ttl"`
	MaxRetries int       `yaml:"max_retries"`
	DLXEnabled bool      `yaml:"dlx_enabled"`
}

type User struct {
	gorm.Model
	ID             uuid.UUID `gorm:"column:id"`
	FirstName      string    `gorm:"column:first_name"`
	LastName       string    `gorm:"column:last_name"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	CreatedAt      string    `gorm:"column:created_at"`
	UpdatedAt      string    `gorm:"column:updated_at"`
	Status         uint8     `gorm:"column:status"`
	CreatedBy      string    `gorm:"column:created_by"`
	ModifiedUserId string    `gorm:"column:modified_user_id"`
}
