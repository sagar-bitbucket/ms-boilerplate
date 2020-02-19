package sdscale

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

// Registrar registers service instance liveness information to Consul.
type Registrar struct {
	client       Client
	registration *consul.AgentServiceRegistration
}

// NewRegistrar returns a Consul Registrar acting on the provided catalog
// registration.
func NewRegistrar(client Client, r *consul.AgentServiceRegistration) *Registrar {
	return &Registrar{
		client:       client,
		registration: r,
	}
}

// Register implements sd.Registrar interface.
func (p *Registrar) Register() {
	if err := p.client.Register(p.registration); err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println("action", "register")
	}
}

//s Deregister implements sd.Registrar interface.
func (p *Registrar) Deregister() {
	if err := p.client.Deregister(p.registration); err != nil {
		fmt.Println("err", err)
	} else {
		fmt.Println("action", "deregister")
	}
}
