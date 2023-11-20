package db

import "yscloudeBack/source/app/model"

func (dm *DbManager) GetStructures() ([]model.Structure, error) {
	var structures []model.Structure
	err := dm.dbEngine.Find(&structures).Error
	return structures, err
}

// 向数据库添加
func (dm *DbManager) AddStructure(structure model.Structure) error {
	return dm.dbEngine.Create(&structure).Error
}
func (dm *DbManager) GetStructureByHash(hashString string) (strcture model.Structure, err error) {
	result := dm.dbEngine.Where("file_hash = ?", hashString).First(&strcture)
	return strcture, result.Error
}
