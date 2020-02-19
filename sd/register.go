package sd

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

//Register Services
func (s ServiceConfig) Register() error {
	// vailidConfig, err := s.validateConfig()

	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = s.ServiceID
	registration.Name = s.ServiceName
	registration.Address = s.ServiceAddress
	registration.Port = s.ServicePort
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck", s.ServiceAddress, s.ServicePort)
	registration.Check.Interval = s.CheckIntervel
	registration.Check.Timeout = s.CheckTimeout
	registration.Tags = s.Tags

	return consul.Agent().ServiceRegister(registration)
}

func (s ServiceConfig) validateConfig() (bool, error) {
	return true, nil
}
