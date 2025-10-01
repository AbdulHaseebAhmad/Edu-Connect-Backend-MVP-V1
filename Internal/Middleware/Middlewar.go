package Middlewares

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func Authorizer(storage Storage.Storage, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			slog.Info("Session Cookie Missing", "ok:", sessionToken)
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		sessionCookie := sessionToken.Value
		slog.Info("Session Cookie Recieved", "ok:", sessionCookie)

		csrfCookie := r.Header.Get("X-CSRF-Token")
		if csrfCookie == "" {
			slog.Info("CSRF Cookie Missing", "error:", csrfCookie)
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		slog.Info("CSRF Cookie Recieved", "ok:", csrfCookie)

		isAuthorised := storage.AuthorizeSysAdmin(r.Context(), sessionCookie, csrfCookie)
		if !isAuthorised {
			slog.Info("sessionCokie & CSRFToken mismatch", "error:", "unauthorized Access")
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		next.ServeHTTP(w, r)
	}
}
