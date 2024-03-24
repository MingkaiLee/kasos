package model

import (
	"gorm.io/gorm"
)

type HpaModel struct {
	gorm.Model
	ModelName   string `gorm:"column:model_name;type:VARCHAR(128);uniqueIndex"`
	Status      string `gorm:"column:status;type:VARCHAR(32)"`
	TrainScript string `gorm:"column:train_script;type:TEXT"`
	InferScript string `gorm:"column:infer_script;type:TEXT"`
	ErrorInfo   string `gorm:"column:error_info;type:VARCHAR(256)"`
}

func (HpaModel) TableName() string {
	return "hpa_model"
}

func HpaModelCreate(modelName string, trainScript string, inferScript string) error {
	h := &HpaModel{
		ModelName:   modelName,
		Status:      statusTesting,
		TrainScript: trainScript,
		InferScript: inferScript,
	}
	return db.Create(h).Error
}

func HpaModelRecordOk(modelName string) error {
	return db.Model(&HpaModel{}).Where("model_name = ?", modelName).Update("status", statusOk).Error
}

func HpaModelRecordError(modelName string, errorInfo string) error {
	return db.Model(&HpaModel{}).Where("model_name = ?", modelName).Update("status", statusError).Update("error_info", errorInfo).Error
}

func HpaModelGet(modelName string) (*HpaModel, error) {
	h := &HpaModel{}
	err := db.Where("model_name = ?", modelName).First(h).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return h, err
}

func HpaModelGetID(modelName string) (uint, error) {
	h := &HpaModel{}
	err := db.Model(&HpaModel{}).Where("model_name = ?", modelName).Select("id").First(h).Error
	return h.ID, err
}

func HpaModelList(start_idx uint) ([]HpaModel, error) {
	var h []HpaModel
	err := db.Where("status = ? AND id >= ?", statusOk, start_idx).Order("id").Limit(PageSize).Find(&h).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return h, err
}

func HpaModelDelete(modelName string) error {
	return db.Unscoped().Delete(&HpaModel{
		ModelName: modelName,
	}).Error
}
