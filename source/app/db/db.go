package db

import (
	"fmt"
	"gorm.io/gorm"
	"yscloudeBack/source/app/model"
	"yscloudeBack/source/app/utils"
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
	err := dm.dbEngine.AutoMigrate(&model.User{}, &model.Key{}, &model.FbToken{}, &model.Structure{}, &model.Notice{}, model.KeyPrice{})
	if err != nil {
		utils.Error(err.Error())
		return err
	}
	return nil
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

func (dm *DbManager) AddData(data interface{}) error {
	return dm.dbEngine.Create(&data).Error
}
func (dm *DbManager) SaveData(data interface{}) error {
	return dm.dbEngine.Save(&data).Error
}
