package cmd

import (
	"strings"

	"github.com/spf13/viper"
	"go.infratographer.com/x/loggingx"
	"go.uber.org/zap"
)

func initLogger() *zap.Logger {
	sl := loggingx.InitLogger("lmi", loggingx.Config{
		Debug:  viper.GetBool("debug"),
		Pretty: viper.GetBool("pretty"),
	})

	return sl.Desugar()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("lmi")

	viper.AutomaticEnv() // read in environment variables that match
}
