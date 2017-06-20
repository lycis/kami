package intercom

import (
	"net"
	"testing"
	"time"
)

func startParticipant(t *testing.T) Participant {
	p := NewParticipant()
	go func() {
		err := p.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Log("listen failed: %s", err)
			t.Fail()
			return
		}
	}()

	for i := 0; !p.Ready(); i++ {
		if i > 10 {
			t.Log("participant did not start to listen")
			t.Fail()
		}
		time.Sleep(time.Second)
	}

	return p
}

func TestListen(t *testing.T) {
	p := startParticipant(t)

	socket, err := net.Dial(p.Addr().Network(), p.Addr().String())
	if err != nil {
		t.Log("connection failed: %s", err)
		t.Fail()
		return
	}

	if socket.Close() != nil {
		t.Log("socket close failed: %s", err)
		t.Fail()
		return
	}

	if p.Close() != nil {
		t.Log("close connection failed: %s", err)
		t.Fail()
		return
	}
}

func TestStartup(t *testing.T) {
	a := startParticipant(t)
	b := startParticipant(t)

	a.AddNeighbour(b.Addr().String())
	b.AddNeighbour(a.Addr().String())

	a.Close()
	b.Close()
}
