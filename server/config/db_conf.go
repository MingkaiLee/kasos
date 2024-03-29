package config

import (
	"fmt"
	"os"

	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

var DSN string

const dbConfFile = "/etc/config/db.json"

type DBConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DBName   string `json:"dbname"`
	Charset  string `json:"charset"`
}

func initDBConf() {
	var conf DBConf
	var err error

	d, err := os.ReadFile(dbConfFile)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	err = jsoniter.Unmarshal(d, &conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}

	DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName, conf.Charset)
}
