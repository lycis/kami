package subsystem

import (
	"fmt"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

type nwiRest struct {
	iface   string
	handler NetworkRequestHandler
	callbacks map[string]*url.URL
	engine *gin.Engine
}

func (nwi *nwiRest) DeleteToken(c *gin.Context) {
	token := c.Param("token")
	if len(token) < 1 {
		http.Error(c.Writer, "invalid token", http.StatusBadRequest)
		return
	}

	if err := nwi.handler.InvalidateUserToken(nwi, token); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}

func (nwi *nwiRest) UserInput(c *gin.Context) {
	defer c.Request.Body.Close()
	if c.Request.Method != http.MethodPost {
		http.Error(c.Writer, "", http.StatusMethodNotAllowed)
		return
	}

	token := c.Param("token")
	if len(token) < 1 {
		http.Error(c.Writer, "invalid token", http.StatusBadRequest)
		return
	}

	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(payload) < 1 {
		http.Error(c.Writer, "no payload given", http.StatusBadRequest)
		return
	}

	if err := nwi.handler.UserInputProvided(nwi, token, string(payload)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusAccepted)
}

func (nwi *nwiRest) NewUserToken(c *gin.Context) {
	if c.Request.Method == http.MethodPut {
		token, err := nwi.handler.UserTokenRequested(nwi)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		rstruct := struct{ Token string }{token}
		j, err := json.Marshal(rstruct)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(j)
	} else {
		http.Error(c.Writer, "", http.StatusMethodNotAllowed)
	}
}

type callbackRegisterPayload struct {
	Protocol string
	URL      string
}
func (nwi *nwiRest) RegisterCallback(c *gin.Context) {
	token := c.Param("token")
	if !nwi.handler.IsValidToken(token) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	defer c.Request.Body.Close()

	var payload callbackRegisterPayload
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(body) < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	url, err := url.Parse(payload.URL)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nwi.callbacks[token] = url
	c.Status(http.StatusAccepted)
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
		if err := nwi.engine.Run(nwi.iface); err != nil {
			panic(err)
		}
	}()

	return e
}

func (nwi nwiRest) Close() {
	// TODO implement
}

type  eventOccurance struct {
	Token string `json:"token"`
	Payload []byte `json:"payload"`
}

func (nwi nwiRest) RouteEvent(token string, payload []byte) error {
	url, ok := nwi.callbacks[token]
	if !ok {
		return fmt.Errorf("token has no callback")
	}

	str := &eventOccurance{token, payload}
	b, err := json.Marshal(str)
	if err != nil {
		return err
	}

	response, err := resty.R().SetBody(b).Post(url.String())
	if err != nil {
		return err
	}
	
	if response.StatusCode() > 299 {
		return fmt.Errorf("event routing failed")
	}

	return nil
}

func (nwi *nwiRest) SetHandler(h NetworkRequestHandler) {
	nwi.handler = h
}

func createRestInterface(iface string) NetworkingInterface {
	nwi := &nwiRest{iface: iface, callbacks: make(map[string]*url.URL)}

	nwi.engine = gin.Default()

	nwi.engine.PUT("/user", nwi.NewUserToken)
	nwi.engine.DELETE("/user/:token", nwi.DeleteToken)
	nwi.engine.POST("/user/:token", nwi.UserInput)
	nwi.engine.PUT("/user/:token/callback", nwi.RegisterCallback)
 
	return nwi
}
