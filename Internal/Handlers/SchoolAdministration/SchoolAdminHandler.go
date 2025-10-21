package SchoolAdminHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
	"github.com/go-playground/validator/v10"
)

func Login(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Login School")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}
		fmt.Println(r.Body)
		var schoolAdmin Types.SchoolAdminLogin
		err := json.NewDecoder(r.Body).Decode(&schoolAdmin)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		sessiontoken, csrftoken, sysAdminAuth, loginerr := storage.SchoolAdminLogin(r.Context(), schoolAdmin)
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
			// Path: "/schoolAdmin",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false, // decides if it can be read by the browser
			Secure:   false, // decides if it shouldould be sent on http req or https onlyuest
			Path:     "/",
			// SameSite: http.SameSiteNoneMode, // Lax works without HTTPS
			// Path:     "/schoolAdmin",          // works for all paths
		})

		response.WriteJson(w, http.StatusAccepted, sysAdminAuth)

	}
}
func LinkValidation(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			slog.Info("Token Error", "error", "token is missing in request")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("token is invalid")))
			return
		}

		status, err := storage.ValidateLink(r.Context(), token)
		if err != nil {
			slog.Info("Token Error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusAccepted, status)
	}
}

func SubmitInviteData(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var SchoolInformation Types.SchoolInformation

		//get the token from path,
		token := r.PathValue("token")
		if token == "" {
			slog.Info("Token Error", "error", "token is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the id to the form is missing")))
			return
		}
		//validate it
		status, err := storage.ValidateLink(r.Context(), token)
		if err != nil {
			slog.Info("Token Error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if status != "pending" {
			slog.Info("Token Error", "error", "token is already consumed")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("token is already consumed")))
			return

		}

		//if not complete
		derr := json.NewDecoder(r.Body).Decode(&SchoolInformation)
		if derr != nil {
			slog.Info("Decode Error", "error", "Couldnt decode body")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("body is missing/invalid")))
			return
		}
		verr := validator.New().Struct(SchoolInformation)
		if verr != nil {
			validateErrors := verr.(validator.ValidationErrors)
			slog.Info("There was an error in the body sent", "error", validateErrors)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(validateErrors))
			return
		}
		//update scholInvites db with information
		// mark completed in school_invites table
		status, qerr := storage.SubmitInvite(r.Context(), SchoolInformation, token)
		if qerr != nil {
			slog.Info("Db Error", "error", qerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(qerr))
			return
		}
		// return response
		response.WriteJson(w, http.StatusCreated, status)
	}

}
