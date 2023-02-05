package misc

import (
	"context"

	"github.com/PanosXY/xm-platform/api/v1/misc"
	"github.com/PanosXY/xm-platform/utils/postgres"
)

type miscRepository struct {
	dbClient *postgres.PostgresClient
}

// NewPostgresMiscRepository returns a miscellaneous postgres repo
func NewPostgresMiscRepository(dbClient *postgres.PostgresClient) misc.MiscRepository {
	return &miscRepository{
		dbClient: dbClient,
	}
}

// DoHealthCheck returns whether the db is healthy or not
func (r *miscRepository) DoHealthCheck(ctx context.Context) error {
	if err := r.dbClient.Ping(ctx); err != nil {
		return err
	}

	return nil
}
