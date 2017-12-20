package subsystem

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type nwiRest struct {
	iface   string
	handler NetworkRequestHandler
}

func (nwi *nwiRest) HandleUser(w http.ResponseWriter, r *http.Request) {
	subpath := r.URL.Path[len("/user"):]
	if len(subpath) <= 0 {
		subpath = "/"
	}

	if subpath == "/" {
		nwi.NewUserToken(w, r)
		return
	}

	rex, err := regexp.Compile("^/([A-Za-z0-9\\-]+)/{0,1}$")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rex.MatchString(subpath) {
		token := rex.FindStringSubmatch(subpath)
		switch(r.Method) {
		case http.MethodPost:
			nwi.UserInput(w, r, token[1])
		case http.MethodDelete:
			nwi.DeleteToken(w, r, token[1])
		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "", http.StatusNotFound)
}

func (nwi *nwiRest) DeleteToken(w http.ResponseWriter, r *http.Request, token string) {
	if len(token) < 1 {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	if err := nwi.handler.InvalidateUserToken(nwi, token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (nwi *nwiRest) UserInput(w http.ResponseWriter, r *http.Request, token string) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	if len(token) < 1 {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(payload) < 1 {
		http.Error(w, "no payload given", http.StatusBadRequest)
		return
	}

	if err := nwi.handler.UserInputProvided(nwi, token, string(payload)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (nwi *nwiRest) NewUserToken(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		token, err := nwi.handler.UserTokenRequested(nwi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rstruct := struct{ Token string }{token}
		j, err := json.Marshal(rstruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(j)
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func (nwi nwiRest) Listen() error {
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

func (nwi nwiRest) Close() {
	// TODO implement
}

func (nwi *nwiRest) SetHandler(h NetworkRequestHandler) {
	nwi.handler = h
}

func createRestInterface(iface string) NetworkingInterface {
	nwi := &nwiRest{iface: iface}

	http.HandleFunc("/user/", nwi.HandleUser)

	return nwi
}
