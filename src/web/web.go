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

//     	"url" => "redfish/v1/Systems/1/Bios/Settings/", PATCH, {"BootMode": "LegacyBios"}
//		"url" => "redfish/v1/Systems/1/Bios/Settings/", PATCH, {"PowerProfile": "MaxPerf"}
//		"url" => "redfish/v1/Systems/1/", POST, {"Action": "Reset", "ResetType": "On"}
// 		"url" => "redfish/v1/Systems/1/", POST, {"Action": "Reset", "ResetType": "ForceRestart"}

package web

import (
	"bytes"
	"crypto/tls"
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

type Credential struct {
	UserName string `json:"UserName"`
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
	// Check if no Host to revent panic
	// if ILOHostname == "" {
	// 	fmt.Println("Error Hostname")
	// 	ILOHostname = "127.0.0.1"
	// 	// http.Error(res, "Error Hostname", 500)
	// 	http.Redirect(res, req, "/index", http.StatusSeeOther)
	// 	return
	// }
	Username := req.FormValue("Username")
	// Set Default Username if not provideed, prevent panic
	if Username == "" {
		Username = "Username"
	}
	// Set Default Password if not provideed, prevent panic
	Password := req.FormValue("Password")
	if Password == "" {
		Password = "Password"
	}
	UEFI := req.FormValue("UEFI")
	Legacy := req.FormValue("Legacy")
	Useradd := req.FormValue("Useradd")
	PowerHigh := req.FormValue("PowerHigh")
	FastBoot := req.FormValue("FastBoot")

	RebootQuick := req.FormValue("RebootQuick")
	fmt.Println("Power>", RebootQuick)

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

	// Error if JSON not exist and Hostname provided
	if JSON == "" && ILOHostname == "" {
		http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}

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

		// Loop on Servers
		for i := range Servers {

			url := "https://" + Servers[i].ILOHostname + "/redfish/v1/SessionService/Sessions/"
			fmt.Println("URL:>", url)
			// Retrieve X-Auth-Token
			// Create my Body
			jsonStr := Credential{Servers[i].Username, Servers[i].Password}
			theJson, _ := json.Marshal(jsonStr)
			fmt.Println("Body:>", jsonStr)

			// Disable self certificate check
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(theJson))
			// req.Header.Set("X-Custom-Header", "")
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Transport: tr}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error Connection: ", Servers[i].ILOHostname)
				//http.Redirect(res, req, "/index", http.StatusSeeOther)
				//return
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			fmt.Println("AUTH:", resp.Header.Get("x-auth-token"))
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))

			// Retrieve x-auth-token
			token := resp.Header.Get("x-auth-token")
			fmt.Println(token)

			// url2 := "https://" + Servers[i].ILOHostname + "/redfish/v1/Systems/1/Bios/Settings/"
			// jsonStr2 := []byte(`{"BootMode":""}`)
			// req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			// client2 := &http.Client{Transport: tr}

			if Legacy == "on" && UEFI == "on" {
				fmt.Println("Error Legacy and UEFI in the same time")
			}

			if UEFI == "on" && Legacy == "" {
				fmt.Println("------> Launch MASSIVE API UEFI on", Servers[i].ILOHostname)

			}
			if Legacy == "on" && UEFI == "" {
				fmt.Println("------> Launch MASSIVE API Legacy on", Servers[i].ILOHostname)
			}
			if Useradd == "on" {
				fmt.Println("------> Launch MASSIVE API Useradd on", Servers[i].ILOHostname)
			}
			if PowerHigh == "on" {
				fmt.Println("------> Launch MASSIVE API PowerHigh on", Servers[i].ILOHostname)
			}
			if FastBoot == "on" {
				fmt.Println("------> Launch MASSIVE API FastBoot on", Servers[i].ILOHostname)
			}
		}

	} else {

		// Call to API ILO x-auth-token

		url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
		fmt.Println("URL:>", url)
		// Retrieve X-Auth-Token
		// Create my Body
		jsonStr := Credential{Username, Password}
		theJson, _ := json.Marshal(jsonStr)
		fmt.Println("Body:>", jsonStr)

		// Disable self certificate check
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(theJson))
		// req.Header.Set("X-Custom-Header", "")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: ", err)
			http.Redirect(res, req, "/index", http.StatusSeeOther)
			return
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("AUTH:", resp.Header.Get("x-auth-token"))
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

		// Retrieve x-auth-token
		token := resp.Header.Get("x-auth-token")

		url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/Bios/Settings/"
		jsonStr2 := []byte(`{"BootMode":""}`)
		req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
		client2 := &http.Client{Transport: tr}

		// Execute Differente API Call for ILO Configuration

		// Prevent UEFI and Legacy BIOS in the same time
		if Legacy == "on" && UEFI == "on" {
			fmt.Println("Error Legacy and UEFI in the same time")
		}

		if UEFI == "on" && Legacy == "" {
			fmt.Println("------> Launch API UEFI")
			// Send BIOS = UEFI
			jsonStr2 = []byte(`{"BootMode":"UEFI"}`)
			req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			fmt.Println("URL:>", url2)
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body2, _ := ioutil.ReadAll(resp2.Body)
			fmt.Println("response Status:", resp2.Status)
			fmt.Println("response Headers:", resp2.Header)
			fmt.Println("response Body:", string(body2))
		}
		if Legacy == "on" && UEFI == "" {
			fmt.Println("------> Launch API Legacy")

			// Send BIOS = Legacy
			jsonStr2 = []byte(`{"BootMode":"LegacyBios"}`)
			req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			fmt.Println("URL:>", url2)
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body3, _ := ioutil.ReadAll(resp2.Body)
			fmt.Println("response Status:", resp2.Status)
			fmt.Println("response Headers:", resp2.Header)
			fmt.Println("response Body:", string(body3))

		}

		if Useradd == "on" {
			fmt.Println("------> Launch API Useradd")
			// Need to unscope to Administrator and retrieve token instead of openstack user

		}
		if PowerHigh == "on" {
			fmt.Println("------> Launch API PowerHigh")
			// Send BIOS = MaxPerf
			jsonStr2 = []byte(`{"PowerProfile":"MaxPerf"}`)
			req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			fmt.Println("URL:>", url2)
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body4, _ := ioutil.ReadAll(resp2.Body)
			fmt.Println("response Status:", resp2.Status)
			fmt.Println("response Headers:", resp2.Header)
			fmt.Println("response Body:", string(body4))

		}
		if FastBoot == "on" {
			fmt.Println("------> Launch API FastBoot")
			// Send BIOS = Extended Memory Test Off
			jsonStr2 = []byte(`{"ExtendedMemTest":"Disabled"}`)
			req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			fmt.Println("URL:>", url2)
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body5, _ := ioutil.ReadAll(resp2.Body)
			fmt.Println("response Status:", resp2.Status)
			fmt.Println("response Headers:", resp2.Header)
			fmt.Println("response Body:", string(body5))
		}

	}

	// Execute Power Setting action

	// Get Information
	// Check if we launch at each reload or launch on demand
	//fmt.Println("------> Launch API Information")

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
	req.ParseForm()
	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API RebootQuick")
	req.ParseForm()

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")

	fmt.Println("ILO", ILOHostname)
	fmt.Println("User", Username)
	fmt.Println("Password", Password)

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
	url := "https:///redfish/v1/SessionService/Sessions/"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"UserName":"","Password":""}`)
	fmt.Println("Body:>", jsonStr)

	// Disable self certificate check
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	// req.Header.Set("X-Custom-Header", "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("AUTH:", resp.Header.Get("x-auth-token"))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	// Retrieve x-auth-token
	token := resp.Header.Get("x-auth-token")

	// New Session
	url2 := "https:///redfish/v1/Chassis/1/"
	req2, err := http.NewRequest("GET", url2, nil)
	req2.Header.Set("X-Auth-Token", token)
	fmt.Println("URL:>", url2)
	client2 := &http.Client{Transport: tr}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()
	fmt.Println("response Status:", resp2.Status)
	fmt.Println("response Headers:", resp2.Header)
	body2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println("response Body:", string(body2))
	var data map[string]interface{}
	var data2 map[string]map[string]interface{}
	erro := json.Unmarshal([]byte(body2), &data)
	if erro != nil {
		panic(err)
	}
	fmt.Println("Model>", data["Model"])
	erro2 := json.Unmarshal([]byte(body2), &data2)
	if erro2 != nil {
	}
	fmt.Println("Health>", data2["Status"]["Health"])

}
