package company

// Company represents a company
type Company struct {
	ID          string `json::id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Employees   int    `json:"employees"`
	Registered  bool   `json:"registered"`
}

// CompanyCreateRequest represents a Create request
type CompanyCreateRequest struct {
	ID          *string `json::id" valid:"uuid"`
	Description *string `json:"description" valid:"stringlength(0|3000)"`
	Name        string  `json:"name" valid:"required,stringlength(1|15)"`
	Type        string  `json:"type" valid:"required,in(Corporations|NonProfit|Cooperative|Sole Proprietorship)"`
	Employees   int     `json:"employees" valid:"required"`
	Registered  bool    `json:"registered"`
}

// CompanyPatchRequest represents a Patch request
type CompanyPatchRequest struct {
	Description *string `json:"description" valid:"stringlength(0|3000)"`
	Name        *string `json:"name" valid:"stringlength(1|15)"`
	Type        *string `json:"type" valid:"in(Corporations|NonProfit|Cooperative|Sole Proprietorship)"`
	Employees   *int    `json:"employees"`
	Registered  *bool   `json:"registered"`
}
