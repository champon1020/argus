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
	MaxViewImageNum   int `json:"maxViewImageNum"`
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

var (
	GlobalConfig    Config
	ConfigLoadError = NewError(ConfigFailedLoadError)
)

func NewConfig(args string) Config {
	configurations := new(Configurations)
	if err := configurations.load(); err != nil {
		StdLogger.ErrorLog(Errors)
		os.Exit(1)
	}

	config := new(Config)
	if args == "" {
		*config = configurations.Deploy
	} else if args == "stg" {
		*config = configurations.Staging
	} else if args == "dev" {
		*config = configurations.Dev
	} else {
		StdLogger.Fatalf("%s is not confortable, required '' or 'stg' or 'dev'.\n", args)
	}
	return *config
}

func (config *Configurations) load() (err error) {
	if os.Getenv("IS_TRAVIS") == "on" {
		return
	}

	var row []byte
	configPath := EnvVars.Get("config")
	if row, err = ioutil.ReadFile(configPath); err != nil {
		current, _ := os.Getwd()
		ConfigLoadError.
			SetErr(err).
			SetValues("current path", current).
			AppendTo(&Errors)
		return
	}

	json.Unmarshal(row, &config)
	return
}
