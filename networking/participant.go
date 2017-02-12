package networking

import (
	"net"
	"sync"
	"github.com/Sirupsen/logrus"
	"bufio"
	"io"
)


type Participant interface {
	// Attach attaches a component to the participant and
	// returns the ID that the component will be identified with
	Attach(c Component) string
}

type _participant struct {
	// components holds a mapping of an ID to the attached component
	components map[string]Component

	// this is the mutex for all locking actions in a concurrent
	// environment
	mutex sync.Mutex

	listener net.Listener
}

func NewParticipant() *Participant {
	p := &_participant{
		components: make(map[string]Component),
	}

	return p
}

func (self *_participant) Listen(net, address string) error {
	self.listener, err := net.Listen(net, address)
	if err != nil {
		return err
	}


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
		} else if err != nil {
			logrus.WithField("remote", conn.RemoteAddr()).Warn("Connection closed.")
		}


	}
}

func (self *_participant) Attach(c Component) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.components[c.ID()] = c
}
