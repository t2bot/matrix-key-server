package api

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/t2bot/matrix-key-server/api/common"
)

func NotFoundHandler(r *http.Request, log *logrus.Entry) interface{} {
	return common.NotFoundError()
}

func MethodNotAllowedHandler(r *http.Request, log *logrus.Entry) interface{} {
	return common.MethodNotAllowed()
}