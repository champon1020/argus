package argus

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DbConf struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Port   string `json:"port"`
	Host   string `json:"host"`
	DbName string `json:"dbname"`
}

type WebConf struct {
	MaxViewArticleNum int `json:"maxViewArticleNum"`
}

type Config struct {
	Db  DbConf  `json:"db"`
	Web WebConf `json:"web"`
}

type Configurations struct {
	Deploy  Config `json:"deploy"`
	Staging Config `json:"staging"`
	Dev     Config `json:"dev"`
}

func (config *Configurations) Load() Config {
	row, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/champon1020/argus/config.json")
	if err != nil {
		current, _ := os.Getwd()
		logger.ErrorMsgPrintf(current, err)
	}
	json.Unmarshal(row, &config)
	return config.Dev
}
