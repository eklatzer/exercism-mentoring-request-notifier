package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func SetupLogging(log *logrus.Logger, logLevel, logfileName string) error {
	log.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	log.SetLevel(level)

	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err = os.Mkdir("log", 0666)
		if err != nil {
			return err
		}
	}

	logFilePath := filepath.Join("log", logfileName)

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}
