package subsystem

import (
	"fmt"
)

// NetworkingInterface is the abstract definition
// for some kind of network input interface that
// the driver responds to
type NetworkingInterface interface {

	// starts the NWI to await and pass over
	// requests
	Listen() error

	// stops the NWI and cancels all connections
	Close()

	// defines the handler that will take care of
	// incoming request
	SetHandler(handler NetworkRequestHandler)
}

// A NetworkRequestHandler is the interface between
// the driver and the NetworkingInterface. The defined
// functions will be called when requests are received
// via the NetworkingInterface tht the handler is attached
// to
type NetworkRequestHandler interface {
	// UserTokenRequested will be called whenever a new
	// user connects and thus requires a token for interaction
	UserTokenRequested(nwi NetworkingInterface) (string, error)

	// UserInputProvided is called when new input from the user that
	// is identified by the token is available for processing
	UserInputProvided(nwi NetworkingInterface, token, input string) error
}

const (
	NWI_REST = 0
	NWI_TCP  = 1
)

// CreateNetworkingInterface provides a new and prepared networking interface for a
// given type on the given address and port
func CreateNetworkingInterface(kind int, listenAddress string, port int) NetworkingInterface {
	switch kind {
	case NWI_REST:
		return nwi_create_rest(fmt.Sprintf("%s:%d", listenAddress, port))
	default:
		panic(fmt.Errorf("unsupported networking interface type: %d", kind))
	}
}
