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
	Host                     string `json:"host"`
	Pickup                   []int  `json:"pickup"`
	MaxViewArticleNum        int    `json:"maxViewArticleNum"`
	MaxViewImageNum          int    `json:"maxViewImageNum"`
	MaxViewSettingArticleNum int    `json:"maxViewSettingArticleNum"`
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

func NewConfig() *Config {
	configurations := new(Configurations)
	if err := configurations.load(); err != nil {
		StdLogger.ErrorLog(Errors)
		os.Exit(1)
	}

	config := new(Config)
	if EnvVars.Get("mode") == "deploy" {
		*config = configurations.Deploy
		Logger.Println("deploy mode")
		return config
	}
	if EnvVars.Get("mode") == "staging" {
		*config = configurations.Staging
		Logger.Println("staging mode")
		return config
	}
	*config = configurations.Dev
	Logger.Println("dev mode")
	return config
}

func (config *Configurations) load() (err error) {
	if EnvVars.Get("mode") == "test" {
		config.Dev = Config{
			Db: DbConf{},
			Web: WebConf{
				MaxViewArticleNum:        4,
				MaxViewImageNum:          12,
				MaxViewSettingArticleNum: 10,
			},
		}
		return
	}

	var row []byte
	configPath := EnvVars.Get("config")
	if row, err = ioutil.ReadFile(configPath); err != nil {
		ConfigLoadError.
			SetErr(err).
			SetValues("configPath", configPath).
			AppendTo(&Errors)
		return
	}

	json.Unmarshal(row, &config)
	return
}
