package intercom

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net"
	"sync"
)

type Participant interface {
	// Attach attaches a component to the participant. Duplicate component
	// ids might me rejected
	Attach(c Component) error

	// Starts the participant to listen and handle snippet deliveries
	Listen(netname, address string) error

	// Disconnect stops to listen for incoming snippets
	Close() error

	// Ready checks checks whether the partcipant is ready for
	// connections and responsive.
	Ready() bool
}

type _participant struct {
	// components holds a mapping of an ID to the attached component
	components map[string]Component

	// this is the mutex for all locking actions in a concurrent
	// environment
	mutex sync.Mutex

	listener net.Listener

	listening bool
}

func NewParticipant() Participant {
	p := &_participant{
		components: make(map[string]Component),
	}

	return p
}

func (self *_participant) Listen(netname, address string) error {
	l, err := net.Listen(netname, address)
	if err != nil {
		return err
	}

	self.listener = l
	self.listening = true

	for {
		c, err := self.listener.Accept()
		if err != nil {
			return err
		}

		go self.handleConnection(c)
	}

	return nil
}

func (self *_participant) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			logrus.WithError(err).WithField("remote", conn.RemoteAddr()).Warn("Error reading from connection.")
			continue
		} else if err != nil {
			logrus.WithField("remote", conn.RemoteAddr()).Warn("Connection closed.")
			return
		}

		jsonStr, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			logrus.WithError(err).WithField("remote", conn.RemoteAddr()).Error("Received invalid bas64 payload.")
			continue
		}

		var m Message
		if json.Unmarshal([]byte(jsonStr), &m) != nil {
			continue
			logrus.WithError(err).WithField("remote", conn.RemoteAddr()).Error("Received invalid json payload.")
		}

		if m.Envelope != nil {
			logrus.WithFields(logrus.Fields{"remote": conn.RemoteAddr(), "address": m.Envelope.To}).Info("Received information snippet.")
		} else {
			logrus.WithFields(logrus.Fields{"remote": conn.RemoteAddr()}).Info("Received participant information snippet.")
		}
	}
}

func (self *_participant) Attach(c Component) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if _, exists := self.components[c.ID()]; exists {
		return fmt.Errorf("duplicate component '%s' registered", c.ID())
	}

	self.components[c.ID()] = c

	return nil
}

func (self *_participant) Close() error {
	return self.listener.Close()
}

func (self *_participant) Ready() bool {
	return self.listener != nil && self.listening
}
