package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CronTask struct {
	gorm.Model
	ID         uuid.UUID `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	RoutingKey string    `gorm:"column:routing_key"`
	//Queue      string    `gorm:"column:queue"`
	Cron       string `gorm:"column:cron"`
	RetryTTL   int    `gorm:"column:retry_ttl"`
	MaxRetries int    `gorm:"column:max_retries"`
	DLXEnabled bool   `gorm:"column:dlx_enabled"`
}
