package company

import (
	"context"
	"database/sql"
	_ "embed"

	sq "github.com/Masterminds/squirrel"
	"github.com/PanosXY/xm-platform/api/v1/company"
	"github.com/PanosXY/xm-platform/utils/postgres"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	pqDuplicateErrorCode pq.ErrorCode = "23505"

	tableCompanies = "companies"

	companiesColumnID          = "id"
	companiesColumnName        = "name"
	companiesColumnDescription = "description"
	companiesColumnEmployees   = "employees_amount"
	companiesColumnRegistered  = "registered"
	companiesColumnType        = "type"
)

var (
	//go:embed queries/find_company_by_id.sql
	findCompanyByID string

	//go:embed queries/insert_company.sql
	insertCompany string

	//go:embed queries/delete_company_by_id.sql
	deleteCompanyByID string

	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type companyRepository struct {
	dbClient *postgres.PostgresClient
}

// NewPostgresCompanyRepository returns a company postgres repo
func NewPostgresCompanyRepository(dbClient *postgres.PostgresClient) company.CompanyRepository {
	return &companyRepository{
		dbClient: dbClient,
	}
}

// GetCompanyByID returns a company by the given id
func (r *companyRepository) GetCompanyByID(ctx context.Context, id string) (*company.Company, error) {
	record := new(company.Company)

	err := r.dbClient.Database().QueryRowContext(ctx, findCompanyByID, id).Scan(
		&record.ID,
		&record.Name,
		&record.Description,
		&record.Employees,
		&record.Registered,
		&record.Type,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return record, err
}

// CreateCompany inserts a company
func (r *companyRepository) CreateCompany(ctx context.Context, companyReq *company.CompanyCreateRequest) error {
	if companyReq.ID == nil {
		id := uuid.NewString()
		companyReq.ID = &id
	}

	var (
		id   string
		name string
	)

	err := r.dbClient.Database().QueryRowContext(ctx, insertCompany,
		companyReq.ID,
		companyReq.Name,
		companyReq.Description,
		companyReq.Employees,
		companyReq.Registered,
		companyReq.Type,
	).Scan(&id, &name)

	if err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == pqDuplicateErrorCode {
			return company.ErrDuplicateKey
		}
		return err
	}

	return nil
}

// DeleteCompanyByID deletes a company by the given id
func (r *companyRepository) DeleteCompanyByID(ctx context.Context, id string) error {
	_, err := r.dbClient.Database().ExecContext(ctx, deleteCompanyByID, id)

	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

// PatchCompanyByID patches a company by the given id
func (r *companyRepository) PatchCompanyByID(ctx context.Context, id string, data *company.CompanyPatchRequest) error {
	patchData := make(map[string]interface{})

	if data.Description != nil {
		patchData[companiesColumnDescription] = *data.Description
	}
	if data.Name != nil {
		patchData[companiesColumnName] = *data.Name
	}
	if data.Employees != nil {
		patchData[companiesColumnEmployees] = *data.Employees
	}
	if data.Registered != nil {
		patchData[companiesColumnRegistered] = *data.Registered
	}
	if data.Type != nil {
		patchData[companiesColumnType] = *data.Type
	}

	sql, args, err := psql.Update(tableCompanies).SetMap(patchData).Where(sq.Eq{companiesColumnID: id}).ToSql()
	if err != nil {
		return err
	}

	if _, err := r.dbClient.Database().ExecContext(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}
