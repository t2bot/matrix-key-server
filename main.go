package main

import (
	"github.com/namsral/flag"
	"github.com/t2bot/matrix-key-server/api"
	"github.com/t2bot/matrix-key-server/logging"
)

func main() {
	logging.Setup()

	//domainName := flag.String("domain", "localhost", "The domain name for the key server")
	//pgUrl := flag.String("postgres", "postgres://username:password@localhost/dbname?sslmode=disable", "PostgreSQL database URI")
	listenHost := flag.String("address", "0.0.0.0", "Address to listen for requests on")
	listenPort := flag.Int("port", 8080, "Port to listen for requests on")
	flag.Parse()

	api.Run(*listenHost, *listenPort)
}
