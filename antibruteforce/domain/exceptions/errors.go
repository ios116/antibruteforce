package exceptions

type DomainError string

func (ee DomainError) Error() string {
	return string(ee)
}

var (
	BucketsNil      = DomainError("Bucket is nil")
	KeyRequired      = DomainError("Key is required")
	TypeNotFound      = DomainError("Type not found")
)