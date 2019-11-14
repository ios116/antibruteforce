package exceptions

// DomainError common errors
type DomainError string

// Error implementation error interface
func (ee DomainError) Error() string {
	return string(ee)
}

var (
	NilValue         = DomainError("Value is nil")
	KeyRequired      = DomainError("Key is required")
	TypeNotFound     = DomainError("Type not found")
	ObjectNoteFound  = DomainError("Object not found")
	IPRequired       = DomainError("Ip address is required")
	LoginRequired    = DomainError("Login address is required")
	PasswordRequired = DomainError("Password address is required")
	LimitReached = DomainError("Request limit has been reached")
)
