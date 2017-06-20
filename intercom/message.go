package intercom

type Message struct {
	Envelope   *Envelope  `json:"envelope,omitempty""`
	Advertise  *Advertise `json:"advertise,omitempty"`
	Neighbours Neighbours `json:"neighbours,omitempty"`
}

type Envelope struct {
	To    string   `json:"to""`
	From  string   `json:"from""`
	SeqId int      `json:"seq-id""`
	Relay []string `json:"relay"`
}

type Advertise struct {
	Address    []string `json:"address"`
	Extensions string   `json:"extensions"`
}

type Neighbours map[string]*Advertise
