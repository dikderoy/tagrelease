package tagrelease

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	GlobalConfig.Log.Debug = true
	InitLogger()

	os.Exit(m.Run())
}
