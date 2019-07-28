package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/common"
)

type handler struct {
	h      func(r *http.Request, entry *logrus.Entry) interface{}
	action string
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextLog := logrus.WithFields(logrus.Fields{
		"method":   r.Method,
		"host":     r.Host,
		"resource": r.URL.Path,
	})
	contextLog.Info("Received request")

	w.Header().Set("Server", "matrix-key-server")

	// Process response
	res := h.h(r, contextLog)
	if res == nil {
		res = &common.EmptyResponse{}
	}

	contextLog.Info(fmt.Sprintf("Replying with result: %T %+v", res, res))

	statusCode := http.StatusOK
	switch result := res.(type) {
	case *common.ErrorResponse:
		statusCode = result.HttpStatus
		break
	default:
		break
	}

	// Order is important: Set headers before sending responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.Encode(res)
}
