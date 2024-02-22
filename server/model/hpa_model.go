package model

import (
	"gorm.io/gorm"
)

type HpaModel struct {
	gorm.Model
	ModelName   string `gorm:"column:model_name"`
	Status      string `gorm:"column:status"`
	TrainScript string `gorm:"column:train_script"`
	InferScript string `gorm:"column:infer_script"`
	ErrorInfo   string `gorm:"column:error_info"`
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

func HpaModelList() ([]HpaModel, error) {
	var h []HpaModel
	err := db.Where("status = ?", statusOk).Order("id").Find(&h).Error
	return h, err
}

func HpaModelDelete(modelName string) error {
	return db.Where("model_name = ?", modelName).Delete(&HpaModel{}).Error
}
