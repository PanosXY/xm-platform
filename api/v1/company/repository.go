package company

import (
	"context"
)

// CompanyRepository includes the company methods
type CompanyRepository interface {
	GetCompanyByID(ctx context.Context, id string) (*Company, error)
	CreateCompany(ctx context.Context, company *CompanyCreateRequest) error
	DeleteCompanyByID(ctx context.Context, id string) error
	PatchCompanyByID(ctx context.Context, id string, data *CompanyPatchRequest) error
}
