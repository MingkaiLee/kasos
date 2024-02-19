package model

import (
	"sort"
	"strings"

	"gorm.io/gorm"
)

type HpaService struct {
	gorm.Model
	ServiceName string   `gorm:"column:service_name;uniqueIndex"`
	Tags        string   `gorm:"column:tags;index"`
	Status      string   `gorm:"column:status"`
	ThreshQPS   uint     `gorm:"column:thresh_qps"`
	ModelId     uint     `gorm:"column:model_id"`
	ErrorInfo   string   `gorm:"column:error_info"`
	HpaModel    HpaModel `gorm:"foreignKey:ModelId"`
}

func (HpaService) TableName() string {
	return "hpa_service"
}

func HpaServiceCreate(serviceName string, tags map[string]string, modelId uint) error {
	h := &HpaService{
		ServiceName: serviceName,
		Tags:        convertTags(tags),
		Status:      statusTesting,
		ModelId:     modelId,
	}
	return db.Create(h).Error
}

func convertTags(tags map[string]string) string {
	kvs := make([]string, 0, len(tags))
	for t := range tags {
		kvs = append(kvs, t)
	}
	sort.Strings(kvs)
	for i := range kvs {
		kvs[i] = kvs[i] + "=" + tags[kvs[i]]
	}
	tagsStr := strings.Join(kvs, "&")

	return tagsStr
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
	return h, err
}

func HpaServiceList() ([]HpaService, error) {
	h := make([]HpaService, 0)
	err := db.Find(&h).Error
	return h, err
}

func HpaServiceDelete(serviceName string) error {
	return db.Where("service_name = ?", serviceName).Delete(&HpaService{}).Error
}
