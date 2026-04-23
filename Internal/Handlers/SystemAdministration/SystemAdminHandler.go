package SysAdminHandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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
		sessiontoken, _, sysAdminAuth, loginerr := storage.SysAdminLogin(r.Context(), SysAdmin)
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

		message := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Access Your GEOS Greatness</title></head>
<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Peace and Blessings! Your link to access greatness on GEOS is here.</span>
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
  <tr><td align="center">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
      <tr><td style="height:4px;background:linear-gradient(to right,#10b981,#059669);font-size:0;">&nbsp;</td></tr>
      <tr><td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr>
          <td><table cellpadding="0" cellspacing="0" border="0"><tr>
            <td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span></td>
            <td style="padding-left:7px;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span></td>
          </tr></table></td>
          <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f0fdf8;color:#059669;border:1px solid #a7f3d0;text-transform:uppercase;letter-spacing:0.5px;">Link Ready ✓</span></td>
        </tr></table>
      </td></tr>
      <tr><td style="background:#ffffff;padding:28px 32px;text-align:center;">
        <p style="margin:0 0 8px;font-size:32px;line-height:1;">✨</p>
        <p style="margin:8px 0 16px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">Peace and Blessings</p>
        <h1 style="margin:0 0 24px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">Here is your link to access greatness</h1>
        
        <!-- Main CTA Button - Full Width Focus -->
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:32px;"><tr><td align="center">
          <a href="%s" style="display:inline-block;background:linear-gradient(135deg,#0099e6,#5b3fcc);color:#ffffff;font-family:Arial,sans-serif;font-weight:700;font-size:16px;text-decoration:none;padding:20px 50px;border-radius:12px;box-shadow:0 6px 24px rgba(0,153,230,0.4);letter-spacing:0.5px;">Access Greatness →</a>
        </td></tr></table>

        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr><td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td></tr></table>
        <p style="margin:24px 0 0;text-align:center;font-size:13px;color:#4b5a7a;">Questions? <a href="mailto:support@geos.com" style="color:#5b3fcc;font-weight:600;text-decoration:none;">support@geos.com</a></p>
      </td></tr>
      <tr><td style="background:#f8f9fc;padding:16px 32px;border-top:1px solid #f0f2f8;text-align:center;">
        <p style="margin:0 0 7px;"><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Privacy Policy</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Help Centre</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Unsubscribe</a></p>
        <p style="margin:0;font-size:10.5px;color:#c0c8d8;">© 2026 GEOS Ltd · London, UK</p>
      </td></tr>
    </table>
  </td></tr>
</table>
</body>
</html>`, link)
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
			// message = fmt.Sprintf("Peace and Blessings be upon you. Congratulations Your application id Number: %s has been accepted. Here is your email & password to access account. Email: %s /n Password: %s. below is attached your school code to share with your students on sign up %s", schoolInfo.SchoolId, schoolInfo.Sys_Eamil, generatePassword, schoolInfo.Code)
			message = fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Your School is Verified – GEOS</title></head>
		<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
		<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Your school is verified on GEOS! Your portal is live and your student referral code is inside.</span>
		<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
  <tr>
  	<td align="center">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
      <tr>
	  	<td style="height:4px;background:linear-gradient(to right,#10b981,#059669);font-size:0;">&nbsp;
	  	</td>
	  </tr>
      
	  <tr>
	  <td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr>   
		<td>
			<table cellpadding="0" cellspacing="0" border="0">
				<tr>
					<td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span></td>
					<td style="padding-left:7px;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span></td>
				</tr>
			</table>
		  </td>
          
		  <td style="text-align:right;">
		  	<span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f0fdf8;color:#059669;border:1px solid #a7f3d0;text-transform:uppercase;letter-spacing:0.5px;">School Verified ✓
			</span>
			</td>
        </tr></table>
      </td></tr>
      <tr>
	  	<td style="background:#ffffff;padding:28px 32px;">
        <p style="margin:0 0 4px;font-size:32px;line-height:1;">🎉</p>
        <p style="margin:8px 0 6px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">Welcome to GEOS, %s!</p>
        <h1 style="margin:0 0 14px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">Your school is verified and live.</h1>
        <p style="margin:0 0 12px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">Congratulations — <strong style="color:#131c36;">%s</strong> has been successfully verified on GEOS. Your School Portal is now active and ready to use.</p>
        <p style="margin:0 0 18px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">Below is your school's unique referral code. Share this with your students when they register — it unlocks <strong style="color:#131c36;">premium features at no extra cost</strong> and links them directly to your institution.</p>
       
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;background:#f5f6fa;border:2px dashed #d0d8f0;border-radius:10px;">
  <tr>
    <td style="padding:20px;">

      <!-- Referral Code -->
      <p style="margin:0 0 7px;font-size:10.5px;color:#94a3b8;text-transform:uppercase;letter-spacing:1px;font-weight:600;text-align:center;">
        Your School Referral Code
      </p>

      <p style="margin:0;text-align:center;font-family:'Courier New',monospace;font-size:30px;font-weight:900;color:#131c36;letter-spacing:6px;">
        %s
      </p>

      <p style="margin:6px 0 16px;font-size:11.5px;color:#94a3b8;text-align:center;">
        Share this with your students when they sign up on GEOS
      </p>

      <!-- Divider -->
      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:16px;">
        <tr>
          <td style="height:1px;background:#dbe3f0;font-size:0;line-height:0;">&nbsp;</td>
        </tr>
      </table>

      <!-- Login Details -->
      <p style="margin:0 0 10px;font-size:10.5px;color:#94a3b8;text-transform:uppercase;letter-spacing:1px;font-weight:600;">
        Portal Login Details
      </p>

      <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background:#ffffff;border:1px solid #e6ebf5;border-radius:8px;">
        <tr>
          <td style="padding:12px 14px;border-bottom:1px solid #eef2f8;">
            <p style="margin:0 0 4px;font-size:11px;color:#94a3b8;font-weight:600;text-transform:uppercase;letter-spacing:0.6px;">
              Email
            </p>
            <p style="margin:0;font-size:13.5px;color:#131c36;font-weight:700;word-break:break-all;">
              %s
            </p>
          </td>
        </tr>
        <tr>
          <td style="padding:12px 14px;">
            <p style="margin:0 0 4px;font-size:11px;color:#94a3b8;font-weight:600;text-transform:uppercase;letter-spacing:0.6px;">
              Temporary Password
            </p>
            <p style="margin:0;font-family:'Courier New',monospace;font-size:15px;color:#131c36;font-weight:700;letter-spacing:1px;">
              %s
            </p>
          </td>
        </tr>
      </table>

      <p style="margin:10px 0 0;font-size:11.5px;color:#94a3b8;line-height:1.6;">
        For security, we recommend changing your password after your first login.
      </p>

    </td>
  </tr>
</table>

		<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;">
			<tr>
				<td style="background:#f0fdf8;border-left:3px solid #10b981;border-radius:0 9px 9px 0;padding:13px 15px;">
				<p style="margin:0;font-size:13px;color:#4b5a7a;line-height:1.7;"><strong style="color:#059669;">What the code unlocks for your students</strong><br>Students who register using your code automatically receive <strong style="color:#131c36;">premium access</strong> — unlimited AI university matches, priority scholarship alerts, and dedicated GEOS advisor sessions.</p>
				</td>
			</tr>
		</table>

		<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:22px;"><tr><td align="center">
          <a href="https://geosedutest.web.app/school/login" style="display:inline-block;background:#131c36;color:#ffffff;font-family:Arial,sans-serif;font-weight:700;font-size:14px;text-decoration:none;padding:13px 40px;border-radius:9px;">Go to School Portal →</a>
        </td></tr></table>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;"><tr><td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td></tr></table>
        <!-- Portal features -->
        <p style="margin:0 0 12px;font-size:11px;font-weight:700;color:#94a3b8;letter-spacing:1.5px;text-transform:uppercase;">Your School Portal Includes</p>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;">
          <tr><td style="padding:9px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f0fdf8;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">👥</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Student Management</p><p style="margin:0;font-size:12px;color:#94a3b8;">View and manage all students linked to your school</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:9px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f0f7ff;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">📊</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Application Analytics</p><p style="margin:0;font-size:12px;color:#94a3b8;">Track your students' application progress in real time</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:9px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#fffbf0;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">🏆</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Scholarship Broadcasts</p><p style="margin:0;font-size:12px;color:#94a3b8;">Push relevant scholarship alerts to your entire student cohort</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:9px 0;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f8f5ff;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">📅</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Event Management</p><p style="margin:0;font-size:12px;color:#94a3b8;">Register your school and students for GEOS events and university fairs</p></td>
          </tr></table></td></tr>
				</table>
				<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:14px;"><tr><td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td></tr></table>
				<p style="margin:0;font-size:13px;color:#4b5a7a;line-height:1.8;">Need help getting started? Our partnerships team is here — <a href="mailto:schools@geos.com" style="color:#5b3fcc;font-weight:600;text-decoration:none;">schools@geos.com</a></p>
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
		`,

				schoolInfo.School,
				schoolInfo.School,
				schoolInfo.Code,
				schoolInfo.Sys_Eamil,
				generatePassword,
			)
		} else {
			message = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Application Response</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #222;">
	<h2>Peace and Blessings be upon you</h2>
	<p>Your application to join the system was rejected.</p>
</body>
</html>`
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
		email, password, fname, err := storage.RespondApplication(r.Context(), action, id)

		if err != nil {
			slog.Error("there was an error in db operation", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		message := ""
		if action == "approved" {
			message = fmt.Sprintf(
				`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Welcome to GEOS</title></head>
<body style="margin:0;padding:0;background-color:#f0f2f8;font-family:Arial,sans-serif;">
<span style="display:none;max-height:0;overflow:hidden;mso-hide:all;">Welcome to GEOS! Your account is live — discover universities, find scholarships and track applications.</span>
<table width="100%%" cellpadding="0" cellspacing="0" border="0" style="background-color:#f0f2f8;padding:40px 16px;">
  <tr><td align="center">
    <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="max-width:540px;border-radius:10px;overflow:hidden;box-shadow:0 2px 16px rgba(19,28,54,0.08);">
      <tr><td style="height:4px;background:linear-gradient(to right,#10b981,#059669);font-size:0;">&nbsp;</td></tr>
      <tr><td style="background:#ffffff;padding:20px 32px 18px;border-bottom:1px solid #f0f2f8;">
        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr>
          <td><table cellpadding="0" cellspacing="0" border="0"><tr>
            <td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:7px;width:26px;height:26px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:12px;color:#fff;line-height:26px;display:block;">G</span></td>
            <td style="padding-left:7px;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:16px;color:#131c36;letter-spacing:1.5px;">GEOS</span></td>
          </tr></table></td>
          <td style="text-align:right;"><span style="font-size:10px;font-weight:700;padding:4px 10px;border-radius:20px;background:#f0fdf8;color:#059669;border:1px solid #a7f3d0;text-transform:uppercase;letter-spacing:0.5px;">Account Active ✓</span></td>
        </tr></table>
      </td></tr>
      <tr><td style="background:#ffffff;padding:28px 32px;">
        <p style="margin:0 0 4px;font-size:32px;line-height:1;">🎉</p>
        <p style="margin:8px 0 6px;font-size:11.5px;font-weight:600;color:#94a3b8;letter-spacing:1px;text-transform:uppercase;">Welcome, Student!</p>
        <h1 style="margin:0 0 14px;font-family:Georgia,serif;font-size:22px;font-weight:900;color:#131c36;line-height:1.25;">You're officially on GEOS.</h1>
        <p style="margin:0 0 18px;font-size:13.5px;color:#4b5a7a;line-height:1.8;">Your account is live. GEOS is your all-in-one platform to discover universities, track applications, and find scholarships around the world — and we're thrilled to have you.</p>
        <!-- Profile strip -->
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:20px;background:#f8f9fc;border:1px solid #eef1f8;border-radius:9px;"><tr><td style="padding:13px 14px;">
          <table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:42px;vertical-align:middle;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:linear-gradient(135deg,#0099e6,#5b3fcc);border-radius:9px;width:40px;height:40px;text-align:center;vertical-align:middle;"><span style="font-family:Georgia,serif;font-weight:900;font-size:14px;color:#fff;line-height:40px;display:block;">%s</span></td></tr></table></td>
            <td style="padding-left:12px;vertical-align:middle;"><p style="margin:0;font-size:13.5px;font-weight:700;color:#131c36;">%s.</p><p style="margin:2px 0 0;font-size:11.5px;color:#94a3b8;">Email: %s</p><p style="margin:2px 0 0;font-size:11.5px;color:#94a3b8;">Password: %s</p></td>
            <td style="text-align:right;vertical-align:middle;"><span style="font-size:10.5px;font-weight:700;color:#059669;background:#f0fdf8;border:1px solid #a7f3d0;padding:3px 9px;border-radius:6px;">✓ Verified</span></td>
          </tr></table>
        </td></tr></table>
        <!-- Features -->
        <p style="margin:0 0 12px;font-size:11px;font-weight:700;color:#94a3b8;letter-spacing:1.5px;text-transform:uppercase;">Here's what's waiting for you</p>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;">
          <tr><td style="padding:10px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f0f7ff;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">✦</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">AI University Matching</p><p style="margin:0;font-size:12px;color:#94a3b8;">Get matched to universities that fit your grades and goals — up to 94%% accuracy.</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:10px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f0fdf8;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">🏆</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Scholarship Finder</p><p style="margin:0;font-size:12px;color:#94a3b8;">Discover over £2 billion in tracked grants, bursaries, and scholarships.</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:10px 0;border-bottom:1px solid #f0f2f8;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#fffbf0;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">📋</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">Application Tracker</p><p style="margin:0;font-size:12px;color:#94a3b8;">Manage all your applications, deadlines, and documents in one place.</p></td>
          </tr></table></td></tr>
          <tr><td style="padding:10px 0;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0" width="100%%"><tr>
            <td style="width:38px;vertical-align:top;"><table cellpadding="0" cellspacing="0" border="0"><tr><td style="background:#f8f5ff;border-radius:8px;width:34px;height:34px;text-align:center;vertical-align:middle;font-size:16px;line-height:34px;">🤖</td></tr></table></td>
            <td style="padding-left:12px;vertical-align:top;"><p style="margin:0 0 2px;font-size:13px;font-weight:600;color:#131c36;">UniGPT 3.0 — AI Advisor</p><p style="margin:0;font-size:12px;color:#94a3b8;">Get essay feedback, interview prep, and instant university answers.</p></td>
          </tr></table></td></tr>
        </table>
        <!-- Stats -->
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:22px;border:1px solid #eef1f8;border-radius:9px;overflow:hidden;">
          <tr>
            <td style="text-align:center;padding:13px 8px;border-right:1px solid #eef1f8;"><p style="margin:0;font-family:Georgia,serif;font-size:20px;font-weight:900;color:#131c36;">15k+</p><p style="margin:2px 0 0;font-size:10.5px;color:#94a3b8;line-height:1.3;">Programs worldwide</p></td>
            <td style="text-align:center;padding:13px 8px;border-right:1px solid #eef1f8;"><p style="margin:0;font-family:Georgia,serif;font-size:20px;font-weight:900;color:#131c36;">£2B+</p><p style="margin:2px 0 0;font-size:10.5px;color:#94a3b8;line-height:1.3;">In scholarships</p></td>
            <td style="text-align:center;padding:13px 8px;"><p style="margin:0;font-family:Georgia,serif;font-size:20px;font-weight:900;color:#131c36;">80+</p><p style="margin:2px 0 0;font-size:10.5px;color:#94a3b8;line-height:1.3;">Countries covered</p></td>
          </tr>
        </table>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0" style="margin-bottom:18px;"><tr><td align="center">
          <a href="https://geosedutest.web.app/student/login" style="display:inline-block;background:#131c36;color:#ffffff;font-family:Arial,sans-serif;font-weight:700;font-size:14px;text-decoration:none;padding:13px 40px;border-radius:9px;">Go to My Dashboard →</a>
        </td></tr></table>
        <table width="100%%" cellpadding="0" cellspacing="0" border="0"><tr><td style="height:1px;background:#f0f2f8;font-size:0;">&nbsp;</td></tr></table>
        <p style="margin:16px 0 0;text-align:center;font-size:13px;color:#4b5a7a;">Questions? We're here — <a href="mailto:support@geos.com" style="color:#5b3fcc;font-weight:600;text-decoration:none;">support@geos.com</a></p>
      </td></tr>
      <tr><td style="background:#f8f9fc;padding:16px 32px;border-top:1px solid #f0f2f8;text-align:center;">
        <p style="margin:0 0 7px;"><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Privacy Policy</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Help Centre</a><a href="#" style="font-size:11px;color:#94a3b8;text-decoration:underline;margin:0 8px;">Unsubscribe</a></p>
        <p style="margin:0;font-size:10.5px;color:#c0c8d8;">© 2026 GEOS Ltd · London, UK</p>
      </td></tr>
    </table>
  </td></tr>
</table>
</body>
</html>
`, strings.ToUpper(string(fname[0])), fname, email, password)
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

func GetScholarships(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		scholarships, dberr := storage.GetScholarships(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, scholarships)
	}
}
func AddScholarships(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var scholarship Types.Scholarship

		decerr := json.NewDecoder(r.Body).Decode(&scholarship)

		if decerr != nil {
			slog.Error("Decoding error", "error", decerr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(decerr))
			return
		}

		dberr := storage.AddScholarship(r.Context(), scholarship)
		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "scholarship added succesfully!")
	}
}

func UpdateScholarship(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var scholarship Types.Scholarship
		scholarshipId := r.URL.Query().Get("scholarship_id")

		if scholarshipId == "" {
			slog.Error("missing parameter", "error", "scholarship id is missing")
			response.WriteJson(w, http.StatusInternalServerError, "invalid URL requested")
			return
		}

		decerr := json.NewDecoder(r.Body).Decode(&scholarship)

		if decerr != nil {
			slog.Error("Decoding error", "error", decerr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(decerr))
			return
		}

		dberr := storage.UpdateScholarship(r.Context(), scholarship, scholarshipId)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "scholarship updated succesfully!")
	}
}

func DeleteScholarship(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		scholarshipId := r.URL.Query().Get("scholarship_id")

		if scholarshipId == "" {
			slog.Error("missing parameter", "error", "scholarship id is missing")
			response.WriteJson(w, http.StatusInternalServerError, "invalid URL requested")
			return
		}

		dberr := storage.DeleteScholarship(r.Context(), scholarshipId)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "scholarship deleted succesfully!")
	}
}

func CreateWebinar(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var webinar Types.Webinar

		decoderrr := json.NewDecoder(r.Body).Decode(&webinar)

		if decoderrr != nil {
			slog.Error("body error", "error", decoderrr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(decoderrr))
			return

		}
		webinar_code, dberr := storage.CreateWebinar(r.Context(), webinar)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, webinar_code)

	}
}
func GetWebinars(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		webinars, dberr := storage.GetWebinars(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, webinars)

	}
}

func UpdateWebinar(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webinar_id := r.URL.Query().Get("webinar_id")
		var webinar Types.Webinar

		decoderrr := json.NewDecoder(r.Body).Decode(&webinar)

		if decoderrr != nil {
			slog.Error("body error", "error", decoderrr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(decoderrr))
			return

		}
		if webinar_id == "" {
			slog.Error("query error", "error", "Missing Webinar Id")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL request")
			return
		}

		dberr := storage.UpdateWebinar(r.Context(), webinar_id, webinar)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "Webinar update succesfully")

	}
}

func DeleteWebinar(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webinar_id := r.URL.Query().Get("webinar_id")

		if webinar_id == "" {
			slog.Error("query error", "error", "Missing Webinar Id")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL request")
			return
		}

		dberr := storage.DeleteWebinar(r.Context(), webinar_id)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "Webinar Deleted succesfully")
	}
}

func GetUniversities(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		listOfUniversities, dberr := storage.GetUniversities(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, listOfUniversities)
	}
}

func AddFeaturedPartners(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var partners []Types.FeaturedPartner
		err := json.NewDecoder(r.Body).Decode(&partners)
		if err != nil {
			slog.Error("body error", "error", err)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		dberr := storage.AddFeaturedPartners(r.Context(), partners)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "Partners Added Succesfully")
	}
}

func GetFeaturedPartners(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		listOfUniversities, dberr := storage.GetFeaturedPartners(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, listOfUniversities)
	}
}

func DeleteFeaturedPartner(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		partner_id := r.URL.Query().Get("partner_id")

		if partner_id == "" {
			slog.Error("query error", "error", "Missing Partner Id")
			response.WriteJson(w, http.StatusBadRequest, "Invalid URL request")
			return
		}

		dberr := storage.DeleteFeaturedPartner(r.Context(), partner_id)

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, "Partner. Succesfully Deleted")
	}
}

func GetUniversitiesCommissions(storage Storage.SysAdmin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commisions, dberr := storage.GetUniversitiesCommissions(r.Context())

		if dberr != nil {
			slog.Error("Db error", "error", dberr)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(dberr))
			return
		}

		response.WriteJson(w, http.StatusOK, commisions)
	}
}
