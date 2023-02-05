package company

import (
	"context"
)

// CompanyService includes the company operations
type CompanyService interface {
	GetCompanyByID(ctx context.Context, id string) (*Company, error)
	CreateCompany(ctx context.Context, company *CompanyCreateRequest) error
	DeleteCompanyByID(ctx context.Context, id string) error
	PatchCompanyByID(ctx context.Context, id string, data *CompanyPatchRequest) error
}

type companyService struct {
	postgresRepo CompanyRepository
}

// NewCompanyService returns a company service
func NewCompanyuService(postgresRepo CompanyRepository) CompanyService {
	return &companyService{
		postgresRepo: postgresRepo,
	}
}

// GetCommpanyByID returns a company by the given id
func (s *companyService) GetCompanyByID(ctx context.Context, id string) (*Company, error) {
	return s.postgresRepo.GetCompanyByID(ctx, id)
}

// CreateCompany creates a company
func (s *companyService) CreateCompany(ctx context.Context, company *CompanyCreateRequest) error {
	return s.postgresRepo.CreateCompany(ctx, company)
}

// DeleteCommpanyByID deletes a company by the given id
func (s *companyService) DeleteCompanyByID(ctx context.Context, id string) error {
	return s.postgresRepo.DeleteCompanyByID(ctx, id)
}

// PatchCommpanyByID patches a company by the given id
func (s *companyService) PatchCompanyByID(ctx context.Context, id string, data *CompanyPatchRequest) error {
	return s.postgresRepo.PatchCompanyByID(ctx, id, data)
}
