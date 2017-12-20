package main

import (
	"strings"
	"gitlab.com/lycis/kami/util"
	"github.com/reiver/go-telnet"
	log "github.com/Sirupsen/logrus"
	flag "github.com/ogier/pflag"
	"fmt"
	"net/url"
	"gopkg.in/resty.v1"
	"encoding/json"
)

var logger *log.Logger
var backendUrl string
type telnetHandler struct {}

type backend interface {
	Connect() (error)
	Disconnect() error
	SendInput(input string) error
	SubscribeToEvents() (chan string, error)
}

func main() {
	util.PrintLicenseHint("Kami Telnet Frontend v0.1\n")

	port := flag.Int("port", 23, "port to listen on for user connections")
	nwif := flag.String("interface", "0.0.0.0", "network interface to listen on")
	backendConnection := flag.String("backend", "rest:localhost:8080", "connection to the kami backend in the form of <protocol>:<host>:<port>")
	resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(20))

	logger = log.New()
	
	logger.WithField("backend", *backendConnection).Info("Using specified backend")
	backendUrl = *backendConnection
	logger.WithFields(log.Fields{"port": *port, "interface": *nwif}).Info("Waiting for user connections.")
	log.Fatal(telnet.ListenAndServe(fmt.Sprintf("%s:%d", *nwif, *port), telnetHandler{}))
}

func (t telnetHandler) ServeTELNET(context telnet.Context, w telnet.Writer, r telnet.Reader) {
	logger.Info("New user connected.")
	server, err := makeBackend(backendUrl)
	if err != nil {
		fmt.Fprintf(w, "Internal error. Please consult the server log.\n")
		logger.WithError(err).Error("Failed to make backend. Possibly an invalid configuration?")
		return
	}

	if err := server.Connect(); err != nil {
		fmt.Fprintf(w, "Failed to connect to backend: %s\n", err.Error())
		logger.WithError(err).Error("Connecting to backend failed.")
		return
	}

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
	resp, err := resty.R().SetHeader("Accept", "application/json").Put(fmt.Sprintf("%s%s", r.url.String(), "user/"))
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
	resp, err := resty.R().Delete(fmt.Sprintf("%suser/%s/", r.url.String(), r.token))
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("Backend request failed: %s", resp.StatusCode())
	}

	return nil
}

func (r *restBackend) SendInput(input string) error {
	resp, err := resty.R().SetBody(input).Post(fmt.Sprintf("%suser/%s/", r.url, r.token))
	if err != nil {
		return err
	}

	if resp.StatusCode() > 299 {
		return fmt.Errorf("Backend request failed: %s", resp.StatusCode())
	}

	logger.WithField("input", input).Info("Forwarded input to backend")

	return nil
}

func (r *restBackend) SubscribeToEvents() (chan string, error) {
	return nil, nil
}

