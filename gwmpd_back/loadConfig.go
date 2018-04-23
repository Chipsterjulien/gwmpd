package main

import (
	// "fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"
	"path"
)

func loadConfig(filenamePath *string, filename *string) {
	log := logging.MustGetLogger("log")

	viper.SetConfigName(*filename)
	viper.AddConfigPath(*filenamePath)

	if err := viper.ReadInConfig(); err != nil {
		log.Criticalf("Unable to load config \"%s\" file: %v", path.Join(*filenamePath, *filename), err)
		os.Exit(1)
	}

	switch viper.GetString("logtype") {
	case "critical":
		logging.SetLevel(0, "")
	case "error":
		logging.SetLevel(1, "")
	case "warning":
		logging.SetLevel(2, "")
	case "notice":
		logging.SetLevel(3, "")
	case "info":
		logging.SetLevel(4, "")
	case "debug":
		logging.SetLevel(5, "")
		log.Debug("\"debug\" is selected")
	default:
		logging.SetLevel(2, "")
	}

	log.Debug("loadConfig func:")
	log.Debugf("  path: %s", *filenamePath)
	log.Debugf("  filename: %s", *filename)
	log.Debugf("  logtype in file config is \"%s\"", viper.GetString("logtype"))
}
