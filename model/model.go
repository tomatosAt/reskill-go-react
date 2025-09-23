package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/tomatosAt/reskill-go-react/pkg/util"
	"gorm.io/gorm"
)

type Model struct {
	Id        uuid.UUID       `gorm:"column:id;primary_key;index"`
	CreatedAt time.Time       `gorm:"column:created_at;index"`
	UpdatedAt time.Time       `gorm:"column:updated_at;index"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CreatedBy string          `gorm:"column:created_by;size:50;"`
	UpdatedBy string          `gorm:"column:updated_by;size:50;"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.Id == uuid.Nil {
		m.Id = util.GenUniqueIdV7()
	}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}
