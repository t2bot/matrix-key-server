package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Setup() {
	logrus.SetFormatter(&utcFormatter{
		&logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.000 Z07:00",
			FullTimestamp:    true,
			ForceColors:      true,
			DisableColors:    false,
			DisableTimestamp: false,
			QuoteEmptyFields: true,
		},
	})
	logrus.SetOutput(os.Stdout)
}
