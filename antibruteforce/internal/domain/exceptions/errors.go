package exceptions

// DomainError common errors
type DomainError string

// Error implementation error interface
func (ee DomainError) Error() string {
	return string(ee)
}

// Domain errors
const (
	NilValue         = DomainError("Value is nil")
	KeyRequired      = DomainError("Key is required")
	KindRequired     = DomainError("Bucket Hash is required")
	TypeNotFound     = DomainError("Type not found")
	ObjectNoteFound  = DomainError("Object not found")
	IPRequired       = DomainError("Ip address is not correct")
	IPInBlackList    = DomainError("Ip address in black list")
	LoginRequired    = DomainError("Login address is required")
	PasswordRequired = DomainError("Password address is required")
	LimitReached     = DomainError("Request limit has been reached")
)
