package main

import (
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/logging"
)

func main() {
	logging.Setup()

	logrus.Info("Hello world!")
}
