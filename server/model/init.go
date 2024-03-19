package model

import (
	"github.com/MingkaiLee/kasos/server/config"
	"github.com/MingkaiLee/kasos/server/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

const (
	PageSize      = 10
	statusTesting = "testing"
	statusOk      = "ok"
	statusError   = "error"
)

func InitModel() {
	var err error

	db, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	err = db.AutoMigrate(&HpaService{}, &HpaModel{})
	if err != nil {
		util.LogErrorf("panic: %v", err)
	}
}
