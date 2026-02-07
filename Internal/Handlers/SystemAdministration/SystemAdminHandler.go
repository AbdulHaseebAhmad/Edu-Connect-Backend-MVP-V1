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
		slog.Info("Request to Login system Admin")
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
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Domain:   ".pigeos.com",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrftoken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Domain:   ".pigeos.com",
		})

		response.WriteJson(w, http.StatusAccepted, sysAdminAuth)

	}
}

func Signup(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creatingsystem Admin")
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
		invitation_id := r.PathValue("invitation_id")
		berr := json.NewDecoder(r.Body).Decode(&link)

		if berr != nil {
			slog.Info("Body error", "message", "Body could not be converted to  json")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("body could not be converted to  json")))
			return
		}

		// fmt.Println(invitation_id, link)
		if invitation_id == "" {
			slog.Info("Token error", "message", "invitation_id is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the invitation_id is missing")))
			return
		}

		if link == "" {
			slog.Info("Link error", "message", "Link is missing")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("the Link is missing")))
			return
		}

		email, qerr := sysadminStore.GetInviteData(r.Context(), invitation_id)
		if qerr != nil {
			slog.Info("Token error", "error", qerr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(qerr))
		}

		message := fmt.Sprintf("Peace and Blessings be upon you. Here is your link to access greatness %s", link)
		smtperr := smtp.Send(email, "Invitation", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(smtperr))
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

func GetSchoolApplications(sysAdminStore Storage.SysAdmin) http.HandlerFunc {
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

		listOfApplications, dberr := sysAdminStore.GetSchoolApplications(r.Context(), limit, offlimit)

		if dberr != nil {
			slog.Info("Error in Db operation", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, listOfApplications)

	}
}
func GetSchoolApplication(sysAdminStore Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		application_id := r.URL.Query().Get("application_id")

		if application_id == "" {
			slog.Info("invalid url request", "error", "missing query parameter application_id")
			response.WriteJson(w, http.StatusInternalServerError, "invaalid url request")
			return
		}
		listOfApplications, dberr := sysAdminStore.GetSchoolApplication(r.Context(), application_id)

		if dberr != nil {
			slog.Info("Error in Db operation", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, listOfApplications)

	}
}

func GetAllInvites(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allInvitesSent, dberr := storage.GetAllInvites(r.Context())

		if dberr != nil {
			slog.Info("Error in Db operation", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, dberr)
			return
		}

		response.WriteJson(w, http.StatusOK, allInvitesSent)

	}
}

func RespondToSchoolApplication(sysAdminStore Storage.SysAdmin, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		application_id := r.URL.Query().Get("application_id")
		status := r.URL.Query().Get("status")
		var message string

		if application_id == "" || status == "" {
			slog.Info("Query Error", "message", "appid or status is missing", "appid= ", application_id, "status= ", status)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("query Error")))
			return
		}

		schoolInfo, generatePassword, err := sysAdminStore.RespondToSchoolInvite(r.Context(), application_id, status)
		if err != nil {
			slog.Info("invite db  Error", "message", "There was an error accepting or rejecting invite", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if status == "approved" {
			message = fmt.Sprintf("Peace and Blessings be upon you. Congratulations Your application id Number: %s has been accepted. Here is your email & password to access account. Email: %s /n Password: %s. below is attached your school code to share with your students on sign up %s", schoolInfo.SchoolId, schoolInfo.Sys_Eamil, generatePassword, schoolInfo.Code)
		} else {
			message = "Peace and Blessings be upon you. Your application to join the system was rejected"
		}

		smtperr := smtp.Send(schoolInfo.Email, "Response", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(smtperr))
			return
		}
		response.WriteJson(w, http.StatusCreated, response.GeneralSuccess("Response sent succesfully"))

	}
}

func GetAnalyticalLists(sysAdminStore Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := sysAdminStore.GetAnalyticsList(r.Context())
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, list)
	}
}

func GetStudentsRegistry(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := r.URL.Query().Get("status")
		if status == "" {
			slog.Error("The status param is missing", "error", "the student status is missing")
			response.WriteJson(w, http.StatusBadRequest, "The status is missing")
			return
		}

		studentsList, err := storage.GetStudentsRegistry(r.Context(), status)

		if err != nil {
			slog.Error("Error in Db", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, studentsList)
	}
}

func RespondApplication(storage Storage.SysAdmin, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		id := r.URL.Query().Get("id")

		if action == "" || id == "" {
			slog.Error("Missing param", "error", "there is a missing param in url")
			response.WriteJson(w, http.StatusBadRequest, "Incomplete Action")
			return
		}

		fmt.Println(action, id)
		email, password, err := storage.RespondApplication(r.Context(), action, id)

		if err != nil {
			slog.Error("there was an error in db operation", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		message := ""
		if action == "approved" {
			message = fmt.Sprintf("Peace and Blessings be upon you. Congratulations Your application has been accepted. Here is your email & password to access account. Email: %s /n Password: %s", email, password)
		} else {
			message = "Peace and Blessings be upon you. Your application to join the system was rejected"

		}

		smtp.Send(email, "Acceptance of Application", message)
		response.WriteJson(w, http.StatusOK, "The Operation was succesfull")
	}
}

func GetStudentDocuments(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentName := r.URL.Query().Get("docname")
		studentId := r.URL.Query().Get("studentId")
		documentmime := r.URL.Query().Get("docmime")

		if studentId == "" || documentName == "" {
			slog.Error("error in requested url", "error", "missing parameters")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		document, err := storage.GetStudentsDocument(r.Context(), studentId, documentName, documentmime)

		if err != nil {
			slog.Error("Error in DB", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		w.Header().Set("Content-Type", document.MimeType)
		w.Header().Set("Content-Disposition", `inline; filename="requested-document"`)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(document.Data)

	}
}

func GetReceipts(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			receiptsList, err := storage.GetAllReceipts(r.Context())

			if err != nil {
				slog.Error("Db error", "error", err)
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}

			response.WriteJson(w, http.StatusOK, receiptsList)
			return
		}
		receipt, err := storage.GetReceiptDetails(r.Context(), student_id)

		if err != nil {
			slog.Error("Db error", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, receipt)

	}
}

func RespondToReceipts(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		receipt_id := r.URL.Query().Get("receipt_id")
		status := r.URL.Query().Get("status")

		if status == "" || receipt_id == "" {
			slog.Error("error in requested url", "error", "missing parameters")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		dberr := storage.RespondToReceipts(r.Context(), receipt_id, status)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "receipt status updated")

	}
}

func GetRegisteredStudents(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registeredStudents, dberr := storage.GetRegisteredStudents(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, registeredStudents)

	}
}
