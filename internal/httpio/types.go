package httpio

import (
	"fmt"
)

// Error represents an application-specific error structure.
// It implements the error interface, allowing it to be used as an error type.
type Error struct {
	Status int
	Code   string
	Desc   string
}

func (e Error) Error() string {
	return fmt.Sprintf("respond.Error{status:%d,code:%s,desc:%s}", e.Status, e.Code, e.Desc)
}

// Message represents a generic message structure used for communication.
// Code is message code, Desc is message description
type Message struct {
	Code string `json:"code"`
	Desc string `json:"desc,omitempty"`
}
