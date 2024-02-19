package config

import (
	"fmt"
	"os"

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
		panic(err)
	}
	err = jsoniter.Unmarshal(d, &conf)
	if err != nil {
		panic(err)
	}

	DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName, conf.Charset)
}
