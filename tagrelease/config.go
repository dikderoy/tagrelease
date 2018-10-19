package tagrelease

import (
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Branches struct {
		Master []string
		Trunk  []string
	}
	Log struct {
		Level  string
		Colors bool
	}
}

var GlobalConfig Config

var configEnvPrefix string

func fullEnvVarName(name string) string {
	if configEnvPrefix != "" {
		name = configEnvPrefix + "_" + name
	}
	return name
}

func defineConfigValue(key string, defaultValue interface{}, envVarName string) {
	viper.SetDefault(key, defaultValue)
	viper.BindEnv(key, fullEnvVarName(envVarName))
}

func DefineConfig() {
	viper.SetTypeByDefaultValue(true)
	defineConfigValue("Log.Level", "info", "LOG_LEVEL")
	viper.AutomaticEnv()

}

func DefineCommandLineConfig() {
	flag.StringVar(
		&configEnvPrefix, "env-prefix", "BDZ",
		"redefine env prefix for conflicting environments",
	)
	viper.SetEnvPrefix(configEnvPrefix)

	flag.String("bin", "", "binary name to run as supervised process")
	viper.BindPFlag("process.name", flag.Lookup("bin"))

	flag.StringSlice("args", []string{}, "program arguments to pass at launch time")
	viper.BindPFlag("process.args", flag.Lookup("args"))
}

func LoadConfig() {
	DefineConfig()
	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to load configuration")
	}
}

func InitLogger() {
	log.SetOutput(os.Stdout)
	level, err := log.ParseLevel(GlobalConfig.Log.Level)
	if err != nil {
		log.WithError(err).
			WithField("input", GlobalConfig.Log.Level).
			Panic("cannot determine log level")
	}
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            GlobalConfig.Log.Colors,
		DisableLevelTruncation: true,
	})
}
