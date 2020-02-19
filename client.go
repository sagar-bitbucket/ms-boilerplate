package main

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	sdscale "gitlab.com/scalent/ms-boilerplate/sdscale"
)

func main() {

	stdClient, err := consulapi.NewClient(&consulapi.Config{
		Address: "",
	})
	if err != nil {
		fmt.Println(err)
	}

	client := sdscale.NewClient(stdClient)

	// Produce a fake service registration.
	r := &consulapi.AgentServiceRegistration{
		ID:                "my-service-ID",
		Name:              "my-service-name",
		Tags:              []string{"alpha", "beta"},
		Port:              8080,
		Address:           "192.168.0.106",
		EnableTagOverride: false,
		// skipping check(s)
	}

	// Build a registrar for r.
	sdscale.NewRegistrar(client, r).Register()

	//defer registrar.Deregister()

	fmt.Println("Service API CALLED")

	// serviceURL, err := sd.LookupService("product-service")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Service URL : %s", serviceURL)

	arr, index, err := sdscale.NewDiscovery(client, "order-service", "", true).Discover()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("array", arr)
	fmt.Println("index", index)

}
