package app

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func setLogger(level, output string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
	})
	if output == "stdout" {
		logrus.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		logrus.SetOutput(file)
	}
}
