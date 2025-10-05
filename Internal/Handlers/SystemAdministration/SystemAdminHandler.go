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

func Login(storage Storage.SysAdmin) http.HandlerFunc {
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
		sessiontoken, csrftoken, sysAdminAuth, loginerr := storage.SysAdminLogin(r.Context(), SysAdmin)
		if loginerr != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(loginerr))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessiontoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,  // decides if it can be read by the browser
			Secure:   false, // decides if it should be sent on http request or https only
			Path:     "/",
			// SameSite: http.SameSiteNoneMode, // allow cross-origin read/write this requires secure true
			// Path: "/sysadmin",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false, // decides if it can be read by the browser
			Secure:   false, // decides if it shouldould be sent on http req or https onlyuest
			Path:     "/",
			// SameSite: http.SameSiteNoneMode, // Lax works without HTTPS
			// Path:     "/sysadmin",          // works for all paths
		})

		response.WriteJson(w, http.StatusAccepted, sysAdminAuth)

	}
}

func Signup(storage Storage.SysAdmin) http.HandlerFunc {
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
			slog.Info("There was an error in the request body validation", "error:", validationerr)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}
		dberr := storage.SysAdminSignup(r.Context(), SysAdmin)
		if dberr != nil {
			slog.Info("There was an error in saving the Sys Admin to db", "error:", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}
		slog.Info("Sys Admin Created")

		response.WriteJson(w, http.StatusCreated, "User Registered Successsfuly")
	}
}

func CreateInvite(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var schoolInvite Types.SchoolInvite
		//get theID
		const sysAdminIDKey Types.SysAdminKey = "sysAdminID"
		sysAdminId := Types.SysAdminId(r.Context().Value(sysAdminIDKey).(string))

		//get the body
		err := json.NewDecoder(r.Body).Decode(&schoolInvite)
		if err != nil {
			slog.Info("There was an error in the body sent", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validate the body
		verr := validator.New().Struct(schoolInvite)
		if verr != nil {
			validateErrors := verr.(validator.ValidationErrors)
			slog.Info("There was an error in the body sent", "error", validateErrors)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(validateErrors))
			return
		}

		// call the db
		generatedData, terr := storage.GenerateInvite(r.Context(), sysAdminId, schoolInvite)
		if terr != nil {
			slog.Info("There was an error in Generating Token", "error", terr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(terr))
			return
		}

		//return response
		response.WriteJson(w, http.StatusCreated, generatedData)
	}
}
