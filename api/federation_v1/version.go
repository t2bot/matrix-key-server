package federation_v1

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/meta"
)

type FederationServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type FederationVersionResponse struct {
	Server FederationServerInfo `json:"server"`
}

func FederationVersion(r *http.Request, log *logrus.Entry) interface{} {
	return &FederationVersionResponse{
		Server: FederationServerInfo{
			Name:    "matrix-key-server",
			Version: meta.AppVersion,
		},
	}
}
