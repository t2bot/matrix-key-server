package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/federation_v1"
	"github.com/t2bot/matrix-key-server/api/health"
)

type route struct {
	method  string
	handler handler
}

func Run(listenHost string, listenPort int) {
	rtr := mux.NewRouter()

	versionHandler := handler{federation_v1.FederationVersion, "federation_version"}
	healthzHandler := handler{health.Healthz, "healthz"}

	routes := make(map[string]route)
	routes["/_matrix/federation/v1/version"] = route{"GET", versionHandler}

	for routePath, route := range routes {
		logrus.Info("Registering route: " + route.method + " " + routePath)
		rtr.Handle(routePath, route.handler).Methods(route.method)

		// This is a hack to a ensure that trailing slashes also match the routes correctly
		rtr.Handle(routePath+"/", route.handler).Methods(route.method)
	}

	rtr.Handle("/healthz", healthzHandler).Methods("OPTIONS", "GET")
	rtr.NotFoundHandler = handler{NotFoundHandler, "not_found"}
	rtr.MethodNotAllowedHandler=handler{MethodNotAllowedHandler, "method_not_allowed"}

	address := fmt.Sprintf("%s:%d", listenHost, listenPort)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", rtr)

	logrus.WithField("address", address).Info("Started up. Listening at http://" + address)
	logrus.Fatal(http.ListenAndServe(address, httpMux))
}
