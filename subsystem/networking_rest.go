package subsystem

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

type nwi_rest struct {
	iface   string
	handler NetworkRequestHandler
}

func (nwi *nwi_rest) HandleUser(w http.ResponseWriter, r *http.Request) {
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
		nwi.UserInput(w, r, token[1])
		return
	}

	http.Error(w, "", http.StatusNotFound)
}

func (nwi *nwi_rest) UserInput(w http.ResponseWriter, r *http.Request, token string) {
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
		http.Error(w, "no payload given", http.StatusInternalServerError)
		return
	}

	if err := nwi.handler.UserInputProvided(nwi, token, string(payload)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (nwi *nwi_rest) NewUserToken(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/user/", nwi.HandleUser)

	return nwi
}
