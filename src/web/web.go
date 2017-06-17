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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func Index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(res, req)
}

func Help(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/help.html")
	t.Execute(res, req)
}

func Debug(res http.ResponseWriter, req *http.Request) {
	url := "https://ilorestfulapiexplorer.ext.hpe.com/redfish/v1/SessionService/Sessions/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"Username":"demousername","Password":"edx4qqmgeld7fu"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("ATUH:", resp.Header.Get("x-auth-token"))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}
