package UniversityHandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func Login(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request to Login University")
		if r.Header.Get("Content-Type") != "application/json" {
			response.WriteJson(w, http.StatusUnsupportedMediaType, map[string]string{
				"error": "Content-Type must be application/json",
			})
			return
		}
		fmt.Println(r.Body)
		var universityLogin Types.UniversityLogin
		err := json.NewDecoder(r.Body).Decode(&universityLogin)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		sessiontoken, csrftoken, universityAuth, loginerr := storage.UniversityLogin(r.Context(), universityLogin)
		if loginerr != nil {
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(loginerr))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    sessiontoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
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

		response.WriteJson(w, http.StatusAccepted, universityAuth)

	}
}

func GetStudntsApplications(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		university_id := r.URL.Query().Get("university_id")
		status := r.URL.Query().Get("status")

		if university_id == "" || status == "" {
			slog.Error("requested invalid url", "error", "missing parameteers")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		studentdDetailList, dberr := storage.GetStudntsApplications(r.Context(), university_id, status)
		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusBadRequest, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, studentdDetailList)
	}
}
func RespondStudntsApplications(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		application_id := r.URL.Query().Get("application_id")
		status := r.URL.Query().Get("status")

		if application_id == "" || status == "" {
			slog.Error("requested invalid url", "error", "missing parameteers")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		dberr := storage.RespondStudntsApplications(r.Context(), application_id, status)
		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusBadRequest, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, "status update succesfully")
	}
}

func GetUniversityProgramsList(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		university_id := r.URL.Query().Get("university_id")

		if university_id == "" {
			slog.Error("requested invalid url", "error", "missing parameteers")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		programsList, dberr := storage.GetUniversityProgramsList(r.Context(), university_id)

		if dberr != nil {
			slog.Error("database error ", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, programsList)
	}
}
func GetProgramDetails(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		program_id := r.URL.Query().Get("program_id")

		if program_id == "" {
			slog.Error("requested invalid url", "error", "missing parameteers")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		programDetails, dberr := storage.GetProgramDetails(r.Context(), program_id)

		if dberr != nil {
			slog.Error("database error ", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, programDetails)
	}
}

func AddNewProgram(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var programDetails Types.Programe

		decrr := json.NewDecoder(r.Body).Decode(&programDetails)

		if decrr != nil {
			slog.Error("missing body", "error", "missing body")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decrr))
			return
		}

		dberr := storage.AddNewProgram(r.Context(), programDetails)

		if dberr != nil {
			slog.Error("database error ", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Program Added succesfully")
	}
}
func UpdateProgram(storage Storage.UniversityPortal) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		program_id := r.URL.Query().Get("program_id")

		if program_id == "" {
			slog.Error("requested invalid url", "error", "missing parameteers")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}
		var programDetails Types.Programe

		decrr := json.NewDecoder(r.Body).Decode(&programDetails)

		if decrr != nil {
			slog.Error("missing body", "error", "missing body")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decrr))
			return
		}

		dberr := storage.UpdateProgram(r.Context(), programDetails, program_id)

		if dberr != nil {
			slog.Error("database error ", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Program Updated succesfully")
	}
}
