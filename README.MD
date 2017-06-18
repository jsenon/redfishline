# RedFishLine

Simple Web Server in GO for communication through API redfish on HP Server

Aim is to have only one console for massive update server through redfish API

### Prerequisite

You need to have:

* Go 1.8
* Go Environment properly set
* ILO Interface
* ILO Username
* ILO Password

### Compilation

```sh
git clone https://github.com/jsenon/redfishline.git
go build -o redfishline
```

### Access

Access through favorite web browser on http://YOURIP:YOURPORT ie http://127.0.0.1:9040


### API

Auto Health Check API status available at /healthy/am-i-up

### ToDo

- [ ] Add Redfish API Call
- [x] Multi Server target selection
- [x] Web part to get information details
- [ ] API Doc
- [ ] Build Docker App
- [x] Add API monitoring
- [x] Power Action
- [ ] Integration with cloudtab project
- [ ] Firmware update
