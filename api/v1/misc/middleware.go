package misc

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PanosXY/xm-platform/response"
	"github.com/golang-jwt/jwt"
)

func (h *miscHandler) JWTAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			response.JSONAPIResponseWithError(w, r, http.StatusUnauthorized, h.response.GetMessage("GenericNotAuthorized"))
			return
		}

		authSlc := strings.Split(auth, " ")
		if len(authSlc) != 2 || authSlc[0] != "Bearer" {
			http.Error(w, "Wrong Authoriazation header format", http.StatusBadRequest)
			response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("BadAuthenticationHeader"))
			return
		}

		var claims jwtClaims
		token := authSlc[1]
		tkn, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect signing method: %v", t.Header["alg"])
			}
			return h.configuration.HttpServer.JWTSecretKey, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				response.JSONAPIResponseWithError(w, r, http.StatusUnauthorized, h.response.GetMessage("GenericNotAuthorized"))
				return
			}

			response.JSONAPIResponseWithError(w, r, http.StatusBadRequest, h.response.GetMessage("BadAuthenticationHeader"))
			return
		}

		if !tkn.Valid {
			response.JSONAPIResponseWithError(w, r, http.StatusUnauthorized, h.response.GetMessage("GenericNotAuthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
