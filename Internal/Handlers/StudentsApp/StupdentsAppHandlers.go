package StudentAppHandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Email"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	Helper "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Helpers"
	response "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Responses"
)

func StudentSignup(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var studentData Types.StudentsSignupData
		err := json.NewDecoder(r.Body).Decode(&studentData)

		if err != nil {
			slog.Error("error in decoding", "Error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		passportData, _ := Helper.Base64ToBytes(studentData.Passport)
		transcriptData, _ := Helper.Base64ToBytes(studentData.Transcript)
		var studentsDocuments Types.StudentSignupDocuments
		studentsDocuments.Passport = passportData
		studentsDocuments.Transcript = transcriptData

		dberr := storage.StudentSignup(r.Context(), studentData, studentsDocuments)

		if dberr != nil {
			slog.Error("Error in saving data to db ", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusAccepted, "Submission Sent Succesfully")
	}
}

func StudentSignin(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Student Sign in Requested")
		var loginDetails Types.StudentsSignIn
		readerr := json.NewDecoder(r.Body).Decode(&loginDetails)
		if readerr != nil {
			slog.Error("there was an error in the requested Body", "error", readerr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(readerr))
			return
		}

		sessiontoken, _, studentAuth, loginerr := storage.StudentsLogin(r.Context(), loginDetails)
		if loginerr != nil {
			slog.Error("error in login ", "error", loginerr)
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
			// Domain:   "www.pigeos.com",
		})

		// http.SetCookie(w, &http.Cookie{
		// 	Name:     "csrf_token",
		// 	Value:    csrftoken,
		// 	Expires:  time.Now().Add(24 * time.Hour),
		// 	HttpOnly: false,
		// 	Secure:   true,
		// 	Path:     "/",
		// 	SameSite: http.SameSiteNoneMode,
		// 	// Domain:   "www.pigeos.com",
		// })

		response.WriteJson(w, http.StatusAccepted, studentAuth)
	}
}

func GetCountriesList(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countryId := r.URL.Query().Get("countryId")
		countries, err := storage.GetCountryList(r.Context(), countryId)

		if err != nil {
			slog.Error("db error", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, countries)
	}
}
func GetUniversitiesList(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countryId := r.URL.Query().Get("countryId")
		countries, err := storage.GetUniversitiesList(r.Context(), countryId)
		fmt.Println(countryId)
		if err != nil {
			slog.Error("db error", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, countries)
	}
}

func GetUniversityProfile(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		university_id := r.URL.Query().Get("university_id")

		if university_id == "" {
			slog.Error("invalid url request", "error", "missing query parameters")
			response.WriteJson(w, http.StatusBadRequest, "invalid url request")
			return
		}

		uniprofile, dberr := storage.GetUniversityProfile(r.Context(), university_id)

		if dberr != nil {
			slog.Error("db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, uniprofile)
	}
}

func GetUniversityPrograms(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		program_id := r.URL.Query().Get("program_id")

		if program_id == "" {
			slog.Error("RequestedURL error", "Error", "There is a missing parameter")
			response.WriteJson(w, http.StatusBadRequest, "There is a missing Parameter")
			return
		}

		programs, err := storage.GetUniversityPrograms(r.Context(), program_id)

		if err != nil {
			slog.Error("Db error", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, programs)
	}
}

func GetStudentProfileDetails(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			slog.Error("RequestedURL error", "Error", "There is a missing parameter")
			response.WriteJson(w, http.StatusBadRequest, "There is a missing Parameter")
			return
		}

		studentDetails, err := storage.GetStudentProfileDetails(r.Context(), student_id)

		if err != nil {
			slog.Error("Db error", "Error", err)
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, studentDetails)
	}
}

func UpdateStudentProfileDetails(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detailType := r.URL.Query().Get("detail_type")
		fieldName := r.URL.Query().Get("field_name")
		fieldValue := r.URL.Query().Get("field_value")
		student_id := r.URL.Query().Get("student_id")

		// Validate required params
		// if student_id == "" || fieldName == "" || fieldValue == "" {
		// 	response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("missing required params")))
		// 	return
		// }

		var err error
		switch detailType {
		case "profile":
			err = storage.UpdateProfile(r.Context(), student_id, fieldName, fieldValue)
		case "contact":
			err = storage.UpdateContact(r.Context(), student_id, fieldName, fieldValue)
		case "education":
			err = storage.UpdateEducation(r.Context(), student_id, fieldName, fieldValue)
		case "preferences": // Fixed typo
			err = storage.UpdatePreferences(r.Context(), student_id, fieldName, fieldValue)
		default:
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid detail_type: %s", detailType)))
			return
		}

		if err != nil {
			slog.Error("Update failed", "detail_type", detailType, "student_id", student_id, "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "Field updated successfully"})
	}
}

func GetstudentsDocuments(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			slog.Error("requested invalid url", "error", "missing parameter")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		documentsList, err := storage.GetstudentsDocuments(r.Context(), student_id)

		if err != nil {
			slog.Error("db error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusAccepted, documentsList)
	}
}

func UploadStudentDocuments(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		var uploadedDocuments Types.UploadDocument

		if student_id == "" {
			slog.Error("requested invalid url", "error", "missing parameter")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		err := json.NewDecoder(r.Body).Decode(&uploadedDocuments)
		if err != nil {
			slog.Error("encoding error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}

		documentDataBytes, _ := Helper.Base64ToBytes(uploadedDocuments.Data)

		qerr := storage.UploadStudentDocuments(r.Context(), uploadedDocuments, documentDataBytes, student_id)

		if qerr != nil {
			slog.Error("db error", "error", qerr)
			response.WriteJson(w, http.StatusBadRequest, qerr)
			return
		}
		response.WriteJson(w, http.StatusOK, "document uploaded")
	}
}

func GetStudentsDocument(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		document_id := r.URL.Query().Get("document_id")
		studentId := r.URL.Query().Get("student_id")

		if studentId == "" || document_id == "" {
			slog.Error("error in requested url", "error", "missing parameters")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		document, err := storage.GetStudentsDocument(r.Context(), studentId, document_id)

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

func UploadApplicationReceipt(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Document Types.UploadDocument
		student_id := r.URL.Query().Get("student_id")
		university_id := r.URL.Query().Get("university_id")
		program_id := r.URL.Query().Get("program_id")
		paid_amount := r.URL.Query().Get("paid_amount")

		if university_id == "" || student_id == "" || program_id == "" || paid_amount == "" {
			slog.Error("Invalid Ul Request", "error", fmt.Sprintf("missing Query Parameter studnet_id = %s and university_id = %s", student_id, university_id))
			response.WriteJson(w, http.StatusBadRequest, "missing parameters")
			return
		}

		err := json.NewDecoder(r.Body).Decode(&Document)

		if err != nil {
			slog.Error("Body decode error", "error", err)
			response.WriteJson(w, http.StatusBadRequest, err)
			return
		}

		dberr := storage.UploadApplicationReceipt(r.Context(), student_id, university_id, program_id, paid_amount, Document)

		if dberr != nil {
			slog.Error("db operation error", "error", dberr)
			// response.WriteJson(w, http.StatusBadRequest, dberr)
			return
		}

		response.WriteJson(w, http.StatusCreated, "file upload successfull")
	}
}

func ApplyToUniversity(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		university_id := r.URL.Query().Get("university_id")
		program_id := r.URL.Query().Get("program_id")

		if university_id == "" || student_id == "" || program_id == "" {
			slog.Error("Invalid Ul Request", "error", fmt.Sprintf("missing Query Parameter studnet_id = %s and university_id = %s", student_id, university_id))
			response.WriteJson(w, http.StatusBadRequest, "missing parameters")
			return
		}

		dberr := storage.ApplyToUniversity(r.Context(), student_id, university_id, program_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, "Application Submitted Succesfully")
	}
}

func GetApplicationsData(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			slog.Error("Invalid Ul Request", "error", fmt.Sprintf("missing Query Parameter studnet_id = %s ", student_id))
			response.WriteJson(w, http.StatusBadRequest, "missing parameters")
			return
		}

		applications, dberr := storage.GetApplicationsData(r.Context(), student_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusCreated, applications)

	}
}

func VerifyApplication(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		program_id := r.URL.Query().Get("program_id")
		student_id := r.URL.Query().Get("student_id")
		university_id := r.URL.Query().Get("university_id")

		if student_id == "" || program_id == "" || university_id == "" {

			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		existingrows, dberr := storage.VerifyApplication(r.Context(), student_id, program_id, university_id)
		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}
		response.WriteJson(w, http.StatusOK, existingrows)
	}
}

func ShortListProgram(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		program_id := r.URL.Query().Get("program_id")
		student_id := r.URL.Query().Get("student_id")
		university_id := r.URL.Query().Get("university_id")

		if program_id == "" || student_id == "" || university_id == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		id, dberr := storage.ShortListProgram(r.Context(), student_id, program_id, university_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, id)

	}
}

func GetShortListProgram(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		listOfPrograms, dberr := storage.GetShortListProgram(r.Context(), student_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, listOfPrograms)
	}
}
func DeleteShortListProgram(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		shortList_id := r.URL.Query().Get("shortlist_id")

		if student_id == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		dberr := storage.DeleteShortListProgram(r.Context(), student_id, shortList_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "program deleted from shortlist")
	}
}

func RegisterForEvent(storage Storage.StudentsApp, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		webinar_code := r.URL.Query().Get("webinar_code")
		var message string

		if student_id == "" || webinar_code == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		student_email, webinar, dberr := storage.RegisterForEvent(r.Context(), student_id, webinar_code)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		message = fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Event Registration Confirmed – GEOS</title></head>
		<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
		<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Your spot is confirmed! Here's your joining link and calendar invite for the UK Virtual University Fair.</span>
		<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
		<tr><td align="center">
			<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
			<tr><td style="height:4px;background:linear-gradient(to right,#5b3fcc,#7c5fe6);font-size:0;">&nbsp;</td></tr>
			<tr><td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
				<table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr>
				<td><table cellpadding="0" cellspacing="0" border="0"><tr>
					<td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span></td>
					<td style="padding-left:7px;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span></td>
				</tr></table></td>
				<td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f8f5ff;color:#5b3fcc;border:1px solid #ddd6fe;text-transform:uppercase;letter-spacing:0.5px;">Registration Confirmed</span></td>
				</tr></table>
			</td></tr>
			<tr><td style="background:#ffffff;padding:28px 32px;">
				<p style="margin:0 0 4px;font-size:32px;line-height:1;">🎓</p>
				<p style="margin:8px 0 6px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">You're registered!</p>
				<h1 style="margin:0 0 14px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">See you at the %s</h1>
				<p style="margin:0 0 18px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">Your spot is confirmed. Below you'll find your joining link, event details, and calendar buttons to save the date.</p>
				<!-- Event details card -->
				<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;border:1px solid #eef1f8;border-radius:9px;overflow:hidden;">
				<tr><td style="background:#131c36;padding:10px 14px;"><span style="color:#ffffff;font-size:12.5px;font-weight:600;">📅 Event Details</span></td></tr>
				<tr><td style="padding:14px 16px;">
					<table width="100%%" cellpadding="0" cellspacing="0" border="0">
					<tr><td style="padding:5px 0;border-bottom:1px solid #f0f2f8;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr><td style="width:68px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Event</td><td style="font-size:13px;color:#131c36;font-weight:600;">%s</td></tr></table></td></tr>
					<tr><td style="padding:5px 0;border-bottom:1px solid #f0f2f8;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr><td style="width:68px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Date</td><td style="font-size:13px;color:#131c36;font-weight:500;">%s</td></tr></table></td></tr>
					<tr><td style="padding:5px 0;border-bottom:1px solid #f0f2f8;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr><td style="width:68px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Time</td><td style="font-size:13px;color:#131c36;font-weight:500;">%s (GMT)</td></tr></table></td></tr>
					<tr><td style="padding:5px 0;border-bottom:1px solid #f0f2f8;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr><td style="width:68px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Format</td><td style="font-size:13px;color:#131c36;font-weight:500;">%s</td></tr></table></td></tr>
					<tr><td style="padding:5px 0;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr><td style="width:68px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Host</td><td style="font-size:13px;color:#131c36;font-weight:500;">%s</td></tr></table></td></tr>
					</table>
				</td></tr>
				</table>
				<!-- Joining link -->
				<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;background:#f8f5ff;border-left:3px solid #5b3fcc;border-radius:0 9px 9px 0;"><tr><td style="padding:14px 16px;">
				<p style="margin:0 0 5px;font-size:11.5px;font-weight:700;color:#5b3fcc;">Your joining link</p>
				<a href="#" style="font-size:12.5px;color:#5b3fcc;word-break:break-all;font-weight:600;text-decoration:none;">%s</a>
				<p style="margin:6px 0 0;font-size:11.5px;color:#94a3b8;">Webinar ID: %s &nbsp;</p>
				</td></tr></table>
				
				<table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f8f5ff;border:1px solid #ddd6fe;border-radius:9px;padding:12px 14px;">
				<p style="margin:0;font-size:12.5px;color:#6b7280;line-height:1.7;">📌 <strong style="color:#131c36;">Can't make it?</strong> A recording will be sent to all registered attendees within 48 hours of the event.</p>
				</td></tr></table>
			</td></tr>
			<tr><td style="background:#f8f9fc;padding:16px 32px;border-top:1px solid #f0f2f8;text-align:center;">
				<p style="margin:0 0 7px;"><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Privacy Policy</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Help Centre</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Cancel Registration</a></p>
				<p style="margin:0;font-size:10.5px;color:#c0c8d8;">© 2026 GEOS Ltd · London, UK</p>
			</td></tr>
			</table>
		</td></tr>
		</table>
		</body>
		</html>
		`,
			webinar.Title,
			webinar.Title,
			webinar.Date,
			webinar.Time,
			webinar.Platform,
			webinar.Host,
			webinar.Link,
			webinar_code,
		)

		smtperr := smtp.Send(student_email, "Event Registration Confirmed", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(smtperr))
			return
		}
		response.WriteJson(w, http.StatusOK, "Student Registered For Event Succesfully")

	}
}

func EventRegisterationCheck(storage Storage.StudentsApp, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		webinar_code := r.URL.Query().Get("webinar_code")

		if student_id == "" || webinar_code == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		registered, dberr := storage.EventRegisterationCheck(r.Context(), student_id, webinar_code)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, registered)

	}

}

func SetScholarshipReminder(storage Storage.StudentsApp, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		scholarship_id := r.URL.Query().Get("scholarship_id")
		var message string

		if student_id == "" || scholarship_id == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		student_email, scholarship_name, opens_date, dberr := storage.SetScholarshipReminder(r.Context(), student_id, scholarship_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		message = fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Scholarship Reminder Set – GEOS</title>
		</head>
		<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
		<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">
			Your scholarship reminder has been set. You’ll receive an email notification when it opens.
		</span>

		<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
			<tr>
			<td align="center">
				<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
				<tr>
					<td style="height:4px;background:linear-gradient(to right,#10b981,#059669);font-size:0;">&nbsp;</td>
				</tr>

				<tr>
					<td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
					<table width="100%%" cellpadding="0" cellspacing="0" border="0">
						<tr>
						<td>
							<table cellpadding="0" cellspacing="0" border="0">
							<tr>
								<td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;">
								<span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span>
								</td>
								<td style="padding-left:7px;vertical-align:middle;">
								<span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span>
								</td>
							</tr>
							</table>
						</td>
						<td style="text-align:right;">
							<span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f0fdf8;color:#059669;border:1px solid #a7f3d0;text-transform:uppercase;letter-spacing:0.5px;">
							Reminder Set ✓
							</span>
						</td>
						</tr>
					</table>
					</td>
				</tr>

				<tr>
					<td style="background:#ffffff;padding:28px 32px;">
					<p style="margin:0 0 4px;font-size:32px;line-height:1;">🔔</p>
					<p style="margin:8px 0 6px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">
						Scholarship alert
					</p>
					<h1 style="margin:0 0 14px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">
						Your reminder has been set
					</h1>
					<p style="margin:0 0 18px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">
						We’ve saved your scholarship reminder. When the scholarship opens, you’ll receive an email notification right away.
					</p>

					<!-- Reminder details card -->
					<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;border:1px solid #eef1f8;border-radius:9px;overflow:hidden;">
						<tr>
						<td style="background:#131c36;padding:10px 14px;">
							<span style="color:#ffffff;font-size:12.5px;font-weight:600;">📌 Reminder Details</span>
						</td>
						</tr>
						<tr>
						<td style="padding:14px 16px;">
							<table width="100%%" cellpadding="0" cellspacing="0" border="0">
							<tr>
								<td style="padding:5px 0;border-bottom:1px solid #f0f2f8;">
								<table cellpadding="0" cellspacing="0" border="0" width="100%%">
									<tr>
									<td style="width:88px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Scholarship</td>
									<td style="font-size:13px;color:#131c36;font-weight:600;">%s</td>
									</tr>
								</table>
								</td>
							</tr>
							<tr>
								<td style="padding:5px 0;border-bottom:1px solid #f0f2f8;">
								<table cellpadding="0" cellspacing="0" border="0" width="100%%">
									<tr>
									<td style="width:88px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Opens On</td>
									<td style="font-size:13px;color:#131c36;font-weight:500;">%s</td>
									</tr>
								</table>
								</td>
							</tr>
							<tr>
								<td style="padding:5px 0;border-bottom:1px solid #f0f2f8;">
								<table cellpadding="0" cellspacing="0" border="0" width="100%%">
									<tr>
									<td style="width:88px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Reminder</td>
									<td style="font-size:13px;color:#131c36;font-weight:500;">Email notification enabled</td>
									</tr>
								</table>
								</td>
							</tr>
							<tr>
								<td style="padding:5px 0;">
								<table cellpadding="0" cellspacing="0" border="0" width="100%%">
									<tr>
									<td style="width:88px;font-size:10px;color:#94a3b8;text-transform:uppercase;letter-spacing:0.5px;padding-top:1px;">Status</td>
									<td style="font-size:13px;color:#131c36;font-weight:500;">Ready</td>
									</tr>
								</table>
								</td>
							</tr>
							</table>
						</td>
						</tr>
					</table>

					<!-- Highlight box -->
					<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;background:#f0fdf8;border-left:3px solid #059669;border-radius:0 9px 9px 0;">
						<tr>
						<td style="padding:14px 16px;">
							<p style="margin:0;font-size:12.5px;color:#4b5a7a;line-height:1.7;">
							✅ <strong style="color:#131c36;">Youre all set.</strong> We ll send you an email as soon as this scholarship opens so you can apply without missing the deadline.
							</p>
						</td>
						</tr>
					</table>

				

					<table width="100%%" cellpadding="0" cellspacing="0" border="0">
						<tr>
						<td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td>
						</tr>
					</table>

					<p style="margin:16px 0 0;text-align:center;font-size:13px;color:#4b5a7a;">
						Questions? We re here — <a href="mailto:hello@geosedtech.com" style="color:#5b3fcc;font-weight:600;text-decoration:none;">hello@geosedtech.com</a>
					</p>
					</td>
				</tr>

				<tr>
					<td style="background:#f8f9fc;padding:16px 32px;border-top:1px solid #f0f2f8;text-align:center;">
					<p style="margin:0 0 7px;">
						<a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Privacy Policy</a>
						<a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Help Centre</a>
						<a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Manage Reminders</a>
					</p>
					<p style="margin:0;font-size:10.5px;color:#c0c8d8;">© 2026 GEOS Ltd · London, UK</p>
					</td>
				</tr>
				</table>
			</td>
			</tr>
		</table>
		</body>
		</html>
		`,
			scholarship_name,
			opens_date,
		)

		smtperr := smtp.Send(student_email, "Scholarship Reminder Set", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(smtperr))
			return
		}
		response.WriteJson(w, http.StatusOK, "Student Reminder Set For Scholarship Succesfully")

	}
}

func ScholarshipReminderCheck(storage Storage.StudentsApp, smtp Email.EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student_id := r.URL.Query().Get("student_id")
		scholarship_id := r.URL.Query().Get("scholarship_id")

		if student_id == "" || scholarship_id == "" {
			slog.Error("invalid URL requested", "error", "missing query parameter")
			response.WriteJson(w, http.StatusBadRequest, "invalid url requested")
			return
		}

		registered, dberr := storage.ScholarshipReminderCheck(r.Context(), student_id, scholarship_id)

		if dberr != nil {
			slog.Error("Db error", "erorr", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, registered)

	}

}

func GetFreeApplicationCount(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		student_id := r.URL.Query().Get("student_id")

		if student_id == "" {
			slog.Error("invalid url", "error", "Invalid URL")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL Requested")
			return
		}

		freeApplicationCoiunt, dberr := storage.GetFreeApplicationCount(r.Context(), student_id)

		if dberr != nil {
			slog.Error("db error", "error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, freeApplicationCoiunt)
	}
}

func SearchPrograms(storage Storage.StudentsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search_term := r.URL.Query().Get("search_term")

		country_ids, dberr := storage.SearchPrograms(r.Context(), search_term)

		if dberr != nil {
			slog.Error("db error", "error", dberr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, country_ids)
	}
}
