package contactrepo

import (
	"addressBook/models/contact"

	"addressBook/models/user"

	"gorm.io/gorm"
)

type ContactRepo struct {
	db *gorm.DB
}

func NewContactRepo(sqldb *gorm.DB) *ContactRepo {
	return &ContactRepo{db: sqldb}
}

func (contactRepo *ContactRepo) CreateContact(newContact *contact.Contact, userEmail string) error {
	var user user.User
	err := contactRepo.db.Select("ID").Where("email = ?", userEmail).Find(&user).Error
	if err != nil {
		return err
	}

	// Set the user ID in the newContact
	newContact.UserID = user.ID
	// Now perform the actual create operation
	err = contactRepo.db.Create(newContact).Error
	if err != nil {
		return err
	}

	return nil
}

func (contactRepo *ContactRepo) GetContacts(userEmail string) ([]contact.Contact, error) {

	contactList := []contact.Contact{}
	var user user.User

	err := contactRepo.db.Select("ID").Where("email = ?", userEmail).Find(&user).Error
	if err != nil {
		return nil, err
	}

	err = contactRepo.db.Where("userId = ?", user.ID).Find(&contactList).Error
	if err != nil {
		return nil, err
	}

	return contactList, err
}

func (contactRepo *ContactRepo) UpdateContact(cid uint, updatedContact contact.Contact) error {

	err := contactRepo.db.Model(&contact.Contact{}).Where("ID = ?", cid).Updates(updatedContact).Error
	if err != nil {

		return err
	}
	return nil
}

func (contactRepo *ContactRepo) DeleteContact(contactId uint) error {

	err := contactRepo.db.Delete(&contact.Contact{}, contactId).Error

	if err != nil {
		return err
	}
	return nil
}
