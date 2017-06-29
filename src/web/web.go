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

//      Inventory Mac /redfish/v1/Systems/1/NetworkAdapters/
//      PowerState /redfish/v1/Systems/1/
//      Raid /redfish/v1/Systems/1/SmartStorage/ https://10.67.224.5/redfish/v1/Systems/1/SmartStorage/ArrayControllers/0/LogicalDrives/
//
//
//      With Help of @Merrick28

package web

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
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

type InventoryServer struct {
	Hostname     string         `json:"Hostname"`
	Memory       float64        `json:"Memory"`
	CPUNum       float64        `json:"CPUNum"`
	CPUModel     string         `json:"CPUModel"`
	Model        string         `json:"Model"`
	SerialNumber string         `json:"SerialNumber"`
	Health       string         `json:"Health"`
	Power        string         `json:"Power"`
	PowerState   string         `json:"PowerState"`
	MacInfo      []InventoryMac `json:"MacInfo"`
}

type InventoryMac struct {
	CardName string   `json:"CardName"`
	Mac      []string `json:"Mac"`
	Position []string `json:"Position"`
}

// Structure to create an account
type AccountILO struct {
	UserName  string `json:"UserName"`
	Password  string `json:"Password"`
	LoginName string `json:"LoginName"`
	Oem       struct {
		Hp struct {
			LoginName  string `json:"LoginName"`
			Privileges struct {
				RemoteConsolePriv        bool `json:"RemoteConsolePriv"`
				VirtualMediaPriv         bool `json:"VirtualMediaPriv"`
				UserConfigPriv           bool `json:"UserConfigPriv"`
				iLOConfigPriv            bool `json:"iLOConfigPriv"`
				VirtualPowerAndResetPriv bool `json:"VirtualPowerAndResetPriv"`
			} `json:"Privileges"`
		} `json:"Hp"`
	} `json:"Oem"`
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
	AllowReset := req.FormValue("AllowReset")

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

		// Loop on Servers
		for i := range Servers {

			url := "https://" + Servers[i].ILOHostname + "/redfish/v1/SessionService/Sessions/"
			// Retrieve X-Auth-Token
			// Create my Body
			jsonStr := Credential{Servers[i].Username, Servers[i].Password}
			theJson, _ := json.Marshal(jsonStr)

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

			// body, _ := ioutil.ReadAll(resp.Body)

			// Retrieve x-auth-token
			token := resp.Header.Get("x-auth-token")
			session := resp.Header.Get("location")

			url2 := "https://" + Servers[i].ILOHostname + "/redfish/v1/Systems/1/Bios/Settings/"
			jsonStr2 := []byte(`{"BootMode":""}`)
			req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
			client2 := &http.Client{Transport: tr}

			if Legacy == "on" && UEFI == "on" {
				fmt.Println("Error Legacy and UEFI in the same time")
			}

			if UEFI == "on" && Legacy == "" {
				fmt.Println("------> Launch MASSIVE API UEFI on", Servers[i].ILOHostname)
				// Send BIOS = UEFI
				jsonStr2 = []byte(`{"BootMode":"UEFI"}`)
				req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
				if err2 != nil {
					http.Redirect(res, req, "/index", http.StatusSeeOther)
					return
				}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				_, err2 := client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				// body2, _ := ioutil.ReadAll(resp2.Body)

			}
			if Legacy == "on" && UEFI == "" {
				fmt.Println("------> Launch MASSIVE API Legacy on", Servers[i].ILOHostname)
				// Send BIOS = Legacy
				jsonStr2 = []byte(`{"BootMode":"LegacyBios"}`)
				req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				_, err2 := client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				// body3, _ := ioutil.ReadAll(resp2.Body)

			}
			if Useradd == "on" {
				fmt.Println("------> Launch MASSIVE API Useradd on", Servers[i].ILOHostname)
				err := AddUser(token, Servers[i].ILOHostname)
				if err != nil {
					fmt.Println("Error Creation User")
				}

			}
			if PowerHigh == "on" {
				fmt.Println("------> Launch MASSIVE API PowerHigh on", Servers[i].ILOHostname)
				// Send BIOS = MaxPerf
				jsonStr2 = []byte(`{"PowerProfile":"MaxPerf"}`)
				req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")

				_, err2 := client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				// body4, _ := ioutil.ReadAll(resp2.Body)

			}
			if FastBoot == "on" {
				fmt.Println("------> Launch MASSIVE API FastBoot on", Servers[i].ILOHostname)
				fmt.Println("------> Launch API FastBoot")
				// Send BIOS = Extended Memory Test Off
				jsonStr2 = []byte(`{"ExtendedMemTest":"Disabled"}`)
				req2, err2 = http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				_, err2 := client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error Connection: ", Servers[i].ILOHostname)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					// return
				}
				// body5, _ := ioutil.ReadAll(resp2.Body)

			}
			// Perform A server reset if checked
			if AllowReset == "on" {
				fmt.Println("------> Launch API Apply Setting by reseting server")

				// We need to check status in order to launch reset

				fmt.Println("------> Launch API Check Power Status")

				url2 := "https://" + Servers[i].ILOHostname + "/redfish/v1/Systems/1/"

				req2, err2 := http.NewRequest("GET", url2, bytes.NewBuffer(jsonStr2))
				client2 := &http.Client{Transport: tr}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				resp2, err2 := client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error: ", err)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					return
				}
				body6, _ := ioutil.ReadAll(resp2.Body)

				var data map[string]interface{}
				json.Unmarshal([]byte(body6), &data)
				state := data["PowerState"]

				// if off we start system
				if state == "Off" {
					jsonStr2 = []byte(`{"Action": "Reset", "ResetType": "On"}`)
					req2, err2 = http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
					client2 = &http.Client{Transport: tr}
					req2.Header.Set("X-Auth-Token", token)
					req2.Header.Set("Content-Type", "application/json")
					resp2, err2 = client2.Do(req2)
					if err2 != nil {
						fmt.Println("Error: ", err)
						// http.Redirect(res, req, "/index", http.StatusSeeOther)
						return
					}
					// body7, _ := ioutil.ReadAll(resp2.Body)

				} else {

					// if on we reset system
					jsonStr2 = []byte(`{"Action": "Reset", "ResetType": "ForceRestart"}`)
					req2, err2 = http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
					client2 = &http.Client{Transport: tr}
					req2.Header.Set("X-Auth-Token", token)
					req2.Header.Set("Content-Type", "application/json")
					resp2, err2 = client2.Do(req2)
					if err2 != nil {
						fmt.Println("Error: ", err)
						// http.Redirect(res, req, "/index", http.StatusSeeOther)
						return
					}
					// body8, _ := ioutil.ReadAll(resp2.Body)

				}

			}
			// Close session
			req5, err := http.NewRequest("DELETE", session, nil)
			req5.Header.Set("X-Auth-Token", token)
			resp5, err5 := client2.Do(req5)
			if err5 != nil {
				panic(err5)
			}
			defer resp5.Body.Close()
			defer client2.Do(req5)
		}

	} else {

		// Call to API ILO x-auth-token

		url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
		// Retrieve X-Auth-Token
		// Create my Body
		jsonStr := Credential{Username, Password}
		theJson, _ := json.Marshal(jsonStr)

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

		// body, _ := ioutil.ReadAll(resp.Body)

		// Retrieve x-auth-token
		token := resp.Header.Get("x-auth-token")
		session := resp.Header.Get("location")

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
			_, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			// body2, _ := ioutil.ReadAll(resp2.Body)

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
			_, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			// body3, _ := ioutil.ReadAll(resp2.Body)

		}

		if Useradd == "on" {
			fmt.Println("------> Launch API Useradd")
			// Need to unscope to Administrator and retrieve token instead of openstack user
			err := AddUser(token, ILOHostname)
			if err != nil {
				fmt.Println("Error Creation User")
			}

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
			_, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			// body4, _ := ioutil.ReadAll(resp2.Body)

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
			_, err2 := client2.Do(req2)
			if err2 != nil {
				http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			// body5, _ := ioutil.ReadAll(resp2.Body)

		}
		// Perform A server reset if checked
		if AllowReset == "on" {
			fmt.Println("------> Launch API Apply Setting by reseting server")

			// We need to check status in order to launch reset

			fmt.Println("------> Launch API Check Power Status")

			url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/"

			req2, err2 := http.NewRequest("GET", url2, bytes.NewBuffer(jsonStr2))
			client2 := &http.Client{Transport: tr}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				fmt.Println("Error: ", err)
				// http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body6, _ := ioutil.ReadAll(resp2.Body)

			var data map[string]interface{}
			json.Unmarshal([]byte(body6), &data)
			state := data["PowerState"]

			// if off we start system
			if state == "Off" {
				jsonStr2 = []byte(`{"Action": "Reset", "ResetType": "On"}`)
				req2, err2 = http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
				client2 = &http.Client{Transport: tr}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				resp2, err2 = client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error: ", err)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					return
				}
				// body7, _ := ioutil.ReadAll(resp2.Body)

			} else {

				// if on we reset system
				jsonStr2 = []byte(`{"Action": "Reset", "ResetType": "ForceRestart"}`)
				req2, err2 = http.NewRequest("POST", url2, bytes.NewBuffer(jsonStr2))
				client2 = &http.Client{Transport: tr}
				req2.Header.Set("X-Auth-Token", token)
				req2.Header.Set("Content-Type", "application/json")
				resp2, err2 = client2.Do(req2)
				if err2 != nil {
					fmt.Println("Error: ", err)
					// http.Redirect(res, req, "/index", http.StatusSeeOther)
					return
				}
				// body8, _ := ioutil.ReadAll(resp2.Body)

			}

		}
		// Close session
		req5, err := http.NewRequest("DELETE", session, nil)
		req5.Header.Set("X-Auth-Token", token)
		resp5, err5 := client2.Do(req5)
		if err5 != nil {
			panic(err5)
		}
		defer resp5.Body.Close()
		defer client2.Do(req5)

	}

	// Get Information
	// Check if we launch at each reload or launch on demand
	//fmt.Println("------> Launch API Information")

	// http redirect to index
	http.Redirect(res, req, "/index", http.StatusSeeOther)
}

func Inventory(res http.ResponseWriter, req *http.Request) {

	// Retrieve info from form
	// http://127.0.0.1/redfish/v1/Systems/1/

	var data map[string]map[string]interface{}
	var data2 map[string]interface{}
	var data3 map[string]interface{}
	var data4 map[string][]map[string]interface{}

	myinventory := []InventoryServer{}

	req.ParseForm()

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")
	JSON := req.FormValue("JSON")
	fmt.Println("JSON>", JSON)

	// Error if JSON not exist and Hostname provided
	if JSON == "" && ILOHostname == "" {
		http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}

	// Massive ILO Credential
	if JSON != "" {
		fmt.Println("------> Launch Massive API Inventory")

		// Remove old value if single hostname has been used before
		Server.ILOHostname = ""
		Server.Username = ""
		Server.Password = ""

		s := []ILODefinition{}

		err := json.Unmarshal([]byte(JSON), &s)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		fmt.Println("Unmarshal>", s)

		// Loop on all node in config file json and store information in Server struct ILODefinition

		// Loop on s
		for i := range s {

			fmt.Println("ILOHostname: ", s[i].ILOHostname)

			url := "https://" + s[i].ILOHostname + "/redfish/v1/SessionService/Sessions/"
			// Retrieve X-Auth-Token
			// Create my Body
			jsonStr := Credential{s[i].Username, s[i].Password}

			fmt.Println("jsonstr: ", jsonStr)

			theJson, _ := json.Marshal(jsonStr)

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

			// body, _ := ioutil.ReadAll(resp.Body)

			// Retrieve x-auth-token
			token := resp.Header.Get("x-auth-token")
			session := resp.Header.Get("location")

			url2 := "https://" + s[i].ILOHostname + "/redfish/v1/Systems/1/"
			req2, err2 := http.NewRequest("GET", url2, nil)
			client2 := &http.Client{Transport: tr}
			req2.Header.Set("X-Auth-Token", token)
			req2.Header.Set("Content-Type", "application/json")
			resp2, err2 := client2.Do(req2)
			if err2 != nil {
				fmt.Println("Error URL Creation: ", err)
				// http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body9, _ := ioutil.ReadAll(resp2.Body)

			// BIOS

			url3 := "https://" + s[i].ILOHostname + "/redfish/v1/Systems/1/Bios/"
			req3, err3 := http.NewRequest("GET", url3, nil)
			client3 := &http.Client{Transport: tr}
			req3.Header.Set("X-Auth-Token", token)
			req3.Header.Set("Content-Type", "application/json")
			resp3, err3 := client3.Do(req3)
			if err3 != nil {
				fmt.Println("Error URL Creation BIOS: ", err)
				// http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body10, _ := ioutil.ReadAll(resp3.Body)

			url4 := "https://" + s[i].ILOHostname + "/redfish/v1/Managers/1/EthernetInterfaces/"
			req4, err4 := http.NewRequest("GET", url4, nil)
			client4 := &http.Client{Transport: tr}
			req4.Header.Set("X-Auth-Token", token)
			req4.Header.Set("Content-Type", "application/json")
			resp4, err4 := client4.Do(req4)
			if err4 != nil {
				fmt.Println("Error: ", err)
				// http.Redirect(res, req, "/index", http.StatusSeeOther)
				return
			}
			body11, _ := ioutil.ReadAll(resp4.Body)

			json.Unmarshal([]byte(body9), &data)
			json.Unmarshal([]byte(body9), &data2)
			json.Unmarshal([]byte(body10), &data3)
			json.Unmarshal([]byte(body11), &data4)

			// HTML Rendering

			// tempmem := data["Memory"]["TotalSystemMemoryGB"].(float64)

			var Ethernet []InventoryMac

			Ethernet, err5 := RetrieveMacAddress(token, ILOHostname)

			fmt.Println("err5", err5)
			fmt.Println("Ethernet", Ethernet)

			myinventory = append(myinventory, InventoryServer{

				Hostname:     data4["Items"][0]["FQDN"].(string),
				Memory:       data["Memory"]["TotalSystemMemoryGB"].(float64),
				CPUNum:       data["Processors"]["Count"].(float64),
				CPUModel:     data["Processors"]["ProcessorFamily"].(string),
				Model:        data2["Model"].(string),
				SerialNumber: data2["SerialNumber"].(string),
				Health:       data["Status"]["Health"].(string),
				Power:        data3["PowerRegulator"].(string),
				PowerState:   data2["PowerState"].(string),
				MacInfo:      Ethernet,
			})
			// Close session
			req5, err := http.NewRequest("DELETE", session, nil)
			req5.Header.Set("X-Auth-Token", token)
			resp5, err5 := client2.Do(req5)
			if err5 != nil {
				panic(err5)
			}
			defer resp5.Body.Close()
			defer client2.Do(req5)
		}
		fmt.Println("myinventory: ", myinventory)
		req.ParseForm()
		t, _ := template.ParseFiles("templates/inventory.html")
		t.Execute(res, myinventory)

	} else {

		// Single ILO Credential

		fmt.Println("------> Launch API Inventory")

		url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
		// Retrieve X-Auth-Token
		// Create my Body
		jsonStr := Credential{Username, Password}
		theJson, _ := json.Marshal(jsonStr)

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

		// body, _ := ioutil.ReadAll(resp.Body)

		// Retrieve x-auth-token
		token := resp.Header.Get("x-auth-token")
		session := resp.Header.Get("location")

		url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/"
		req2, err2 := http.NewRequest("GET", url2, nil)
		client2 := &http.Client{Transport: tr}
		req2.Header.Set("X-Auth-Token", token)
		req2.Header.Set("Content-Type", "application/json")
		resp2, err2 := client2.Do(req2)
		if err2 != nil {
			fmt.Println("Error: ", err)
			// http.Redirect(res, req, "/index", http.StatusSeeOther)
			return
		}
		body9, _ := ioutil.ReadAll(resp2.Body)

		// BIOS

		url3 := "https://" + ILOHostname + "/redfish/v1/Systems/1/Bios/"
		req3, err3 := http.NewRequest("GET", url3, nil)
		client3 := &http.Client{Transport: tr}
		req3.Header.Set("X-Auth-Token", token)
		req3.Header.Set("Content-Type", "application/json")
		resp3, err3 := client3.Do(req3)
		if err3 != nil {
			fmt.Println("Error: ", err)
			// http.Redirect(res, req, "/index", http.StatusSeeOther)
			return
		}
		body10, _ := ioutil.ReadAll(resp3.Body)

		// Retrieve Name Server ILO

		url4 := "https://" + ILOHostname + "/redfish/v1/Managers/1/EthernetInterfaces/"
		req4, err4 := http.NewRequest("GET", url4, nil)
		client4 := &http.Client{Transport: tr}
		req4.Header.Set("X-Auth-Token", token)
		req4.Header.Set("Content-Type", "application/json")
		resp4, err4 := client4.Do(req4)
		if err4 != nil {
			fmt.Println("Error: ", err)
			// http.Redirect(res, req, "/index", http.StatusSeeOther)
			return
		}
		body11, _ := ioutil.ReadAll(resp4.Body)

		json.Unmarshal([]byte(body9), &data)
		json.Unmarshal([]byte(body9), &data2)
		json.Unmarshal([]byte(body10), &data3)
		json.Unmarshal([]byte(body11), &data4)

		// HTML Rendering

		// tempmem := data["Memory"]["TotalSystemMemoryGB"].(float64)

		var Ethernet []InventoryMac

		Ethernet, err5 := RetrieveMacAddress(token, ILOHostname)

		fmt.Println("err5", err5)
		fmt.Println("Ethernet", Ethernet)

		myinventory = append(myinventory, InventoryServer{

			Hostname:     data4["Items"][0]["FQDN"].(string),
			Memory:       data["Memory"]["TotalSystemMemoryGB"].(float64),
			CPUNum:       data["Processors"]["Count"].(float64),
			CPUModel:     data["Processors"]["ProcessorFamily"].(string),
			Model:        data2["Model"].(string),
			SerialNumber: data2["SerialNumber"].(string),
			Health:       data["Status"]["Health"].(string),
			Power:        data3["PowerRegulator"].(string),
			PowerState:   data2["PowerState"].(string),
			MacInfo:      Ethernet,
		})

		req.ParseForm()
		t, _ := template.ParseFiles("templates/inventory.html")
		t.Execute(res, myinventory)

		// Close session
		req5, err := http.NewRequest("DELETE", session, nil)
		req5.Header.Set("X-Auth-Token", token)
		resp5, err5 := client2.Do(req5)
		if err5 != nil {
			panic(err5)
		}
		defer resp5.Body.Close()
		defer client2.Do(req5)

	}

}

func Rebootquick(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API RebootQuick")

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")

	// Call to API ILO x-auth-token

	url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
	// Retrieve X-Auth-Token
	// Create my Body
	jsonStr := Credential{Username, Password}
	theJson, _ := json.Marshal(jsonStr)

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
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)

	// Retrieve x-auth-token
	token := resp.Header.Get("x-auth-token")
	session := resp.Header.Get("location")

	url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/"
	jsonStr2 := []byte(`{"Action": "PowerButton", "PushType": "Press", "Target": "/Oem/Hp"}`)
	req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
	client2 := &http.Client{Transport: tr}
	req2.Header.Set("X-Auth-Token", token)
	req2.Header.Set("Content-Type", "application/json")
	_, err2 = client2.Do(req2)
	if err2 != nil {
		fmt.Println("Error: ", err)
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}

	// Close session
	req3, err := http.NewRequest("DELETE", session, nil)
	req3.Header.Set("X-Auth-Token", token)
	resp3, err3 := client2.Do(req3)
	if err3 != nil {
		panic(err3)
	}
	defer resp3.Body.Close()
	defer client2.Do(req3)

}

func Reboothold(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API Reboothold")

	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")

	// Call to API ILO x-auth-token

	url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
	// Retrieve X-Auth-Token
	// Create my Body
	jsonStr := Credential{Username, Password}
	theJson, _ := json.Marshal(jsonStr)

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
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)

	// Retrieve x-auth-token
	token := resp.Header.Get("x-auth-token")
	session := resp.Header.Get("location")

	url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/"
	jsonStr2 := []byte(`{"Action": "PowerButton", "PushType": "PressAndHold", "Target": "/Oem/Hp"}`)
	req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
	client2 := &http.Client{Transport: tr}
	req2.Header.Set("X-Auth-Token", token)
	req2.Header.Set("Content-Type", "application/json")
	_, err2 = client2.Do(req2)
	if err2 != nil {
		fmt.Println("Error: ", err)
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}
	// Close session
	req3, err := http.NewRequest("DELETE", session, nil)
	req3.Header.Set("X-Auth-Token", token)
	resp3, err3 := client2.Do(req3)
	if err3 != nil {
		panic(err3)
	}
	defer resp3.Body.Close()
	defer client2.Do(req3)
}

func Reset(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	t, _ := template.ParseFiles("templates/Reboot.html")
	t.Execute(res, Server)
	fmt.Println("------> Launch API Reset")
	ILOHostname := req.FormValue("ILOHostname")
	Username := req.FormValue("Username")
	Password := req.FormValue("Password")

	// Call to API ILO x-auth-token

	url := "https://" + ILOHostname + "/redfish/v1/SessionService/Sessions/"
	// Retrieve X-Auth-Token
	// Create my Body
	jsonStr := Credential{Username, Password}
	theJson, _ := json.Marshal(jsonStr)

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
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	// body, _ := ioutil.ReadAll(resp.Body)

	// Retrieve x-auth-token
	token := resp.Header.Get("x-auth-token")
	session := resp.Header.Get("location")

	url2 := "https://" + ILOHostname + "/redfish/v1/Systems/1/"
	jsonStr2 := []byte(`{"Action": "Reset", "ResetType": "ForceRestart"}`)
	req2, err2 := http.NewRequest("PATCH", url2, bytes.NewBuffer(jsonStr2))
	client2 := &http.Client{Transport: tr}
	req2.Header.Set("X-Auth-Token", token)
	req2.Header.Set("Content-Type", "application/json")
	_, err2 = client2.Do(req2)
	if err2 != nil {
		fmt.Println("Error: ", err)
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return
	}
	// Close session
	req3, err := http.NewRequest("DELETE", session, nil)
	req3.Header.Set("X-Auth-Token", token)
	resp3, err3 := client2.Do(req3)
	if err3 != nil {
		panic(err3)
	}
	defer resp3.Body.Close()
	defer client2.Do(req3)
}

func Help(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/help.html")
	t.Execute(res, req)
}

func Debug(res http.ResponseWriter, req *http.Request) {

	url := "https:///redfish/v1/SessionService/Sessions/"
	// fmt.Println("URL:>", url)

	var EthernetBis []InventoryMac
	EthernetBis = make([]InventoryMac, 0, 10)

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
	session := resp.Header.Get("location")

	// New Session
	url2 := "https:///redfish/v1/Systems/1/NetworkAdapters"
	req2, err := http.NewRequest("GET", url2, nil)
	req2.Header.Set("X-Auth-Token", token)
	// fmt.Println("URL:>", url2)
	client2 := &http.Client{Transport: tr}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()
	body2, _ := ioutil.ReadAll(resp2.Body)

	// Create Interface for parsing
	var f interface{}

	// Unmarshal json with interface
	erro2 := json.Unmarshal([]byte(body2), &f)
	if erro2 != nil {
		fmt.Println("Error: ", erro2)
		panic(erro2)
	}

	// Go to definition needed. Don t use .(string) assertion at the end
	// fmt.Println("f:", f.(map[string]interface{})["links"].(map[string]interface{})["Member"].([]interface{})[1].(map[string]interface{})["href"])

	// Store Value of links/Member
	l := f.(map[string]interface{})["links"].(map[string]interface{})["Member"]

	// Use reflect to store value of interface
	// ValueOf returns a new Value initialized to the concrete value stored in the interface l
	s := reflect.ValueOf(l)

	// Range over all Member value
	// Loop over Nbr of Card

	for j := 0; j < s.Len(); j++ {
		urlstring := reflect.ValueOf(f.(map[string]interface{})["links"].(map[string]interface{})["Member"].([]interface{})[j].(map[string]interface{})["href"]).String()

		url3 := "https://" + urlstring
		req3, err3 := http.NewRequest("GET", url3, nil)
		if err3 != nil {
			panic(err3)
		}
		req3.Header.Set("X-Auth-Token", token)
		fmt.Println("URL:>", url3)
		client2 := &http.Client{Transport: tr}
		resp3, err4 := client2.Do(req3)
		if err4 != nil {
			panic(err4)
		}
		defer resp3.Body.Close()
		body3, _ := ioutil.ReadAll(resp3.Body)

		var g interface{}

		erro2 := json.Unmarshal([]byte(body3), &g)
		if erro2 != nil {
			fmt.Println("Error: ", erro2)
			panic(erro2)
		}

		fmt.Println("Name Card:", g.(map[string]interface{})["Name"])

		EthernetBis = append(EthernetBis, InventoryMac{CardName: reflect.ValueOf(g.(map[string]interface{})["Name"]).String()})

		nbrports := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"])

		// Loop Over Nbre of Port in that Card

		for i := 0; i < nbrports.Len(); i++ {

			toto := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"].([]interface{})[i].(map[string]interface{})["MacAddress"]).String()
			fmt.Println("toto>", toto)
			EthernetBis[j].Mac = append(EthernetBis[j].Mac, toto)

			titi := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"].([]interface{})[i].(map[string]interface{})["Oem"].(map[string]interface{})["Hp"].(map[string]interface{})["StructuredName"]).String()
			fmt.Println("titi>", titi)
			EthernetBis[j].Position = append(EthernetBis[j].Position, titi)

		}
		fmt.Println("Slice", EthernetBis)

	}

	// type InventoryMac struct {
	// 	CardName string   `json:"CardName"`
	// 	Mac      []string `json:"Mac"`
	// 	Position []string `json:"Position"`
	// }

	fmt.Println("Ethernet>", EthernetBis)

	// Close session
	req3, err := http.NewRequest("DELETE", session, nil)
	req3.Header.Set("X-Auth-Token", token)
	resp3, err3 := client2.Do(req3)
	if err3 != nil {
		panic(err3)
	}
	defer resp3.Body.Close()
	defer client2.Do(req3)

}

func Serialize(res http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/debug.html")
	t.Execute(res, req)
}

func SerializeSend(res http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	MyHost := req.FormValue("MyHost")
	MyUser := req.FormValue("MyUser")
	MyPassword := req.FormValue("MyPassword")

	fmt.Println("MyHost> ", MyHost)
	fmt.Println("MyUser> ", MyUser)
	fmt.Println("MyPassword> ", MyPassword)

	jsonStr := ILODefinition{MyUser, MyPassword, MyHost}
	// theJson, _ := json.MarshalIndent(jsonStr, "", "    ")
	theJson, _ := json.Marshal(jsonStr)

	fmt.Println("theJson> ", string(theJson))

	req.ParseForm()
	t, _ := template.ParseFiles("templates/result.html")
	t.Execute(res, string(theJson))

}

func AddUser(token string, hostname string) error {

	fmt.Println("------> Launch API Add User")
	// Launch API adding user

	url := "https://" + hostname + "/redfish/v1/AccountService/Accounts/"

	// Read from file
	sheetData, err := ioutil.ReadFile("../credential-ilo-airbus.json")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(sheetData))
	if err != nil {
		fmt.Println("Error API")
	}
	// Disable self certificate check
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req.Header.Set("X-Auth-Token", token)
	req.Header.Set("Content-Type", "application/json")
	_, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println("Error: ", err2)
		// http.Redirect(res, req, "/index", http.StatusSeeOther)
		return err2
	}
	// body9, _ := ioutil.ReadAll(resp.Body)
	return nil
}

func RetrieveMacAddress(token string, hostname string) ([]InventoryMac, error) {
	fmt.Println("------> Launch API MAC Address")

	var EthernetBis []InventoryMac
	EthernetBis = make([]InventoryMac, 0, 10)

	url2 := "https://" + hostname + "/redfish/v1/Systems/1/NetworkAdapters"
	req2, _ := http.NewRequest("GET", url2, nil)
	req2.Header.Set("X-Auth-Token", token)
	// fmt.Println("URL:>", url2)
	// Disable self certificate check
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client2 := &http.Client{Transport: tr}
	resp2, err2 := client2.Do(req2)
	if err2 != nil {
		panic(err2)
	}
	defer resp2.Body.Close()
	body2, _ := ioutil.ReadAll(resp2.Body)

	// Create Interface for parsing
	var f interface{}

	// Unmarshal json with interface
	erro2 := json.Unmarshal([]byte(body2), &f)
	if erro2 != nil {
		fmt.Println("Error: ", erro2)
		panic(erro2)
	}

	// Go to definition needed. Don t use .(string) assertion at the end
	// fmt.Println("f:", f.(map[string]interface{})["links"].(map[string]interface{})["Member"].([]interface{})[1].(map[string]interface{})["href"])

	// Store Value of links/Member
	l := f.(map[string]interface{})["links"].(map[string]interface{})["Member"]

	// Use reflect to store value of interface
	// ValueOf returns a new Value initialized to the concrete value stored in the interface l
	s := reflect.ValueOf(l)

	// Range over all Member value
	// Loop over Nbr of Card

	for j := 0; j < s.Len(); j++ {
		urlstring := reflect.ValueOf(f.(map[string]interface{})["links"].(map[string]interface{})["Member"].([]interface{})[j].(map[string]interface{})["href"]).String()

		url3 := "https://10.67.224.23" + urlstring
		req3, err3 := http.NewRequest("GET", url3, nil)
		if err3 != nil {
			panic(err3)
		}
		req3.Header.Set("X-Auth-Token", token)
		fmt.Println("URL:>", url3)
		client2 := &http.Client{Transport: tr}
		resp3, err4 := client2.Do(req3)
		if err4 != nil {
			panic(err4)
		}
		defer resp3.Body.Close()
		body3, _ := ioutil.ReadAll(resp3.Body)

		var g interface{}

		erro2 := json.Unmarshal([]byte(body3), &g)
		if erro2 != nil {
			fmt.Println("Error: ", erro2)
			panic(erro2)
		}

		fmt.Println("Name Card:", g.(map[string]interface{})["Name"])

		EthernetBis = append(EthernetBis, InventoryMac{CardName: reflect.ValueOf(g.(map[string]interface{})["Name"]).String()})

		nbrports := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"])

		// Loop Over Nbre of Port in that Card

		for i := 0; i < nbrports.Len(); i++ {

			toto := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"].([]interface{})[i].(map[string]interface{})["MacAddress"]).String()
			fmt.Println("toto>", toto)
			EthernetBis[j].Mac = append(EthernetBis[j].Mac, toto)

			titi := reflect.ValueOf(g.(map[string]interface{})["PhysicalPorts"].([]interface{})[i].(map[string]interface{})["Oem"].(map[string]interface{})["Hp"].(map[string]interface{})["StructuredName"]).String()
			fmt.Println("titi>", titi)
			EthernetBis[j].Position = append(EthernetBis[j].Position, titi)

		}

	}

	fmt.Println("Ethernet>", EthernetBis)

	return EthernetBis, nil
}
