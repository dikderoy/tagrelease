package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/dikderoy/tagrelease/tagrelease"
)

func rootCommand() {
	if tagrelease.GlobalConfig.WorkDir != "." {
		err := os.Chdir(tagrelease.GlobalConfig.WorkDir)
		if err != nil {
			logrus.WithError(err).Fatal("cannot change working directory")
		}
	}

	strategy, err := tagrelease.StrategyFactory(tagrelease.GlobalConfig.Strategy.Increment)
	if err != nil {
		logrus.WithError(err).Fatal("strategy not supported")
	}
	formatter := tagrelease.NewFormatter(tagrelease.NewConverter(
		&tagrelease.GitAdapter{},
		strategy,
	))
	format := tagrelease.FormatFactory(formatter, tagrelease.GlobalConfig.Strategy.Format)
	output := format()

	var target io.WriteCloser
	switch tagrelease.GlobalConfig.Output {
	case "-":
		target = os.Stdout
	default:
		target, err = os.OpenFile(tagrelease.GlobalConfig.Output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.WithError(err).Fatal("cannot open target")
		}
		defer func() { _ = target.Close() }()
	}

	if tagrelease.GlobalConfig.Strategy.Escape != "" {
		output = tagrelease.EscapeSensitiveChars(output, tagrelease.GlobalConfig.Strategy.Escape)
	}

	_, _ = fmt.Fprintln(target, output)
}

func main() {
	tagrelease.DefineConfig()
	flag.Parse()
	tagrelease.LoadConfig()
	tagrelease.InitLogger()
	logrus.WithField("config", tagrelease.GlobalConfig).Debug("configured")

	rootCommand()
}
