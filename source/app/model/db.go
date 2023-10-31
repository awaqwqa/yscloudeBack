package model

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type DbManager struct {
	dbEngine *gorm.DB
}

func NewDbManager(r *gorm.DB) *DbManager {
	return &DbManager{
		dbEngine: r,
	}
}
func (dm *DbManager) Init() error {
	if dm.dbEngine == nil {
		return fmt.Errorf("db is not exiting,maybe init() will help you")
	}
	err := dm.dbEngine.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	err = dm.dbEngine.AutoMigrate(&Key{})
	if err != nil {
		return err
	}
	err = dm.AddKey("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		return err
	}
	return nil
}

// get all users from db
func (dm *DbManager) GetUsers() (users []User, err error) {
	result := dm.dbEngine.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}
	return
}
func (dm *DbManager) CheckErrorUserNotFound(err error) bool {
	return errors.Is(err, errors.New("user not found"))
}
func (dm *DbManager) GetUserByUserName(name string) (user *User, err error) {
	result := dm.dbEngine.Where("user_name = ?", name).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return user, nil
}
func (dm *DbManager) GetUserByToken(token string) (*User, error) {
	var user User
	result := dm.dbEngine.Where("token = ?", token).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}
func (dm *DbManager) GetUserByFBToken(fbToken string) (*User, error) {
	var user User
	result := dm.dbEngine.Where("fb_token = ?", fbToken).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update 更新数据库中的记录。
// model 参数是一个指向结构体的指针，它包含了要更新的字段。
// conditions 是一个map，包含了用于查找记录的条件。
func (dm *DbManager) UpdateByConditions(model interface{}, conditions map[string]interface{}) error {
	if dm.dbEngine == nil {
		return fmt.Errorf("db is not existing")
	}

	// 使用Where方法添加条件，然后使用Updates方法更新记录。
	result := dm.dbEngine.Where(conditions).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateUserByUsername 根据用户名更新用户信息 传入新的结构体即可
func (dm *DbManager) UpdateUserByUsername(username string, updatedUser *User) error {
	if dm.dbEngine == nil {
		return fmt.Errorf("db is not existing")
	}

	// 首先，尝试寻找对应用户名的用户
	user, err := dm.GetUserByUserName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// 更新找到的用户信息
	result := dm.dbEngine.Model(user).Updates(updatedUser)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, update failed")
	}

	return nil
}

func (dm *DbManager) CreateUser(user *User) error {
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
	result := dm.dbEngine.Where("user_name = ?", username).Delete(&User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, user deletion failed")
	}
	return nil
}
func (dm *DbManager) AddKey(value string) error {
	key := Key{Value: value}
	return dm.dbEngine.Create(&key).Error
}
func (dm *DbManager) CheckKey(key string) (bool, error) {
	var count int64
	err := dm.dbEngine.Model(&Key{}).Where("value = ?", key).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dm *DbManager) GetKeyCount() (int64, error) {
	var count int64
	err := dm.dbEngine.Model(&Key{}).Count(&count).Error
	return count, err
}
func (dm *DbManager) GetAllKeys() ([]Key, error) {
	var keys []Key
	err := dm.dbEngine.Find(&keys).Error
	return keys, err
}

func (dm *DbManager) DeleteKey(value string) error {
	return dm.dbEngine.Where("value = ?", value).Delete(&Key{}).Error
}
