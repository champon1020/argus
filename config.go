package argus

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	errConfReadFailed      = errors.New("argus.config: Failed to read configuration file")
	errConfUnmarshalFailed = errors.New("argus.config: Failed to unmarshal configuration")
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

// Configuration contains all configurations.
type Configuration struct {
	Db  DbConf  `json:"db"`
	Web WebConf `json:"web"`
}

// Configurations contains some configurations splitted by argus mode.
type Configurations struct {
	Deploy Configuration `json:"deploy"`
	Dev    Configuration `json:"dev"`
}

// NewConfig initialize the struct of Config.
func NewConfig() *Configuration {
	// Load from config file.
	configs := new(Configurations)
	if err := configs.load(); err != nil {
		os.Exit(1)
	}

	// Config will be returned.
	var config Configuration

	// Branch by the project built mode.
	switch Env.Get("mode") {
	case "deploy":
		config = configs.Deploy
	case "dev":
		config = configs.Dev
	case "test":
		// If test mode, configuration file wouldn't be mounted on CI.
		// So this function returns the mock Config.
		return &Configuration{
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
func (config *Configurations) load() error {
	var row []byte

	configPath := Env.Get("config")
	row, err := ioutil.ReadFile(configPath)
	if err != nil {
		return NewError(errConfReadFailed, err).
			AppendValue("config path", configPath)
	}

	if err := json.Unmarshal(row, &config); err != nil {
		return NewError(errConfUnmarshalFailed, err).
			AppendValue("config", config)
	}

	return nil
}
