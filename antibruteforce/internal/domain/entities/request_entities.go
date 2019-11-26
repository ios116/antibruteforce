package entities

import (
	"antibruteforce/internal/domain/exceptions"
)

// Request - a request for approved.
type Request struct {
	IP       string
	Login    string
	Password string
}

// Validation - validation request.
func (r *Request) Validation() error {
	if r.IP == "" {
		return exceptions.IPRequired
	}
	if r.Login == "" {
		return exceptions.LoginRequired
	}
	if r.Password == "" {
		return exceptions.PasswordRequired
	}
	return nil
}
