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
			HttpOnly: true,
			Secure:   true, // must be true for cross-site HTTPS
			Path:     "/",
			SameSite: http.SameSiteNoneMode, // allows cross-origin
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
		})

		response.WriteJson(w, http.StatusAccepted, sysAdminAuth)

	}
}
func LinkValidation(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invitation_id := r.URL.Query().Get("invitation_id")
		if invitation_id == "" {
			slog.Info("Token Error", "error", "invitation_id is missing in request")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invitation_id is invalid")))
			return
		}

		status, err := storage.ValidateLink(r.Context(), invitation_id)
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
		invitation_id := r.PathValue("invitation_id")
		if invitation_id == "" {
			slog.Info("Token Error", "error", "invitation_id is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the id to the form is missing")))
			return
		}
		//validate it
		status, err := storage.ValidateLink(r.Context(), invitation_id)
		if err != nil {
			slog.Info("Token Error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if status != "pending" {
			slog.Info("Token Error", "error", "invitation_id is already consumed")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("invitation_id is already consumed")))
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
		status, qerr := storage.SubmitInvite(r.Context(), SchoolInformation)
		if qerr != nil {
			slog.Info("Db Error", "error", qerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(qerr))
			return
		}
		// return response
		response.WriteJson(w, http.StatusCreated, status)
	}

}

func GetUnProcessedStudentsList(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")

		if school_id == "" {
			slog.Error("invalidurl request", "error", "incomplete parameters")
			response.WriteJson(w, http.StatusBadRequest, "parameters missing")
			return
		}

		listOfStudents, dberr := storage.GetUnProcessedStudentsList(r.Context(), school_id)

		if dberr != nil {
			slog.Error("Error inside Db", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, listOfStudents)

	}
}

func VerifyStudentAccount(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("running verify student handler")
		school_id := r.URL.Query().Get("school_id")
		student_id := r.URL.Query().Get("student_id")
		status := r.URL.Query().Get("status")

		if school_id == "" || student_id == "" || status == "" {
			slog.Error("invalidurl request", "error", "incomplete parameters")
			response.WriteJson(w, http.StatusBadRequest, "parameters missing")
			return
		}

		dberr := storage.VerifyStudentAccount(r.Context(), school_id, student_id, status)

		if dberr != nil {
			slog.Error("Db Error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, "status updated")
	}
}

func GetProcessedStudentsList(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")
		status := r.URL.Query().Get("status")

		if school_id == "" {
			slog.Error("invalidurl request", "error", "incomplete parameters")
			response.WriteJson(w, http.StatusBadRequest, "parameters missing")
			return
		}

		listOfStudents, dberr := storage.GetProcessedStudentsList(r.Context(), school_id, status)

		if dberr != nil {
			slog.Error("Error inside Db", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, listOfStudents)

	}
}

func GetSchoolProfileData(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")

		if school_id == "" {
			slog.Error("invalid url request", "error", "parameters missing")
			response.WriteJson(w, http.StatusBadRequest, "Missing Params")
			return
		}

		schoolData, dberr := storage.GetSchoolProfileData(r.Context(), school_id)

		if dberr != nil {
			slog.Error("database error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, schoolData)
	}
}
