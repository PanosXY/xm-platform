package response

import (
	"net/http"

	"github.com/go-chi/render"
)

var messages = map[string]ResponseMessage{
	"GenericSuccess":       {Code: 10000, Description: "GenericSuccess", Message: "Request completed"},
	"GenericError":         {Code: 10001, Description: "GenericError", Message: "Request failed"},
	"GenericBadRequest":    {Code: 10002, Description: "GenericBadRequest", Message: "Bad request"},
	"GenericServerError":   {Code: 10003, Description: "GenericServerError", Message: "Server Error"},
	"GenericHealthError":   {Code: 10005, Description: "GenericHealthError", Message: "Health check failed"},
	"GenericNotAuthorized": {Code: 10006, Description: "GenericNotAuthorized", Message: "Not authorized"},

	"InvalidCompanyID":      {Code: 20000, Description: "InvalidCompanyID", Message: "Invalid company ID"},
	"CompanyGetError":       {Code: 20001, Description: "CompanyGetError", Message: "Failed to get company"},
	"InvalidCompanyPayload": {Code: 20002, Description: "InvalidCompanyPayload", Message: "Invalid company payload"},
	"CompanyCreateError":    {Code: 20003, Description: "CompanyCreateError", Message: "Failed to create company"},
	"CompanyAlreadyExists":  {Code: 20004, Description: "CompanyAlreadyExists", Message: "Company already exists"},
	"CompanyDeleteError":    {Code: 20005, Description: "CompanyDeleteError", Message: "Failed to delete company"},
	"CompanyPatchError":     {Code: 20006, Description: "CompanyPatchError", Message: "Failed to patch company"},

	"LoginCredentialsError":   {Code: 30000, Description: "LoginCredentialsError", Message: "Failed to decode credentials"},
	"LoginUnauthorized":       {Code: 30001, Description: "LoginUnauthorized", Message: "User is not authorized"},
	"BadAuthenticationHeader": {Code: 30002, Description: "BadAuthenticationHeader", Message: "Bad authentication header"},
}

type Response struct{}

// ResponseMessage is used for a metadata section in a response.
type ResponseMessage struct {
	Description string `json:"description,omitempty"`
	Message     string `json:"message,omitempty"`
	Code        int    `json:"code,omitempty"`
}

type Metadata struct {
	Success bool `json:"success"`
}

func (r *Response) GetMessage(code string) ResponseMessage {
	return messages[code]
}

// ResponseData returns a response formatted as JSON.
type ResponseData struct {
	Data   interface{}       `json:"data,omitempty"`
	Meta   interface{}       `json:"meta,omitempty"`
	Errors []ResponseMessage `json:"errors,omitempty"`
}

// JSONAPIResponseWithSuccess returns a successful response.
func JSONAPIResponseWithSuccess(w http.ResponseWriter, r *http.Request, httpStatus int, data interface{}) {
	render.Status(r, httpStatus)
	render.JSON(w, r, ResponseData{
		Data: data,
		Meta: Metadata{true},
	})
}

// JSONAPIResponseWithError returns a failed response.
func JSONAPIResponseWithError(w http.ResponseWriter, r *http.Request, httpStatus int, errors ResponseMessage) {
	render.Status(r, httpStatus)
	render.JSON(w, r, ResponseData{
		Meta:   Metadata{false},
		Errors: []ResponseMessage{errors},
	})
}
