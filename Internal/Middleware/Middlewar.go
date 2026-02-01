package Middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func Authorizer(storage Storage.SysAdmin, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			slog.Info("Session Cookie Missing", "ok:", sessionToken)
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		sessionCookie := sessionToken.Value
		slog.Info("Session Cookie Recieved", "ok:", sessionCookie)

		csrfCookie, _ := r.Cookie("X-CSRF-Token")
		if csrfCookie == "" {
			slog.Info("CSRF Cookie Missing", "error:", csrfCookie)
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		slog.Info("CSRF Cookie Recieved", "ok:", csrfCookie)

		id, isAuthorised := storage.AuthorizeSysAdmin(r.Context(), sessionCookie, csrfCookie)
		if !isAuthorised {
			slog.Info("sessionCokie & CSRFToken mismatch", "error:", "unauthorized Access")
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
			return
		}
		const sysAdminIDKey Types.SysAdminKey = "sysAdminID"

		ctx := context.WithValue(r.Context(), sysAdminIDKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// func AuthorizerStudent(storage Storage.SysAdmin, next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		sessionToken, err := r.Cookie("session_token")
// 		if err != nil {
// 			slog.Info("Session Cookie Missing", "ok:", sessionToken)
// 			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
// 			return
// 		}
// 		sessionCookie := sessionToken.Value
// 		slog.Info("Session Cookie Recieved", "ok:", sessionCookie)

// 		csrfCookie := r.Header.Get("X-CSRF-Token")
// 		if csrfCookie == "" {
// 			slog.Info("CSRF Cookie Missing", "error:", csrfCookie)
// 			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
// 			return
// 		}
// 		slog.Info("CSRF Cookie Recieved", "ok:", csrfCookie)

// 		id, isAuthorised := storage.AuthorizeSysAdmin(r.Context(), sessionCookie, csrfCookie)
// 		if !isAuthorised {
// 			slog.Info("sessionCokie & CSRFToken mismatch", "error:", "unauthorized Access")
// 			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(errors.New("unauthorized Access")))
// 			return
// 		}
// 		const sysAdminIDKey Types.SysAdminKey = "sysAdminID"

// 		ctx := context.WithValue(r.Context(), sysAdminIDKey, id)

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	}
// }
