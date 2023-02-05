package company

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/PanosXY/xm-platform/response"
	"github.com/PanosXY/xm-platform/utils/helpers"
	"github.com/PanosXY/xm-platform/utils/logger"
	"github.com/PanosXY/xm-platform/utils/request"
	"github.com/asaskevich/govalidator"
	"github.com/go-chi/chi/v5"
)

const (
	operationContextTimeout = time.Second * 2
)

// CompanyHandler includes company handler's operations
type CompanyHandler interface {
	GetCompanyByID(w http.ResponseWriter, r *http.Request)
	CreateCompany(w http.ResponseWriter, r *http.Request)
	DeleteCompanyByID(w http.ResponseWriter, r *http.Request)
	PatchCompanyByID(w http.ResponseWriter, r *http.Request)
}

type companyHandler struct {
	log            *logger.Logger
	response       *response.Response
	companyService CompanyService
}

//  NewCompanyHandler returns a new company handler
func NewCompanyHandler(log *logger.Logger, companyService CompanyService) CompanyHandler {
	return &companyHandler{
		log:            log,
		companyService: companyService,
	}
}

// GetCompanyID returns a company by the given id
func (h *companyHandler) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	component := "GetCompanyByID"
	rCtx := r.Context()
	reqID := request.GetRequestID(rCtx)

	companyID := chi.URLParam(r, "id")

	fields := logger.Fields{
		"id": companyID,
	}

	if err := helpers.IsValidUUID(companyID); err != nil {
		h.log.Errorf(reqID, component, "invalid company id", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyID"))
		return
	}

	ctx, cancel := context.WithTimeout(rCtx, operationContextTimeout)
	defer cancel()

	company, err := h.companyService.GetCompanyByID(ctx, companyID)
	if err != nil {
		h.log.Errorf(reqID, component, "failed to get company", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("CompanyGetError"))
		return
	}

	h.log.Infof(reqID, component, "got company", fields)
	response.JSONAPIResponseWithSuccess(w, r, http.StatusOK, company)
}

// CreateCompany creates a new company with the given data
func (h *companyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	component := "CreateCompany"
	rCtx := r.Context()
	reqID := request.GetRequestID(rCtx)

	var payload CompanyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Errorf(reqID, component, "failed to decode payload", err, nil)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyPayload"))
		return
	}

	fields := logger.Fields{
		"payload": payload,
	}

	if _, err := govalidator.ValidateStruct(payload); err != nil {
		h.log.Errorf(reqID, component, "invalid company payload", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyPayload"))
		return
	}

	ctx, cancel := context.WithTimeout(rCtx, operationContextTimeout)
	defer cancel()

	if err := h.companyService.CreateCompany(ctx, &payload); err != nil {
		if errors.Is(err, ErrDuplicateKey) {
			h.log.Errorf(reqID, component, "company already exists", err, fields)
			response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("CompanyAlreadyExists"))
			return
		}

		h.log.Errorf(reqID, component, "failed to get company", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("CompanyCreateError"))
		return
	}

	h.log.Infof(reqID, component, "added company", fields)
	response.JSONAPIResponseWithSuccess(w, r, http.StatusNoContent, nil)
}

// DeleteCompanyByID deletes a company by the given id
func (h *companyHandler) DeleteCompanyByID(w http.ResponseWriter, r *http.Request) {
	component := "DeleteCompanyByID"
	rCtx := r.Context()
	reqID := request.GetRequestID(rCtx)

	companyID := chi.URLParam(r, "id")

	fields := logger.Fields{
		"id": companyID,
	}

	if err := helpers.IsValidUUID(companyID); err != nil {
		h.log.Errorf(reqID, component, "invalid company id", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyID"))
		return
	}

	ctx, cancel := context.WithTimeout(rCtx, operationContextTimeout)
	defer cancel()

	if err := h.companyService.DeleteCompanyByID(ctx, companyID); err != nil {
		h.log.Errorf(reqID, component, "failed to delete company", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("CompanyDeleteError"))
		return
	}

	h.log.Infof(reqID, component, "deleted company", fields)
	response.JSONAPIResponseWithSuccess(w, r, http.StatusNoContent, nil)
}

// PatchCompanyByID patches a company by the given id
func (h *companyHandler) PatchCompanyByID(w http.ResponseWriter, r *http.Request) {
	component := "PatchCompanyByID"
	rCtx := r.Context()
	reqID := request.GetRequestID(rCtx)

	companyID := chi.URLParam(r, "id")

	fields := logger.Fields{
		"id": companyID,
	}

	if err := helpers.IsValidUUID(companyID); err != nil {
		h.log.Errorf(reqID, component, "invalid company id", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyID"))
		return
	}

	var payload CompanyPatchRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Errorf(reqID, component, "failed to decode payload", err, nil)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyPayload"))
		return
	}

	fields["payload"] = payload

	if _, err := govalidator.ValidateStruct(payload); err != nil {
		h.log.Errorf(reqID, component, "invalid company payload", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("InvalidCompanyPayload"))
		return
	}

	ctx, cancel := context.WithTimeout(rCtx, operationContextTimeout)
	defer cancel()

	if err := h.companyService.PatchCompanyByID(ctx, companyID, &payload); err != nil {
		h.log.Errorf(reqID, component, "failed to patch company", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("CompanyPatchError"))
		return
	}

	h.log.Infof(reqID, component, "patched company", fields)
	response.JSONAPIResponseWithSuccess(w, r, http.StatusNoContent, nil)
}
