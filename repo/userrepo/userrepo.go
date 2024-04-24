package repo

import (
	"addressBook/models/user"
	"crypto/md5"
	"encoding/hex"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(sqldb *gorm.DB) *UserRepo {
	return &UserRepo{db: sqldb}
}

func (repo *UserRepo) RegisterUser(email, password string) error {

	userModel := user.User{
		Email:    email,
		Password: password,
	}

	// fmt.Println("hashString : ", hashString)

	err := repo.db.Create(&userModel).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) LoginUser(email, password string) (bool, error) {

	var loginUser user.User
	result := repo.db.Where("email = ?", email).First(&loginUser)
	if result.Error != nil || result.RowsAffected != 1 {

		return false, result.Error
	}
	passwordBytes := []byte(password)
	hash := md5.Sum(passwordBytes)
	hashString := hex.EncodeToString(hash[:])

	if loginUser.Password != hashString {
		return false, nil
	}
	return true, nil

}
