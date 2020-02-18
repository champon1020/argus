package argus

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DbConf struct {
	User string `json: "user"`
	Pass string `json: "pass"`
	Port string `json: "port"`
	Host string `json: "host"`
}

type WebConf struct {
	MaxViewArticleNum int `json: "maxViewArticleNum"`
}

type Config struct {
	Db    DbConf  `json: "db"`
	DevDb DbConf  `json: "devDb"`
	Web   WebConf `json: "web"`
}

func (config *Config) Load() {
	row, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/champon1020/argus/config.json")
	if err != nil {
		current, _ := os.Getwd()
		logger.ErrorMsgPrintf(current, err)
	}
	json.Unmarshal(row, &config)
}
