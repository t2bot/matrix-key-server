/*
 * Copyright 2019 Travis Ralston <travis@t2bot.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os"
	"strings"

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
	var err error
	if strings.HasPrefix(*pgUrl, "/run/secrets") {
		var b []byte
		b, err = os.ReadFile(*pgUrl)
		if err != nil {
			logrus.Fatal(err)
		}
		err = db.Setup(strings.TrimSpace(string(b)))
	} else {
		err = db.Setup(*pgUrl)
	}
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
