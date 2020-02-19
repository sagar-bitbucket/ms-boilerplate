package sd

import (
	"errors"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

//LookupService for Service Look Up
func LookupService(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	qo := consulapi.QueryOptions{
		Datacenter: "dc1",
	}

	services, _, err := consul.Catalog().Service(serviceName, "", &qo)

	if err != nil {
		return "", err
	}
	fmt.Println(services[0].Address, services[0].ServicePort)

	// if srvc, ok := services[serviceName]; ok {
	// 	fmt.Println(srvc)
	// 	// address := srvc.
	// 	// port := srvc.Port
	// 	// return fmt.Sprintf("http://%s:%v", address, port), nil
	// 	return "", nil
	// }
	return "", errors.New("Service Not Found")
}
