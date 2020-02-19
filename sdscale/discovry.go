package sdscale

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type Discovery struct {
	client      Client
	service     string
	tag         string
	passingOnly bool
}

func (d Discovery) Discover() ([]string, uint64, error) {

	type response struct {
		instances []string
		index     uint64
	}

	var (
		errc = make(chan error, 1)
		resc = make(chan response, 1)
	)

	go func() {
		entries, meta, err := d.client.Service(d.service, d.tag, d.passingOnly, &consul.QueryOptions{})
		if err != nil {
			errc <- err
			return
		}

		resc <- response{
			instances: makeInstances(entries),
			index:     meta.LastIndex,
		}
	}()

	select {
	case err := <-errc:
		return nil, 0, err
	case res := <-resc:
		return res.instances, res.index, nil
		// case <-interruptc:
		// 	return nil, 0, errStopped
		// }
	}
}

func makeInstances(entries []*consul.ServiceEntry) []string {
	instances := make([]string, len(entries))
	for i, entry := range entries {
		addr := entry.Node.Address
		if entry.Service.Address != "" {
			addr = entry.Service.Address
		}
		instances[i] = fmt.Sprintf("%s:%d", addr, entry.Service.Port)
	}
	return instances
}

func NewDiscovery(client Client, service string, tag string, passingOnly bool) *Discovery {

	s := &Discovery{
		client:      client,
		service:     service,
		tag:         tag,
		passingOnly: passingOnly,
	}

	return s
}
