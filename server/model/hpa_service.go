package model

import (
	"github.com/MingkaiLee/kasos/server/util"
	"gorm.io/gorm"
)

type HpaService struct {
	gorm.Model
	ServiceName string `gorm:"column:service_name;type:VARCHAR(128);uniqueIndex"`
	Tags        string `gorm:"column:tags;type:VARCHAR(1024)"`
	Status      string `gorm:"column:status;type:VARCHAR(32)"`
	ThreshQPS   uint   `gorm:"column:thresh_qps;type:BIGINT"`
	ModelName   string `gorm:"column:model_name;type:VARCHAR(128)"`
	ErrorInfo   string `gorm:"column:error_info;type:VARCHAR(256)"`
}

func (HpaService) TableName() string {
	return "hpa_service"
}

func HpaServiceCreate(serviceName string, tags map[string]string, modelName string) error {
	h := &HpaService{
		ServiceName: serviceName,
		Tags:        util.ConvertTags(tags),
		Status:      statusTesting,
		ModelName:   modelName,
	}
	return db.Create(h).Error
}

func HpaServiceRecordError(serviceName string, errorInfo string) error {
	return db.Model(&HpaService{}).Where("service_name = ?", serviceName).Updates(map[string]interface{}{"error_info": errorInfo, "status": statusError}).Error
}

func HpaServiceRecordThreshQPS(serviceName string, threshQPS uint) error {
	return db.Model(&HpaService{}).Where("service_name = ?", serviceName).Updates(map[string]interface{}{"thresh_qps": threshQPS, "status": statusOk}).Error
}

func HpaServiceChangeModel(serviceName string, modelId uint) error {
	return db.Model(&HpaService{}).Where("service_name = ?", serviceName).Updates(map[string]interface{}{"model_id": modelId}).Error
}

func HpaServiceGet(serviceName string) (*HpaService, error) {
	h := &HpaService{}
	err := db.Where("service_name = ?", serviceName).First(h).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return h, err
}

func HpaServiceList(start_idx uint) ([]HpaService, error) {
	h := make([]HpaService, 0)
	err := db.Where("status = ? AND id >= ?", statusOk, start_idx).Order("id").Limit(PageSize).Find(&h).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return h, err
}

func HpaServiceDelete(serviceName string) error {
	return db.Unscoped().Delete(&HpaService{
		ServiceName: serviceName,
	}).Error
}
