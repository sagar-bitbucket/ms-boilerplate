package main

import (
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/scalent/ms-boilerplate/sd"
)

func main() {
	portNumber := 8080
	http.HandleFunc("/", helloServer)
	http.HandleFunc("/service", serviceURL)
	http.HandleFunc("/healthcheck", healthcheck)

	sdConfig, _ := sd.DefaultConfig()
	fmt.Println(sdConfig)
	sdConfig.ServiceID = "pankaj-service-" + strconv.Itoa(portNumber)
	sdConfig.ServiceName = "pankaj-service"
	sdConfig.ServicePort = portNumber
	err := sdConfig.Register()
	fmt.Println(err)
	http.ListenAndServe(":"+strconv.Itoa(portNumber), nil)
}

//Healthcheck
func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm running fine")
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func serviceURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service API CALLED")
	serviceURL, err := sd.LookupService("aditya-service")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, "Service URL : %s", serviceURL)
}
