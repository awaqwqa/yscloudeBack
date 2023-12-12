package db

import (
	"errors"
	"gorm.io/gorm"
	"yscloudeBack/source/app/model"
)

func (db *DbManager) GetAnnouncementByID(id uint) (*model.Notice, error) {
	var announcement model.Notice
	result := db.dbEngine.First(&announcement, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 如果找不到公告，则自动创建一个新的公告
			announcement = model.Notice{
				ID:    id,
				Value: "默认内容",
			}
			if err := db.dbEngine.Create(&announcement).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, result.Error
		}
	}
	return &announcement, nil
}
func (db *DbManager) UpdateAnnouncement(id uint, newContent string) error {
	announcement, err := db.GetAnnouncementByID(id)
	if err != nil {
		return err
	}

	announcement.Value = newContent
	result := db.dbEngine.Save(&announcement)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
