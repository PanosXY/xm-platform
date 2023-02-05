package misc

import "context"

// MiscService is the miscellaneous service interface
type MiscService interface {
	DoHealthCheck(context.Context) error
}

type miscService struct {
	postgresRepo MiscRepository
}

// NewMiscService returns a new miscellaneous service
func NewMiscService(postgresRepo MiscRepository) MiscService {
	return &miscService{
		postgresRepo: postgresRepo,
	}
}

// DoHealthCheck returns whether the db is healthy or not
func (s *miscService) DoHealthCheck(ctx context.Context) error {
	return s.postgresRepo.DoHealthCheck(ctx)
}
