package privilege

const (
	// No privileges granted. Meaning no access to any efuns
	PrivilegeNone = 0

	// Basic privilege level for all objects
	PrivilegeBasic = 1

	// Highest privilege level. Allows access to all, even insecure and security functions
	PrivilegeRoot = 10
)

type Level int
