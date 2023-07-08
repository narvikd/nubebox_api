package cfg

import (
	"github.com/narvikd/errorskit"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"reflect"
)

// Config is a struct to unmarshall a .yml config file
type Config struct {
	Database struct {
		Ip     string `yaml:"ip"`
		Port   string `yaml:"port"`
		DBName string `yaml:"dbname"`
		User   string `yaml:"user"`
		//Password string `yaml:"password"`
	} `yaml:"database"`
}

// InitCfg reads a .yml config file returns a Config pointer.
//
// If an error is found it exits the application, except errClose which it only prints.
//
// In the absence of a config file, the application cannot continue any operation.
// Since it cannot connect to any of the necessary services. That's why the application needs to exit.
//
// The reason for errClose only printing instead of crashing the app,
// is that in the case it doesn't correctly handle the error,
// it will kill the application only because it couldn't close the file.
// This error is not critical or of "almost" any importance.
func InitCfg() *Config {
	// Handles the open/close of the file
	file, errOpen := os.Open("config.yml")
	if errOpen != nil {
		errorskit.FatalWrap(errOpen, "couldn't read configuration .yml file")

	}

	defer func(f *os.File) {
		errClose := f.Close()
		if errClose != nil {
			errorskit.LogWrap(errClose, "there was an error closing the configuration .yml file")
		}
	}(file)

	// Handles the decoding of the file into the Config struct
	decoder := yaml.NewDecoder(file)
	config := new(Config)
	errDecode := decoder.Decode(&config)
	if errDecode != nil {
		errorskit.FatalWrap(errDecode, "couldn't decode the configuration .yml file")
	}
	cfgCheckers(config)
	return config
}

// cfgCheckers calls all "config file" checkers to ensure it's all read as intended
func cfgCheckers(config *Config) {
	checkCfgEmpty(config)
	//checkFiberCfg(config)
}

// checkCfgEmpty iterates through a Config object and checks if any parameters are empty.
// If any of the struct values are empty, it crashes the application.
// More info on: https://www.reddit.com/r/golang/comments/mp7qqp/comment/gu8invd/?utm_source=share&utm_medium=web2x&context=3
func checkCfgEmpty(config *Config) {
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		for j := 0; j < v.Field(i).NumField(); j++ {
			structName := v.Type().Field(i).Name               // Name of the struct, for example Database
			originalVarName := v.Field(i).Type().Field(j).Name // Name of the var for example: IP
			composedVarName := structName + "_" + originalVarName
			varValue := v.Field(i).Field(j).String()
			if varValue == "" {
				log.Fatalln("couldn't read " + composedVarName + " in configuration .yml file")
			}
		}
	}
}
