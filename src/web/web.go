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

package web

import (
	"html/template"
	"net/http"
)

func Index(res http.ResponseWriter, req *http.Request) {

	rs, err := db.GetInfo()

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(res, rs)
	if err != nil {
		return
	}
}

func Help(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/help.html")

	t.Execute(res, req)

}
