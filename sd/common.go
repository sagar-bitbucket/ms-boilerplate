package sd

import (
	"log"
	"net"
)

const (

	//defaultCheckIntervel default value
	//if config doesnt have value it takes default value
	defaultCheckIntervel = "5s"

	//defaultCheckTimeout default value
	//if config doesnt have value it takes default value
	defaultCheckTimeout = "3s"
)

//ServiceInfo Struct
type ServiceConfig struct {
	ServiceID      string
	ServiceName    string
	ServiceAddress string
	ServicePort    int
	Tags           []string
	CheckIntervel  string
	CheckTimeout   string
}

//GetConfig Function
func DefaultConfig() (ServiceConfig, error) {

	s := ServiceConfig{}
	s.ServiceAddress = getOutboundIP()
	s.CheckIntervel = defaultCheckIntervel
	s.CheckTimeout = defaultCheckTimeout

	return s, nil
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
