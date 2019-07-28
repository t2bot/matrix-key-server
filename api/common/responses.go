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

package common

import (
	"net/http"
)

type EmptyResponse struct{}

type ErrorResponse struct {
	Code       string `json:"errcode"`
	Message    string `json:"error"`
	HttpStatus int    `json:"http_status"`
}

func InternalServerError(message string) *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", message, http.StatusInternalServerError}
}

func MethodNotAllowed() *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", "Method Not Allowed", http.StatusMethodNotAllowed}
}

func NotFoundError() *ErrorResponse {
	return &ErrorResponse{"M_NOT_FOUND", "Resource Not Found", http.StatusNotFound}
}

func UnauthorizedError() *ErrorResponse {
	return &ErrorResponse{"M_UNAUTHORIZED", "Authentication Failed", http.StatusUnauthorized}
}

func BadRequest(message string) *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", message, http.StatusBadRequest}
}
