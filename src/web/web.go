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
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type ILODefinition struct {
	// Server ILOHostame
	ILOHostname string `json:"ILOHostname"`
	// Server Username
	Username string `json:"Username"`
	// Server Password
	Password string `json:"Password"`
}

// Multiple Server input
var Servers []ILODefinition

// Single Server input
var Server ILODefinition

func Index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(res, Server)
}

func Send(res http.ResponseWriter, req *http.Request) {
	// Retrieve info from form
	req.ParseForm()

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")
	UEFI := req.FormValue("UEFI")
	Legacy := req.FormValue("Legacy")
	Useradd := req.FormValue("Useradd")
	PowerHigh := req.FormValue("PowerHigh")
	FastBoot := req.FormValue("FastBoot")

	fmt.Println("ILO", ILOHostname)
	fmt.Println("User", Username)
	fmt.Println("Password", Password)
	fmt.Println("Uefi", UEFI)
	fmt.Println("Legacy", Legacy)
	fmt.Println("Useradd", Useradd)
	fmt.Println("PowerHigh", PowerHigh)
	fmt.Println("FastBoot", FastBoot)

	Server.ILOHostname = ILOHostname
	Server.Username = Username
	Server.Password = Password

	// Retrieve a JSON Struct with all servers infos
	// If Json is requested reset value server ILOHostname and use []Servers Definition

	JSON := req.FormValue("JSON")

	if JSON != "" {
		fmt.Println("-> With JSON Struct")

		// Remove old value if single hostname has been used before
		Server.ILOHostname = ""
		Server.Username = ""
		Server.Password = ""
		fmt.Println("ILO", ILOHostname)
		fmt.Println("JSON", JSON)

		s := []ILODefinition{}

		err := json.Unmarshal([]byte(JSON), &s)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		// Loop on all node in config file json
		for i := range s {
			Servers = append(Servers, ILODefinition{

				ILOHostname: s[i].ILOHostname,
				Username:    s[i].Username,
				Password:    s[i].Password,
			})
		}

		fmt.Println("Servers", Servers)
		if UEFI == "on" {
			fmt.Println("------> Launch MASSIVE API UEFI")
			// Loop on Servers
		}
		if Legacy == "on" {
			fmt.Println("------> Launch MASSIVE API Legacy")
		}
		if Useradd == "on" {
			fmt.Println("------> Launch MASSIVE API Useradd")
		}
		if PowerHigh == "on" {
			fmt.Println("------> Launch MASSIVE API PowerHigh")
		}
		if FastBoot == "on" {
			fmt.Println("------> Launch MASSIVE API FastBoot")
		}

	}

	// Call to API ILO x-auth-token

	// Execute Differente API Call for ILO Configuration
	if UEFI == "on" {
		fmt.Println("------> Launch API UEFI")
	}
	if Legacy == "on" {
		fmt.Println("------> Launch API Legacy")
	}
	if Useradd == "on" {
		fmt.Println("------> Launch API Useradd")
	}
	if PowerHigh == "on" {
		fmt.Println("------> Launch API PowerHigh")
	}
	if FastBoot == "on" {
		fmt.Println("------> Launch API FastBoot")
	}
	// Execute Power Setting action

	// Get Information
	// Check if we launch at each reload or laumch on demand
	fmt.Println("------> Launch API Information")

	// http redirect to index
	http.Redirect(res, req, "/index", http.StatusSeeOther)
}

func Inventory(res http.ResponseWriter, req *http.Request) {

	// Retrieve info from form
	req.ParseForm()

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")

	fmt.Println("ILO", ILOHostname)
	fmt.Println("User", Username)
	fmt.Println("Password", Password)

	fmt.Println("------> Launch API Inventory")

	http.Redirect(res, req, "/index", http.StatusSeeOther)

}

func Rebootquick(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API RebootQuick")

}

func Reboothold(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API Reboothold")

}

func Reset(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API Reset")

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
