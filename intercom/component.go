package intercom

// Component is the interface that everything that wants to attach to a
// participant has to implement
type Component interface {
	ID() string
}
