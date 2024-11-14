package user

import (
	"addressBook/models/contact"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// UserID   uint              `gorm:"column:userId;primarykey" json:"userId"`
	Email    string            `gorm:"column:email;unique" json:"email" validate:"required,email"`
	Password string            `gorm:"password" json:"password" validate:"required,min=5,max=10"`
	Contacts []contact.Contact `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
