package subsystem

import (
	"encoding/json"
	"fmt"
	"net/http"
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

// --- REST INTERFACE ---
type nwi_rest struct {
	iface   string
	handler NetworkRequestHandler
}

func (nwi *nwi_rest) HandleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		token, err := nwi.handler.UserTokenRequested(nwi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		rstruct := struct{ Token string }{token}
		j, err := json.Marshal(rstruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(j)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func (nwi nwi_rest) Listen() error {
	var e error
	e = nil
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
		}
	}()

	go func() {
		if err := http.ListenAndServe(nwi.iface, nil); err != nil {
			panic(err)
		}
	}()

	return e
}

func (nwi nwi_rest) Close() {
	// TODO implement
}

func (nwi *nwi_rest) SetHandler(h NetworkRequestHandler) {
	nwi.handler = h
}

func nwi_create_rest(iface string) NetworkingInterface {
	nwi := &nwi_rest{iface: iface}

	http.HandleFunc("/user", nwi.HandleUser)

	return nwi
}
