package tagrelease

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	configEnvPrefix = "TR"

	DefineConfig()
	LoadConfig()

	GlobalConfig.Log.Level = "debug"
	InitLogger()

	os.Exit(m.Run())
}
