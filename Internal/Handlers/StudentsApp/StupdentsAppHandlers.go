package StudentAppHandler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

		sessiontoken, csrftoken, studentAuth, loginerr := storage.StudentsLogin(r.Context(), loginDetails)
		if loginerr != nil {
			slog.Error("error in login ", "error", loginerr)
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
			Domain:   "localhost", // ✅ Cross-port

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
			Domain:   "localhost", // ✅ Cross-port

			// SameSite: http.SameSiteNoneMode, // Lax works without HTTPS
			// Path:     "/schoolAdmin",          // works for all paths
		})

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
		documentName := r.URL.Query().Get("docname")
		studentId := r.URL.Query().Get("student_id")

		if studentId == "" || documentName == "" {
			slog.Error("error in requested url", "error", "missing parameters")
			response.WriteJson(w, http.StatusBadRequest, "Missing Parameters")
			return
		}

		document, err := storage.GetStudentsDocument(r.Context(), studentId, documentName)

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
