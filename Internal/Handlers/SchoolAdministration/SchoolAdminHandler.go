package SchoolAdminHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Email"
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
		sessiontoken, _, sysAdminAuth, loginerr := storage.SchoolAdminLogin(r.Context(), schoolAdmin)
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

func SubmitInviteData(storage Storage.SchoolAdmin, smtp Email.EmailSender) http.HandlerFunc {
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

		message := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>School Application Under Review – GEOS</title></head>
<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Your school's GEOS application is under review — we'll be in touch within 2–3 business days.</span>
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
  <tr><td align="center">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
      <tr><td style="height:4px;background:linear-gradient(to right,#131c36,#1e2d55);font-size:0;">&nbsp;</td></tr>
      <tr><td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr>
          <td><table cellpadding="0" cellspacing="0" border="0"><tr>
            <td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span></td>
            <td style="padding-left:7px;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span></td>
          </tr></table></td>
          <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f5f6fa;color:#131c36;border:1px solid #e4e8f4;text-transform:uppercase;letter-spacing:0.5px;">School Application</span></td>
        </tr></table>
      </td></tr>
      <tr><td style="background:#ffffff;padding:28px 32px;">
        <p style="margin:0 0 4px;font-size:32px;line-height:1;">🏫</p>
        <p style="margin:8px 0 6px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">Hello, Admissions Team 👋</p>
        <h1 style="margin:0 0 14px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">Your school application is under review</h1>
        <p style="margin:0 0 12px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">Thank you for registering <strong style="color:#131c36;">%s</strong> on the GEOS platform. We've received your application and our partnerships team is currently reviewing your institution's details.</p>
        <p style="margin:0 0 18px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">This process typically takes <strong style="color:#131c36;">2–3 business days</strong>. Once verified, you'll receive a welcome email with full portal access and your school's unique student referral code.</p>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;"><tr><td style="background:#f5f6fa;border-left:3px solid #131c36;border-radius:0 9px 9px 0;padding:13px 15px;">
          <p style="margin:0;font-size:13px;color:#4b5a7a;line-height:1.7;"><strong style="color:#131c36;">What happens during review?</strong><br>We verify your institution's registration details, contact information, and accreditation status to ensure GEOS is a great fit for your students and community.</p>
        </td></tr></table>
        <!-- Status -->
        <p style="margin:0 0 10px;font-size:11px;font-weight:700;color:#94a3b8;letter-spacing:1.5px;text-transform:uppercase;">Application Progress</p>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;">
          <tr><td style="padding:6px 0;"><table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#f8f9fc;border:1px solid #eef1f8;border-radius:9px;"><tr><td style="padding:11px 13px;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="font-size:18px;width:28px;">✅</td>
            <td style="font-size:13px;color:#131c36;font-weight:500;padding-left:10px;">Application submitted</td>
            <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:3px 9px;border-radius:6px;background:#f0fdf8;color:#059669;border:1px solid #a7f3d0;text-transform:uppercase;">Done</span></td>
          </tr></table></td></tr></table></td></tr>
          <tr><td style="padding:6px 0;"><table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#f5f6fa;border:1px solid #e4e8f4;border-radius:9px;"><tr><td style="padding:11px 13px;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="font-size:18px;width:28px;">⏳</td>
            <td style="font-size:13px;color:#131c36;font-weight:700;padding-left:10px;">Institution verification</td>
            <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:3px 9px;border-radius:6px;background:#fffbf0;color:#d97706;border:1px solid #fde68a;text-transform:uppercase;">In Progress</span></td>
          </tr></table></td></tr></table></td></tr>
          <tr><td style="padding:6px 0;opacity:0.45;"><table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#f8f9fc;border:1px solid #eef1f8;border-radius:9px;"><tr><td style="padding:11px 13px;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="font-size:18px;width:28px;">⭕</td>
            <td style="font-size:13px;color:#131c36;font-weight:500;padding-left:10px;">School portal access granted</td>
            <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:3px 9px;border-radius:6px;background:#f5f6fa;color:#94a3b8;border:1px solid #e4e8f4;text-transform:uppercase;">Waiting</span></td>
          </tr></table></td></tr></table></td></tr>
          <tr><td style="padding:6px 0;opacity:0.45;"><table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#f8f9fc;border:1px solid #eef1f8;border-radius:9px;"><tr><td style="padding:11px 13px;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="font-size:18px;width:28px;">⭕</td>
            <td style="font-size:13px;color:#131c36;font-weight:500;padding-left:10px;">Student referral code issued</td>
            <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:3px 9px;border-radius:6px;background:#f5f6fa;color:#94a3b8;border:1px solid #e4e8f4;text-transform:uppercase;">Waiting</span></td>
          </tr></table></td></tr></table></td></tr>
        </table>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:14px;"><tr><td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td></tr></table>
        <p style="margin:0;font-size:13px;color:#4b5a7a;line-height:1.8;">Questions? Reach our partnerships team at <a href="mailto:hello@geosedtech.com" style="color:#5b3fcc;font-weight:600;text-decoration:none;">hello@geosedtech.com</a> or call <strong style="color:#131c36;">+44 20 1234 5678</strong>.</p>
      </td></tr>
      <tr><td style="background:#f8f9fc;padding:16px 32px;border-top:1px solid #f0f2f8;text-align:center;">
        <p style="margin:0 0 7px;"><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Privacy Policy</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">School Help Centre</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Contact Partnerships</a></p>
        <p style="margin:0;font-size:10.5px;color:#c0c8d8;">© 2026 GEOS Ltd · London, UK</p>
      </td></tr>
    </table>
  </td></tr>
</table>
</body>
</html>
`, SchoolInformation.School)
		smtperr := smtp.Send(SchoolInformation.Email, "School Application Under Review", message)

		if smtperr != nil {
			slog.Info("SMTP error", "message", smtperr)
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(smtperr))
			return
		}
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

func GetSchoolStatistics(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")

		if school_id == "" {
			slog.Error("params rror", "error:", "School Id missing")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL Request")
			return
		}

		statistics, dberr := storage.GetSchoolStatistics(r.Context(), school_id)

		if dberr != nil {
			slog.Error("database error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, statistics)

	}
}

func GetEnrolledStudents(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")

		if school_id == "" {
			slog.Error("params rror", "error:", "School Id missing")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL Request")
			return
		}

		students, dberr := storage.GetEnrolledStudents(r.Context(), school_id)

		if dberr != nil {
			slog.Error("database error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}

func GetRejectedStudents(storage Storage.SchoolAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		school_id := r.URL.Query().Get("school_id")

		if school_id == "" {
			slog.Error("params rror", "error:", "School Id missing")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL Request")
			return
		}

		students, dberr := storage.GetRejectedStudents(r.Context(), school_id)

		if dberr != nil {
			slog.Error("database error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}
