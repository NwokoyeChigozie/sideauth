package models

import (
	"fmt"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

type BusinessType struct {
	ID        uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	Name      string    `gorm:"column:name; type:varchar(255); not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

func (b *BusinessType) CreateBusinessType(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &b)
	if err != nil {
		return fmt.Errorf("BusinessType creation failed: %v", err.Error())
	}
	return nil
}

func (b *BusinessType) GetBusinessTypes(db *gorm.DB) ([]BusinessType, error) {
	types := []BusinessType{}
	err := postgresql.SelectAllFromDb(db, "desc", &types, "")
	if err != nil {
		return types, err
	}
	return types, nil
}
