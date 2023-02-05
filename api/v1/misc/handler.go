package misc

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PanosXY/xm-platform/config"
	"github.com/PanosXY/xm-platform/response"
	"github.com/PanosXY/xm-platform/utils/logger"
	"github.com/PanosXY/xm-platform/utils/request"
	"github.com/golang-jwt/jwt"
)

// Some random users for login purposes
var users = map[string]string{
	"xm":    "xm123",
	"panos": "passwd123",
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// jwtClaims will be encoded to JWT
type jwtClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// MiscHandler miscellaneous endpoints handler
type MiscHandler interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	JWTAuthToken(next http.Handler) http.Handler
}

type miscHandler struct {
	configuration *config.Configuration
	log           *logger.Logger
	miscService   MiscService
	response      *response.Response
}

// NewMiscHandler returns a new miscellaneous endpoints handler
func NewMiscHandler(configuration *config.Configuration, log *logger.Logger, miscService MiscService) MiscHandler {
	return &miscHandler{
		configuration: configuration,
		log:           log,
		miscService:   miscService,
	}
}

// HealthCheck does a db health check
func (h *miscHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	component := "HealthCheck"
	ctx := r.Context()
	reqID := request.GetRequestID(ctx)

	if err := h.miscService.DoHealthCheck(ctx); err != nil {
		h.log.Errorf(reqID, component, "an error occurred during healthcheck", err, nil)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("GenericHealthError"))
		return
	}

	response.JSONAPIResponseWithSuccess(w, r, http.StatusNoContent, nil)
}

// Login authenticates a user
func (h *miscHandler) Login(w http.ResponseWriter, r *http.Request) {
	component := "Login"
	ctx := r.Context()
	reqID := request.GetRequestID(ctx)

	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		h.log.Errorf(reqID, component, "failed to decode credentials", err, nil)
		response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("LoginCredentialsError"))
		return
	}

	fields := logger.Fields{
		"username": creds.Username,
	}

	passwd, ok := users[creds.Username]
	if !ok || passwd != creds.Password {
		h.log.Errorf(reqID, component, "invalid credentials for user", nil, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusUnauthorized, h.response.GetMessage("LoginUnauthorized"))
		return
	}

	expireTime := time.Now().Add(24 * time.Hour)
	claims := &jwtClaims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.configuration.HttpServer.JWTSecretKey)
	if err != nil {
		h.log.Errorf(reqID, component, "failed to  sign token", err, fields)
		response.JSONAPIResponseWithError(w, r, http.StatusInternalServerError, h.response.GetMessage("GenericServerError"))
		return
	}

	res := &Token{
		Token:   tokenString,
		Expires: expireTime,
	}

	response.JSONAPIResponseWithSuccess(w, r, http.StatusOK, res)
}
