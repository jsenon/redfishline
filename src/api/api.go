// Package api RedFishLine.
//
// the purpose of this package is to provide Web HTML Interface
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Julien SENON <julien.senon@gmail.com>

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Status Answer
// swagger:response statusjson
type statusjson struct {
	// The Status message
	// in: body
	Code    int32  `json:"statuscode"`
	Message string `json:"statusmessage"`
}

// swagger:route GET /healthy/am-i-up health amiup
//
// Check health of platform
//
// This will sent information if service is up and running.
//
//     Responses:
//       default: validationError
//       200: statusjson
func Statusamiup(w http.ResponseWriter, req *http.Request) {
	answerjson := statusjson{
		Code:    200,
		Message: "Im awake !",
	}
	b, err := json.Marshal(answerjson)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}

// swagger:route GET /healthy/about health about
//
// Check service availability
//
// This will details status availability
//
//     Responses:
//       default: validationError
//       200: statusjson
func Statusabout(w http.ResponseWriter, req *http.Request) {
	answerjson := statusjson{
		Code:    200,
		Message: "Made by Somebody",
	}
	b, err := json.Marshal(answerjson)
	if err != nil {
		fmt.Println("error:", err)
	}
	w.Write(b)
}
