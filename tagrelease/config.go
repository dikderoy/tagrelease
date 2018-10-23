package tagrelease

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Log struct {
		Debug bool
	}
	Branches struct {
		Master []string
		Trunk  []string
	}
	Strategy struct {
		Format    string
		Increment string
		Escape    string
	}
	WorkDir string
	Output  string
}

var GlobalConfig Config

func DefineConfig() {
	viper.SetTypeByDefaultValue(true)

	flag.StringSlice("rc", []string{"master"},
		"provide list of branches recognized as mainstream,"+
			" all releases on these branches will be marked as RC (release candidate),"+
			" except tagged ones")
	viper.BindPFlag("Branches.Master", flag.Lookup("rc"))

	flag.StringSlice("beta", []string{"trunk"},
		"provide list of branches recognized as trunked,"+
			" all releases on these branches will be marked as B (beta), except tagged ones")
	viper.BindPFlag("Branches.Trunk", flag.Lookup("beta"))

	flag.StringP("format", "f", FormatRelease,
		fmt.Sprintf("select output format: %v \n"+
			"or provide go-template string, available properties are: \n"+
			".SemVer|.Major|.Minor|.Patch|.Short|.Diff|.Revision|.RevisionShort", FormatList))

	viper.BindPFlag("Strategy.Format", flag.Lookup("format"))

	flag.StringP("increment", "i", StrategyUpstream,
		fmt.Sprintf("select increment strategy: %v", StrategyList))
	viper.BindPFlag("Strategy.Increment", flag.Lookup("increment"))

	flag.StringP("workdir", "d", ".",
		"select workdir to look for repository")
	viper.BindPFlag("WorkDir", flag.Lookup("workdir"))

	flag.StringP("output", "o", "-",
		"select output target, default is stdout")
	viper.BindPFlag("Output", flag.Lookup("output"))

	flag.Bool("debug", false, "enable debug output (to stderr)")
	viper.BindPFlag("Log.Debug", flag.Lookup("debug"))

	flag.StringP("escape", "e", "",
		"escape conflicting chars, some systems are sensitive to chars like +,/,~"+
			" and other which may occur in identifiers produced, use this option"+
			" to escape them with char provided")
	viper.BindPFlag("Strategy.Escape", flag.Lookup("escape"))
}

func LoadConfig() {
	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to load configuration")
	}
}

func InitLogger() {
	log.SetOutput(os.Stdout)
	if GlobalConfig.Log.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("enabled debug logging")
	}
	log.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation: true,
	})
}
