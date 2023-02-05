package request

import (
	"context"

	"github.com/google/uuid"
)

type ctxRequestID int

const (
	requestIDHttpHeader = "x-request-id"

	requestIDKey ctxRequestID = 0
)

// GetRequestID if exists returns a request ID from the given context, otherwise creates and returns a new one.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}

	return uuid.NewString()
}
