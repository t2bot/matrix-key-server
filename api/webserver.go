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

package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/federation_v1"
	"github.com/t2bot/matrix-key-server/api/health"
	"github.com/t2bot/matrix-key-server/api/keys_v2"
)

type route struct {
	method  string
	handler handler
}

func Run(listenHost string, listenPort int) {
	rtr := mux.NewRouter()

	healthzHandler := handler{health.Healthz, "healthz"}
	versionHandler := handler{federation_v1.FederationVersion, "federation_version"}
	localKeysHandler := handler{keys_v2.GetLocalKeys, "local_keys"}
	querySingleHandler := handler{keys_v2.QueryKeysSingle, "query_keys_single"}

	routes := make(map[string]route)
	routes["/_matrix/federation/v1/version"] = route{"GET", versionHandler}
	routes["/_matrix/key/v2/server"] = route{"GET", localKeysHandler}
	routes["/_matrix/key/v2/server/{keyId:[^/]+}"] = route{"GET", localKeysHandler}
	routes["/_matrix/key/v2/query/{serverName:[^/]+}"] = route{"GET", querySingleHandler}
	routes["/_matrix/key/v2/query/{serverName:[^/]+}/{keyId:[^/]+}"] = route{"GET", querySingleHandler}

	for routePath, route := range routes {
		logrus.Info("Registering route: " + route.method + " " + routePath)
		rtr.Handle(routePath, route.handler).Methods(route.method)

		// This is a hack to a ensure that trailing slashes also match the routes correctly
		rtr.Handle(routePath+"/", route.handler).Methods(route.method)
	}

	rtr.Handle("/healthz", healthzHandler).Methods("OPTIONS", "GET")
	rtr.NotFoundHandler = handler{NotFoundHandler, "not_found"}
	rtr.MethodNotAllowedHandler = handler{MethodNotAllowedHandler, "method_not_allowed"}

	address := fmt.Sprintf("%s:%d", listenHost, listenPort)
	httpMux := http.NewServeMux()
	httpMux.Handle("/", rtr)

	logrus.WithField("address", address).Info("Started up. Listening at http://" + address)
	logrus.Fatal(http.ListenAndServe(address, httpMux))
}
