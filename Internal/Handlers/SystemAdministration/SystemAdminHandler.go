package SysAdminHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Email"
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

func SendInvite(sysadminStore Storage.SysAdmin, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var link string
		token := r.PathValue("token")
		berr := json.NewDecoder(r.Body).Decode(&link)

		if berr != nil {
			slog.Info("Body error", "message", "Body could not be converted to  json")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("body could not be converted to  json")))
			return
		}

		// fmt.Println(token, link)
		if token == "" {
			slog.Info("Token error", "message", "token is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the token is missing")))
			return
		}

		if link == "" {
			slog.Info("Link error", "message", "Link is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the Link is missing")))
			return
		}

		email, qerr := sysadminStore.GetInviteData(r.Context(), token)
		if qerr != nil {
			slog.Info("Token error", "error", qerr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(qerr))
		}

		message := fmt.Sprintf("Peace and Blessings be upon you. Here is your link to access greatness %s", link)
		smtperr := smtp.Send(email, "Invitation", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the Link is missing")))
			return
		}
		response.WriteJson(w, http.StatusCreated, response.GeneralSuccess("the link has been sent succcesfully"))
	}
}

func GetInvitesAnalytics(sysadminStore Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := sysadminStore.GetInvitesAnalytics(r.Context())
		if err != nil {
			slog.Info("Analytics Error", "message", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, data)
	}
}

func GetInvitesApplications(sysAdminStore Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("runing")
		limitStr := r.URL.Query().Get("limit")
		offLimitStr := r.URL.Query().Get("offlimit")

		limit := 10
		offlimit := 0

		lim, err := strconv.Atoi(limitStr)
		if err == nil && lim > 0 {
			limit = lim
		}

		olimit, err := strconv.Atoi(offLimitStr)
		if err == nil && olimit > 0 && olimit <= 100 {
			offlimit = olimit
		}

		data, dberr := sysAdminStore.GetInvites(r.Context(), limit, offlimit)

		if dberr != nil {
			slog.Info("Error in Db operation", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, data)

	}
}
