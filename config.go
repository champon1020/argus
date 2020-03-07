package argus

import (
	"encoding/json"
	"fmt"
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

var GlobalConfig Config

func (config *Configurations) New(args string) {
	config.Load()

	if args == "" {
		GlobalConfig = config.Deploy
	} else if args == "stg" {
		GlobalConfig = config.Staging
	} else if args == "dev" {
		GlobalConfig = config.Dev
	} else {
		fmt.Printf("%s is not confortable, required '' or 'stg' or 'dev'.\n", args)
		os.Exit(0)
	}
}

func (config *Configurations) Load() {
	row, err := ioutil.ReadFile(os.Getenv("GOPATH") + "/src/github.com/champon1020/argus/config.json")
	if err != nil {
		//current, _ := os.Getwd()
		//Logger.ErrorMsgPrintf(current, err)
	}
	json.Unmarshal(row, &config)
}
