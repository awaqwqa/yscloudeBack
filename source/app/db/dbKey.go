package db

import (
	"errors"
	"gorm.io/gorm"
	"yscloudeBack/source/app/model"
)

// 向数据库添加key密钥
func (dm *DbManager) AddKey(key model.Key) error {
	return dm.dbEngine.Create(&key).Error
}

// 这里其实FindKeyByValue 功能较为一致 区别就是一个返回key一个不返回key
func (dm *DbManager) CheckKey(key string) (bool, error) {
	var count int64
	err := dm.dbEngine.Model(&model.Key{}).Where("value = ?", key).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dm *DbManager) GetKeyCount() (int64, error) {
	var count int64
	err := dm.dbEngine.Model(&model.Key{}).Count(&count).Error
	return count, err
}

// 获取所有的key
func (dm *DbManager) GetAllKeys() ([]model.Key, error) {
	var keys []model.Key
	err := dm.dbEngine.Find(&keys).Error
	return keys, err
}

// 删除key根据value
func (dm *DbManager) DeleteKey(value string) error {
	return dm.dbEngine.Where("value = ?", value).Delete(&model.Key{}).Error
}

// 从数据库根据value获取key的值
func (dm *DbManager) GetKeyByValue(value string) (*model.Key, error) {
	var key model.Key
	result := dm.dbEngine.Where("value = ?", value).First(&key)
	return &key, result.Error
}
func (dm *DbManager) UpdateKeyStatus(key *model.Key, status bool) error {
	return dm.dbEngine.Model(key).Update("Status", status).Error
}
func (dm *DbManager) UpdateKeyFileGroupName(key *model.Key, name string) error {
	return dm.dbEngine.Model(key).Update("FileGroup", name).Error
}
func (db *DbManager) GetKeyPriceByID(id uint) (*model.KeyPrice, error) {
	var Value model.KeyPrice
	result := db.dbEngine.First(&Value, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果找不到keyValue，则自动创建一个新的公告
			Value = model.KeyPrice{
				ID:    id,
				Value: 4,
			}
			if err := db.dbEngine.Create(&Value).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}
	return &Value, nil
}
func (db *DbManager) UpdateKeyPrice(id uint, newContent int) error {
	Value, err := db.GetKeyPriceByID(id)
	if err != nil {
		return err
	}

	Value.Value = newContent
	result := db.dbEngine.Save(&Value)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
