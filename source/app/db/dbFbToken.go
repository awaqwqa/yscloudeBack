package db

import "yscloudeBack/source/app/model"

func (dm *DbManager) GetFbTokens() ([]model.FbToken, error) {
	var tokens []model.FbToken
	err := dm.dbEngine.Find(&tokens).Error
	return tokens, err
}

// 删除fbToken根据value
func (dm *DbManager) DeleteFbToken(value string) error {
	return dm.dbEngine.Where("value = ?", value).Delete(&model.FbToken{}).Error
}

// 向数据库添加fbToken密钥
func (dm *DbManager) AddFbToken(fbToken model.FbToken) error {
	return dm.dbEngine.Create(&fbToken).Error
}
