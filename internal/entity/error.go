package entity

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const (
	EmptyTagParameterReason   = "tag parameter is required"
	InvalidTagParameterReason = "tag parameter should be 'error', 'info', 'debug' or 'trace'"
	EmptyDataParameterReason  = "data parameter is required"
)

var (
	ErrEmptyTagParameter   = errors.New("tag parameter is required")
	ErrInvalidTagParameter = errors.New("tag parameter should be 'error', 'info', 'debug' or 'trace'")
	ErrEmptyDataParameter  = errors.New("data parameter is required")
)

// errResp represents a structure whose fields can be passed as a response to the NewErrorResponse.
type errResp struct {
	Message string `json:"message"`
}

// NewErrorResponse represents an error response for handlers.
func NewErrorResponse(c *gin.Context, codeStatus int, message string) {
	c.AbortWithStatusJSON(codeStatus, errResp{message})
}

func (e *errResp) BuildErrByReason(code int, reason string) error {
	return errors.New(reason)
}
