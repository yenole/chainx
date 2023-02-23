package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yenole/chainx/pkg/chainx"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}

	chainx.New(log).Run()
}
