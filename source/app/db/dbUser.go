package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"yscloudeBack/source/app/model"
)

func (dm *DbManager) CreateUser(user *model.User) error {
	if dm.dbEngine == nil {
		return fmt.Errorf("db is not existing")
	}
	result := dm.dbEngine.Create(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, user creation failed")
	}
	return nil
}
func (dm *DbManager) DeleteUser(username string) error {
	if dm.dbEngine == nil {
		return fmt.Errorf("db is not existing")
	}
	result := dm.dbEngine.Where("user_name = ?", username).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, user deletion failed")
	}
	return nil
}

// get all users from db
func (dm *DbManager) GetUsers() (users []model.User, err error) {
	result := dm.dbEngine.Preload("UserKeys").Preload("Slots").Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return
}
func (dm *DbManager) CheckErrorUserNotFound(err error) bool {
	return errors.Is(err, errors.New("user not found"))
}

func (dm *DbManager) GetUserByUserName(name string) (user *model.User, err error) {
	result := dm.dbEngine.Preload("UserKeys").Preload("Slots").Where("user_name = ?", name).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return user, nil
}

func (dm *DbManager) GetUserKeys(userID int) []model.Key {
	var user model.User
	dm.dbEngine.Preload("UserKeys").First(&user, userID)
	return user.UserKeys
}

// 关联 User 和 Key
func (dm *DbManager) AssociateKeyWithUser(userID uint, keyID uint) error {
	result := dm.dbEngine.Model(&model.Key{}).Where("id = ?", keyID).Update("user_id", userID)
	return result.Error
}
func (dm *DbManager) UpdateUserBalance(userId uint, value int) error {
	user := model.User{}
	db := dm.dbEngine
	err := db.Model(&model.User{}).Where("id = ?", userId).First(&user).Error
	if err != nil {
		return err
	}

	user.Balance = value
	err = db.Model(&user).Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

//
//// 关联 User 和 structure
//func (dm *DbManager) AssociateStuctureWithUser(userID uint, structureID uint) error {
//	result := dm.dbEngine.Model(&model.Structure{}).Where("id = ?", structureID).Update("structure_user_id", userID)
//	return result.Error
//}
