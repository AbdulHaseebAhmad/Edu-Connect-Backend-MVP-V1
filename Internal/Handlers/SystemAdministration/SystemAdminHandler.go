package SysAdminHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
	"github.com/go-playground/validator/v10"
)

func Login(storage Storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Login Student")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}

		var SysAdmin Types.SysAdminLogin
		err := json.NewDecoder(r.Body).Decode(&SysAdmin)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		sessiontoken, csrftoken, loginerr := storage.SysAdminLogin(r.Context(), SysAdmin)
		if loginerr != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(loginerr))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessiontoken,
			Expires: time.Now().Add(24 * time.Hour),

			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
		})

		response.WriteJson(w, http.StatusAccepted, response.GeneralSuccess("User Logged In succesfully"))

	}
}

func Signup(storage Storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating Sys Admin")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}
		var SysAdmin Types.SysAdminSignup
		err := json.NewDecoder(r.Body).Decode(&SysAdmin)
		if err != nil {
			slog.Info("There was an error in the request body", "error:", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		validationerr := validator.New().Struct(SysAdmin)
		if validationerr != nil {
			validateErrors := validationerr.(validator.ValidationErrors)
			slog.Info("There was an error in the request body validation", "error:", err)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}
		dberr := storage.SysAdminSignup(r.Context(), SysAdmin)
		if dberr != nil {
			slog.Info("There was an error in saving the Sys Admin to db", "error:", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		slog.Info("Sys Admin Created")

		response.WriteJson(w, http.StatusCreated, "User Registered Successsfuly")
	}
}
