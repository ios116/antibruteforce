package entities

import (
	"antibruteforce/internal/domain/exceptions"
	"net"
)

// Request - a request for approved.
type Request struct {
	IP       *net.IP
	Login    string
	Password string
}

// Validation - validation request.
func (r *Request) Validation() error {
	if r.IP == nil {
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
