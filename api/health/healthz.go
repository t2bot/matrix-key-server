package health

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type HealthzResponse struct {
	OK     bool   `json:"ok"`
	Status string `json:"status"`
}

func Healthz(r *http.Request, log *logrus.Entry) interface{} {
	return &HealthzResponse{
		OK:     true,
		Status: "Probably not dead",
	}
}
