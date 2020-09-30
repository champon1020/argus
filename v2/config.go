package argus

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Errors returned at this file.
var (
	ErrConfReadFailed      = NewErrorType("arugs", "Failed to read configuration file")
	ErrConfUnmarshalFailed = NewErrorType("argus", "Failed to unmarshal into json")
)

// DbConf contains database configurations.
type DbConf struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Port   string `json:"port"`
	Host   string `json:"host"`
	DbName string `json:"dbname"`
}

// WebConf contains web configurations.
type WebConf struct {
	// Host is the http origin.
	Host string `json:"host"`

	// MaxDisplayArticle is the number of articles
	// which is displayed on the screen.
	MaxDisplayArticle int `json:"maxDisplayArticle"`

	// MaxDisplaySettingImage is the number of images
	// which is displayed on the setting image tab screen.
	MaxDisplaySettingImage int `json:"maxDisplaySettingImage"`

	// MaxDisplaySEttingArticle is the number of articles
	// which is displayed on the setting article tab screen.
	MaxDisplaySettingArticle int `json:"maxDisplaySettingArticle"`
}

// Config contains all configurations.
type Config struct {
	Db  DbConf  `json:"db"`
	Web WebConf `json:"web"`
}

// Configs contains some configurations splitted by argus mode.
type Configs struct {
	Deploy Config `json:"deploy"`
	Dev    Config `json:"dev"`
}

// NewConfig initialize the struct of Config.
func NewConfig() *Config {
	// Load from config file.
	configs := new(Configs)
	if err := configs.load(); err != nil {
		os.Exit(1)
	}

	// Config will be returned.
	var config Config

	// Branch by the project built mode.
	switch Env.Get("mode") {
	case "deploy":
		config = configs.Deploy
	case "dev":
		config = configs.Dev
	case "test":
		// If test mode, configuration file wouldn't be mounted on CI.
		// So this function returns the mock Config.
		return &Config{
			Db: DbConf{},
			Web: WebConf{
				Host:                     "",
				MaxDisplayArticle:        4,
				MaxDisplaySettingImage:   12,
				MaxDisplaySettingArticle: 10,
			},
		}
	}

	// Assign the database configuration.
	config.Db = DbConf{
		User:   Env.Get("dbUser"),
		Pass:   Env.Get("dbPass"),
		Port:   Env.Get("dbPort"),
		Host:   Env.Get("dbHost"),
		DbName: Env.Get("dbName"),
	}

	return &config
}

// Load from configuration file.
func (config *Configs) load() error {
	var row []byte

	configPath := Env.Get("config")
	row, err := ioutil.ReadFile(configPath)
	if err != nil {
		return NewError(ErrConfReadFailed, err)
	}

	if err := json.Unmarshal(row, &config); err != nil {
		return NewError(ErrConfUnmarshalFailed, err)
	}

	return nil
}
