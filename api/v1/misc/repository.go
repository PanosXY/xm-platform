package misc

import "context"

// MiscRepository represents the misc methods
type MiscRepository interface {
	DoHealthCheck(context.Context) error
}
