// contact/contact.go

package contact

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	UserID       uint   `gorm:"column:userId;" json:"userId"`
	FirstName    string `gorm:"firstName" json:"firstName"`
	LastName     string `gorm:"lastName" json:"lastName"`
	Email        string `gorm:"email" json:"email"`
	Phone        string `gorm:"phone" json:"phone"`
	AddressLine1 string `gorm:"addressLine1" json:"addressLine1"`
	AddressLine2 string `gorm:"addressLine2" json:"addressLine2"`
	City         string `gorm:"city" json:"city"`
	State        string `gorm:"state" json:"state"`
	Country      string `gorm:"country" json:"country"`
	Pincode      int    `gorm:"pincode" json:"pincode"`
}

var FieldsRequest struct {
	Fields []string `json:"fields"`
}
