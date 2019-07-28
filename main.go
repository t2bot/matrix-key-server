package main

import (
	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api"
	"github.com/t2bot/matrix-key-server/db"
	"github.com/t2bot/matrix-key-server/keys"
	"github.com/t2bot/matrix-key-server/logging"
)

func main() {
	logging.Setup()
	logrus.Info("Starting up...")

	domainName := flag.String("domain", "localhost", "The domain name for the key server")
	pgUrl := flag.String("postgres", "postgres://username:password@localhost/dbname?sslmode=disable", "PostgreSQL database URI")
	listenHost := flag.String("address", "0.0.0.0", "Address to listen for requests on")
	listenPort := flag.Int("port", 8080, "Port to listen for requests on")
	flag.Parse()

	logrus.Info("Preparing database...")
	err := db.Setup(*pgUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	keys.SelfDomainName = *domainName
	logrus.Infof("This server's domain is %s", keys.SelfDomainName)

	logrus.Info("Preparing own signing key...")
	err = prepareOwnKey()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Starting app...")
	api.Run(*listenHost, *listenPort)
}

func prepareOwnKey() error {
	key, err := keys.GetSelfKey()
	if err != nil {
		return err
	}

	logrus.Infof("Using key %s as the preferred key for this server", key.ID)
	return nil
}
