package exceptions

type DomainError string

func (ee DomainError) Error() string {
	return string(ee)
}

var (
	NilValue        = DomainError("Value is nil")
	KeyRequired     = DomainError("Key is required")
	TypeNotFound    = DomainError("Type not found")
	ObjectNoteFound = DomainError("Object not found")
)
