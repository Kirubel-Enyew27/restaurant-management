package intiator

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(name, path string, log *zap.Logger) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to read config: %v", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed:", zap.String("file", e.Name))
	})
}
