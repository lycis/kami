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

	// RouteEvent will pass an event to the frontend that
	// registered for events of the given token.
	RouteEvent(token string, payload []byte) error
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

	// InvalidateUserToken will be called when the client wants to invalidate
	// a user token (e.g. a log off ocurred) or the user terminated the
	// connection.
	InvalidateUserToken(nwi NetworkingInterface, token string) error

	// IsValidToken has to return 'true' if the given token is valid and equals
	// to a still active token that was previously created.
	IsValidToken(token string) bool
}

const (
	// NWI_REST is the switch for a REST networking interface
	InterfaceREST = 0

	// NWI_TCP is the marker for using the raw tcp protocol
	InterfaceTCP  = 1
)

// CreateNetworkingInterface provides a new and prepared networking interface for a
// given type on the given address and port
func CreateNetworkingInterface(kind int, listenAddress string, port int) NetworkingInterface {
	switch kind {
	case InterfaceREST:
		return createRestInterface(fmt.Sprintf("%s:%d", listenAddress, port))
	default:
		panic(fmt.Errorf("unsupported networking interface type: %d", kind))
	}
}
