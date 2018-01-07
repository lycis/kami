package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"gitlab.com/lycis/kami/util"
	"github.com/reiver/go-telnet"
	log "github.com/Sirupsen/logrus"
	flag "github.com/ogier/pflag"
	"fmt"
	"net/url"
	"gopkg.in/resty.v1"
	"encoding/json"
	"sync"
)

var logger *log.Logger
var backendUrl string
var callbackURL string
type telnetHandler struct {}
var users map[string]chan string
var usersMutex sync.Mutex

func init() {
	users = make(map[string]chan string)
}

type backend interface {
	Connect() (error)
	Disconnect() error
	SendInput(input string) error
	SubscribeToEvents() (chan string, error)
}

type callbackEvent struct {
	Token string `json:"token"`
	Payload []byte `json:"payload"`
}

var callbackEvents chan callbackEvent

func main() {
	util.PrintLicenseHint("Kami Telnet Frontend v0.1\n")

	port := flag.Int("port", 23, "port to listen on for user connections")
	nwif := flag.String("interface", "0.0.0.0", "network interface to listen on")
	backendConnection := flag.String("backend", "rest:localhost:8080", "connection to the kami backend in the form of <protocol>:<host>:<port>")
	cbu := flag.String("callback-address", "localhost", "IP address that the callback server will listen to")
	cbp := flag.Int("callback-port", 80,"port for the callback server")
	resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))

	if len(*cbu) < 1 {
		fmt.Printf("callback-address is missing")
		flag.PrintDefaults()
		return
	}

	callbackEvents = make(chan callbackEvent)

	logger = log.New()

	// callback handling
	go func() {
	    callbackURL = fmt.Sprintf("http://%s:%d/callback", *cbu, *cbp)
		logger.WithField("callback-url", callbackURL).Info("Starting callback webserver.")
		http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request){
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				log.WithFields(log.Fields{"method": r.Method, "endpoint": callbackURL}).Warn(fmt.Sprintf("Access with wrong method."))
				return
			}

			fmt.Printf("~~~~> received event\n")

			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.WithError(err).Error("Reading callback request body failed.")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var cbe callbackEvent
			if err := json.Unmarshal(body, &cbe); err != nil {
				log.WithError(err).WithField("endpoint", callbackURL).Error("Unmarshalling callback event failed.")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			callbackEvents <- cbe
			w.WriteHeader(http.StatusAccepted)
			fmt.Printf("~~~~> received event (end)\n")
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *cbp), nil))
	}()

	// dispatch callback requests
	go func() {
		for event := range callbackEvents {
			fmt.Printf("~~~~> dispatching event\n")
			/*usersMutex.Lock()
			defer usersMutex.Unlock()*/

			queue, ok := users[event.Token]
			if !ok {
				log.WithField("token", event.Token).Error("received event for unknown user token")
			} else {
				queue <- string(event.Payload)
			}

			fmt.Printf("~~~~> dispatching event (end)\n")
		}
	}()
	
	logger.WithField("backend", *backendConnection).Info("Using specified backend")
	backendUrl = *backendConnection
	logger.WithFields(log.Fields{"port": *port, "interface": *nwif}).Info("Waiting for user connections.")
	log.Fatal(telnet.ListenAndServe(fmt.Sprintf("%s:%d", *nwif, *port), telnetHandler{}))
}

// TODO implement callback subscriptions
func (t telnetHandler) ServeTELNET(context telnet.Context, w telnet.Writer, r telnet.Reader) {
	logger.Info("New user connected.")
	server, err := makeBackend(backendUrl)
	if err != nil {
		fmt.Fprintf(w, "Internal error. Please consult the server log.\n")
		logger.WithError(err).Error("Failed to make backend. Possibly an invalid configuration?")
		return
	}

	fmt.Fprintf(w, "Establishing connection...")
	if err := server.Connect(); err != nil {
		fmt.Fprintf(w, "[FAILED]\nerror: %s\n", err.Error())
		logger.WithError(err).Error("Connecting to backend failed.")
		return
	}
	fmt.Fprintf(w, "[OK]\n")


	fmt.Fprintf(w, "Registering callback connection...")
	outputChannel, err := server.SubscribeToEvents()
	if err != nil {
		fmt.Fprintf(w, "[FAILED]\nerror: %s\n", err.Error())
		logger.WithError(err).Error("failed to subscribe to events.")
		return
	}
	fmt.Fprintf(w, "[OK]\n")

	// write events to user
	go func(){
		for message := range outputChannel {
			fmt.Printf("~~~~> writing event\n")
			fmt.Fprintf(w, message)
			fmt.Printf("~~~~> writing event (done)\n")
		}
	}()

	defer func(){
		if err := server.Disconnect(); err != nil {
			logger.WithError(err).Error("Disconnect failed.")
		}
		logger.Info("User disconnected.")
	}()

	var buffer [1]byte
	input := buffer[:]
	var line string
	logger.Info("Waiting for user input.")
	fmt.Fprintf(w, "Connection is ready.\n")
	for {
		_, err := r.Read(input)
		if err != nil {
			logger.WithError(err).Error("Error reading input from client.")
			return
		}
		
		line += fmt.Sprintf("%s", buffer)
		if strings.HasSuffix(line, "\n") {
			line = strings.Trim(line, " \n\r\t")
			if err := server.SendInput(string(line)); err != nil {
				logger.WithError(err).Error("Transmission of input failed.")
				fmt.Fprintf(w, "Error: %s\n", err)
			}
			line = ""
		}
	}
}

func makeBackend(url string) (backend, error) {
	parts := strings.Split(url, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("incomplete backend configuration")
	}

	switch(parts[0]) {
	case "rest":
		return makeRestBackend(parts[1], parts[2])
	}

	return nil, fmt.Errorf("invalid backend protocol")
}

func makeRestBackend(host, port string) (backend, error) {
	u, err := url.Parse(fmt.Sprintf("http://%s:%s/", host, port))
	if err != nil {
		return nil, err
	}

	return &restBackend {
		token: "",
		url: u,
	}, nil
}

type restBackend struct {
	token string
	url *url.URL
}

type Token struct {
	Token string
}

func (r *restBackend) Connect() error {
	resp, err := resty.R().SetHeader("Accept", "application/json").Put(fmt.Sprintf("%s%s", r.url.String(), "user"))
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("Backend request failed: %s", resp.StatusCode())
	}

	var t Token
	err = json.Unmarshal(resp.Body(), &t)
	if err != nil {
		return err
	}

	if len(t.Token) <= 0 {
		return fmt.Errorf("no token received")
	}

	r.token = t.Token
	logger.WithField("token", r.token).Info("Received token for user.")

	return nil
}

func (r *restBackend) Disconnect() error {
	resp, err := resty.R().Delete(fmt.Sprintf("%suser/%s", r.url.String(), r.token))
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("Backend request failed: %s", resp.StatusCode())
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()
	close(users[r.token])

	return nil
}

func (r *restBackend) SendInput(input string) error {
	resp, err := resty.R().SetBody(input).Post(fmt.Sprintf("%suser/%s", r.url, r.token))
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("Backend request failed: %s (%s)", resp.StatusCode(), resp.Body())
	}

	logger.WithField("input", input).Info("Forwarded input to backend")

	return nil
}

type SubscribeEventRequest struct {
	Protocol string `json:'protocol'`
	URL      string `json:"url"`
}

func (r *restBackend) SubscribeToEvents() (chan string, error) {
	fmt.Printf("\nX\n")
	payload := SubscribeEventRequest{"rest", callbackURL}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nX O\n")

	ep := fmt.Sprintf("%suser/%s/callback", r.url, r.token)
    log.WithField("ep", ep).Info("Registering callback.")
	resp, err := resty.R().SetBody(string(body)).Put(ep)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nX O X\n")


	if resp.StatusCode() != http.StatusAccepted {
		return nil, fmt.Errorf("Registering callback failed: %s", resp.StatusCode())
	}

	fmt.Printf("\nX O X O\n")

	usersMutex.Lock()
	fmt.Printf("\nX O X O X\n")
    defer usersMutex.Unlock()
	users[r.token] = make(chan string)

	fmt.Printf("\nX O X O X O\n")
	return users[r.token], nil
}
