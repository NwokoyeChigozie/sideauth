package models

import (
	"fmt"
	"time"

	"github.com/vesicash/auth-ms/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func (ContactUs) TableName() string {
	return "contact_us"
}

type ContactUs struct {
	ID           uint      `gorm:"column:id; type:uint; not null; primaryKey; unique; autoIncrement" json:"id"`
	FirstName    string    `gorm:"column:first_name; type:varchar(250); not null" json:"first_name"`
	LastName     string    `gorm:"column:last_name; type:varchar(250); not null" json:"last_name"`
	Email        string    `gorm:"column:email; type:varchar(250); not null" json:"email"`
	WebsiteUrl   string    `gorm:"column:website_url; type:varchar(250); not null" json:"website_url"`
	BusinessType string    `gorm:"column:business_type; type:varchar(250); not null" json:"business_type"`
	Country      string    `gorm:"column:country; type:varchar(250); not null" json:"country"`
	Message      string    `gorm:"column:message; type:text; not null" json:"message"`
	DeletedAt    time.Time `gorm:"column:deleted_at;" json:"deleted_at"`
	CreatedAt    time.Time `gorm:"column:created_at; autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at; autoUpdateTime" json:"updated_at"`
}

type ContactUsCreateModel struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required"`
	WebsiteUrl   string `json:"website_url" validate:"required,url"`
	BusinessType string `json:"business_type" validate:"required,oneof=ecommerce social_commerce marketplace"`
	Country      string `json:"country" validate:"required"`
	Message      string `json:"message" validate:"required"`
}

func (c *ContactUs) CreateContactUs(db *gorm.DB) error {
	err := postgresql.CreateOneRecord(db, &c)
	if err != nil {
		return fmt.Errorf("contact us creation failed: %v", err.Error())
	}
	return nil
}
