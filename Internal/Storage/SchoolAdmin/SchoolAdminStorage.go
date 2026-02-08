package SchoolAdminStorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
)

type SchoolAdminStore struct {
	*Postgress.Postgress
}

func NewSchoolAdminStore(pg *Postgress.Postgress) *SchoolAdminStore {
	return &SchoolAdminStore{pg}
}

func (p *SchoolAdminStore) SchoolAdminLogin(ctx context.Context, schooladmin Types.SchoolAdminLogin) (sessionToken string, csrfToken string, sysadminauth *Types.SchoolAdminAuthenticated, err error) {
	var hashedPassword string

	SchoolAdminAut := Types.SchoolAdminAuthenticated{
		Authenticated: true,
		Status:        true,
	}

	sessionId, sessionerr := Tokens.GenerateToken(10)
	if sessionerr != nil {
		slog.Info("There was a session ID token generation error", "error", sessionerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, sessionerr
	}
	queryerr := p.DB.QueryRowContext(ctx, "SELECT school_id,sys_email,hashed_password,role,name from school_credentials WHERE sys_email = $1", schooladmin.Email).Scan(&SchoolAdminAut.Id, &SchoolAdminAut.Email, &hashedPassword, &SchoolAdminAut.Role, &SchoolAdminAut.Name)
	if queryerr != nil {
		slog.Info("There was an error in querying hashed password from db", "error", queryerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, queryerr
	}

	passwordmatch, matcherr := HashPassword.Unhashpassword(schooladmin.Password, hashedPassword)

	if matcherr != nil {
		slog.Info("There was an internal error", "error", "Hashin algorithim error")
		return "", "", &Types.SchoolAdminAuthenticated{}, errors.New("authentication Error")
	}
	if !passwordmatch {
		slog.Info("There was an auth error", "error", "Password/Email is wrong")
		return "", "", &Types.SchoolAdminAuthenticated{}, errors.New("authentication Error")
	}
	session_token, stokenerr := Tokens.GenerateToken(10)
	if stokenerr != nil {
		slog.Info("There was a session token generation error", "error", stokenerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, stokenerr
	}
	csrf_token, csrftokenerr := Tokens.GenerateToken(10)
	if csrftokenerr != nil {
		slog.Info("There was a csrf token generation error", "error", stokenerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, csrftokenerr
	}

	SchoolAdminAut.CsrfToken = csrf_token

	_, insertqerr := p.DB.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, session_id, credential_id,role)  VALUES ($1, $2, $3, $4, $5,$6)", session_token, csrf_token, SchoolAdminAut.Email, sessionId, SchoolAdminAut.Id, SchoolAdminAut.Role)
	if insertqerr != nil {
		slog.Info("There was an error inserting data to db", "error", insertqerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, nil
	}

	return session_token, csrf_token, &SchoolAdminAut, nil
}

func (p *SchoolAdminStore) ValidateLink(ctx context.Context, invitation_id string) (string, error) {
	var returneTimeStamp string
	var status string

	dberr := p.DB.QueryRowContext(ctx, "SELECT  created_at, invitation_status from school_invites WHERE invitation_id = $1", invitation_id).Scan(&returneTimeStamp, &status)

	if dberr != nil {
		return "", dberr
	}

	parsedTime, terr := time.Parse(time.RFC3339Nano, returneTimeStamp)

	if terr != nil {
		slog.Info("There was an error parsing time", "error", terr)
		return "", terr
	}

	expired := time.Now().Before(parsedTime)

	if expired {
		return "", errors.New("invitation is expired")
	}

	return status, nil
}

func (p *SchoolAdminStore) SubmitInvite(ctx context.Context, schoolInfo Types.SchoolInformation) (string, error) {
	var status string
	currentTime := time.Now()
	dbtx, txerr := p.DB.BeginTx(ctx, nil)

	if txerr != nil {
		return "", txerr
	}

	qerr := dbtx.QueryRowContext(ctx, "INSERT into school_applications (status,school_name,school_email, admin_name , school_phone , school_country , registration_code , school_curriculum , school_branch , school_city , completed_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING status", "completed",
		schoolInfo.School, schoolInfo.Email, schoolInfo.Admin, schoolInfo.Phone, schoolInfo.Country, schoolInfo.RegistrationCode, schoolInfo.Curriculam, schoolInfo.Branch, schoolInfo.City, currentTime).Scan(&status)

	if qerr != nil {
		slog.Info("Query Error", "message", "there was an error saving school info to db", "error", qerr)
		dbtx.Rollback()
		return "", qerr
	}

	_, qerr2 := dbtx.ExecContext(ctx, `UPDATE school_invites set invitation_status = $1`, "completed")

	if qerr2 != nil {
		slog.Info("Query Error", "message", "there was an error saving school info to db", "error", qerr)
		dbtx.Rollback()
		return "", qerr2
	}

	dbtx.Commit()
	return "completed", nil
}

func (p *SchoolAdminStore) GetUnProcessedStudentsList(ctx context.Context, school_id string) ([]Types.StudentProfile, error) {
	var school_code string

	var studentDetailsList []Types.StudentProfile

	err := p.DB.QueryRowContext(ctx, "SELECT school_code from school_codes where school_id = $1", school_id).Scan(&school_code)
	if err != nil {
		return nil, err
	}

	rows, q2err := p.DB.QueryContext(ctx, `SELECT DISTINCT 
    stc.student_id,
    sp.first_name, sp.last_name, sp.middle_name, sp.dob, sp.gender, sp.nationality, 
    sp.passport_number, sp.passport_expiry, sp.marrital_status,
    sc.email, sc.phone_number, sc.whatsapp_number, sc.permanent_address, sc.street_address,
    sc.city, sc.state_province, sc.zip_postal_code, sc.emergency_phone, 
    sc.emergency_contact_name, sc.emergency_relationship,
    se.school_name, se.curriculum, se.graduation_year, se.cummulative_score,
    se.language_type, se.language_overall_score, se.language_reading, 
    se.language_writting, se.language_speaking, se.language_listening,
    spr.primary_career_interest, spr.degree_level, spr.preferred_start_date, 
    spr.annual_budget, spr.scholarship_interest
	FROM students_credentials AS stc 
	LEFT JOIN student_profile AS sp ON stc.student_id = sp.student_id 
	LEFT JOIN student_contact AS sc ON stc.student_id = sc.student_id
	LEFT JOIN student_education AS se ON stc.student_id = se.student_id
	LEFT JOIN students_preferences AS spr ON stc.student_id = spr.student_id
	WHERE stc.school_id = $1 AND stc.school_verified = $2
`, school_code, "un-verified")

	if q2err != nil {
		return nil, q2err
	}
	defer rows.Close()

	for rows.Next() {
		var studentDetails Types.StudentProfile
		scanerr := rows.Scan(
			&studentDetails.StudentId,
			&studentDetails.StudentPersonalDetails.First_Name,
			&studentDetails.StudentPersonalDetails.Last_Name,
			&studentDetails.StudentPersonalDetails.Middle_Name,
			&studentDetails.StudentPersonalDetails.Dob,
			&studentDetails.StudentPersonalDetails.Gender,
			&studentDetails.StudentCitizenshipDetails.Nationality,
			&studentDetails.StudentCitizenshipDetails.Passport_Number,
			&studentDetails.StudentCitizenshipDetails.Passport_Expiry,
			&studentDetails.StudentPersonalDetails.Marrital_Status,
			&studentDetails.StudentContact.Email,
			&studentDetails.StudentContact.Phone,
			&studentDetails.StudentContact.WhatsApp,
			&studentDetails.StudentContact.Address,
			&studentDetails.StudentContact.Street,
			&studentDetails.StudentContact.City,
			&studentDetails.StudentContact.Province,
			&studentDetails.StudentContact.PostCode,
			&studentDetails.StudentContact.EmergencyPhone,
			&studentDetails.StudentContact.EmergencyName,
			&studentDetails.StudentContact.EmergencyRelation,
			&studentDetails.StudentEducation.SchoolName,
			&studentDetails.StudentEducation.Curriculum,
			&studentDetails.StudentEducation.GraduationYear,
			&studentDetails.StudentEducation.CummulativeScore,
			&studentDetails.StudentEducation.LanguageType,
			&studentDetails.StudentEducation.LanguageOverallScore,
			&studentDetails.StudentEducation.LanguageReading,
			&studentDetails.StudentEducation.LanguageWritting,
			&studentDetails.StudentEducation.LanguageSpeaking,
			&studentDetails.StudentEducation.LanguageListening,
			&studentDetails.StudentPrefferences.PrimaryInterest,
			&studentDetails.StudentPrefferences.Degree,
			&studentDetails.StudentPrefferences.PrefferedDate,
			&studentDetails.StudentPrefferences.AnnualBudget,
			&studentDetails.StudentPrefferences.Schalarship,
		)

		if scanerr != nil {
			return nil, scanerr
		}

		studentDetailsList = append(studentDetailsList, studentDetails)
	}
	return studentDetailsList, nil
}

func (p *SchoolAdminStore) VerifyStudentAccount(ctx context.Context, school_id string, student_id string, status string) error {
	fmt.Println("running db")

	_, qerr := p.DB.ExecContext(ctx, "UPDATE students_credentials sc set school_verified = $1 , school_id = $2  from school_codes scs where sc.school_id = scs.school_code AND scs.school_id = $3 AND sc.student_id = $4 ", status, school_id, school_id, student_id)
	if qerr != nil {
		return qerr
	}
	return nil
}

func (p *SchoolAdminStore) GetProcessedStudentsList(ctx context.Context, school_id string, status string) ([]Types.StudentProfile, error) {

	var studentDetailsList []Types.StudentProfile
	rows, q2err := p.DB.QueryContext(ctx, `SELECT  
	stc.student_id,
	sp.first_name, sp.last_name, sp.middle_name, sp.dob, sp.gender, sp.nationality, 
	sp.passport_number, sp.passport_expiry, sp.marrital_status,
	sc.email, sc.phone_number, sc.whatsapp_number, sc.permanent_address, sc.street_address,
	sc.city, sc.state_province, sc.zip_postal_code, sc.emergency_phone, 
	sc.emergency_contact_name, sc.emergency_relationship,
	se.school_name, se.curriculum, se.graduation_year, se.cummulative_score,
	se.language_type, se.language_overall_score, se.language_reading, 
	se.language_writting, se.language_speaking, se.language_listening,
	spr.primary_career_interest, spr.degree_level, spr.preferred_start_date, 
	spr.annual_budget, spr.scholarship_interest,

	json_agg(
        json_build_object(
            'document_id', doc.document_id,
            'name', doc.document_name,
            'document_name', doc.document_file_name,
            'status', doc.document_status,
            'data', doc.document,
            'type', doc.document_type
        	) 
		)AS documents

	FROM students_credentials AS stc 
	LEFT JOIN student_profile AS sp ON stc.student_id = sp.student_id 
	LEFT JOIN student_contact AS sc ON stc.student_id = sc.student_id
	LEFT JOIN student_education AS se ON stc.student_id = se.student_id
	LEFT JOIN students_preferences AS spr ON stc.student_id = spr.student_id
	LEFT JOIN students_documents AS doc ON doc.student_id = spr.student_id
	WHERE stc.school_id = $1 AND stc.school_verified = $2
	GROUP BY
    stc.student_id,
    sp.first_name, sp.last_name, sp.middle_name, sp.dob, sp.gender, sp.nationality, 
    sp.passport_number, sp.passport_expiry, sp.marrital_status,
    sc.email, sc.phone_number, sc.whatsapp_number, sc.permanent_address, sc.street_address,
    sc.city, sc.state_province, sc.zip_postal_code, sc.emergency_phone, 
    sc.emergency_contact_name, sc.emergency_relationship,
    se.school_name, se.curriculum, se.graduation_year, se.cummulative_score,
    se.language_type, se.language_overall_score, se.language_reading, 
    se.language_writting, se.language_speaking, se.language_listening,
    spr.primary_career_interest, spr.degree_level, spr.preferred_start_date, 
    spr.annual_budget, spr.scholarship_interest
`, school_id, status)

	if q2err != nil {
		return nil, q2err
	}
	defer rows.Close()

	for rows.Next() {
		var studentDetails Types.StudentProfile
		var docBytes []byte

		scanerr := rows.Scan(
			&studentDetails.StudentId,
			&studentDetails.StudentPersonalDetails.First_Name,
			&studentDetails.StudentPersonalDetails.Last_Name,
			&studentDetails.StudentPersonalDetails.Middle_Name,
			&studentDetails.StudentPersonalDetails.Dob,
			&studentDetails.StudentPersonalDetails.Gender,
			&studentDetails.StudentCitizenshipDetails.Nationality,
			&studentDetails.StudentCitizenshipDetails.Passport_Number,
			&studentDetails.StudentCitizenshipDetails.Passport_Expiry,
			&studentDetails.StudentPersonalDetails.Marrital_Status,
			&studentDetails.StudentContact.Email,
			&studentDetails.StudentContact.Phone,
			&studentDetails.StudentContact.WhatsApp,
			&studentDetails.StudentContact.Address,
			&studentDetails.StudentContact.Street,
			&studentDetails.StudentContact.City,
			&studentDetails.StudentContact.Province,
			&studentDetails.StudentContact.PostCode,
			&studentDetails.StudentContact.EmergencyPhone,
			&studentDetails.StudentContact.EmergencyName,
			&studentDetails.StudentContact.EmergencyRelation,
			&studentDetails.StudentEducation.SchoolName,
			&studentDetails.StudentEducation.Curriculum,
			&studentDetails.StudentEducation.GraduationYear,
			&studentDetails.StudentEducation.CummulativeScore,
			&studentDetails.StudentEducation.LanguageType,
			&studentDetails.StudentEducation.LanguageOverallScore,
			&studentDetails.StudentEducation.LanguageReading,
			&studentDetails.StudentEducation.LanguageWritting,
			&studentDetails.StudentEducation.LanguageSpeaking,
			&studentDetails.StudentEducation.LanguageListening,
			&studentDetails.StudentPrefferences.PrimaryInterest,
			&studentDetails.StudentPrefferences.Degree,
			&studentDetails.StudentPrefferences.PrefferedDate,
			&studentDetails.StudentPrefferences.AnnualBudget,
			&studentDetails.StudentPrefferences.Schalarship,
			&docBytes,
		)

		if scanerr != nil {
			return nil, scanerr
		}

		err := json.Unmarshal(docBytes, &studentDetails.Documents)

		if err != nil {
			return nil, err

		}

		studentDetailsList = append(studentDetailsList, studentDetails)
	}
	return studentDetailsList, nil
}

func (p SchoolAdminStore) GetSchoolProfileData(ctx context.Context, school_id string) (Types.SchoolInformation, error) {

	var schoolData Types.SchoolInformation

	dberr := p.DB.QueryRowContext(ctx, `SELECT 
	registration_code,
	school_name,
	school_email,
	admin_name,
	school_phone,
	school_country,
	school_curriculum,
	school_branch,
	school_city,
	school_id,
	status,
	sys_email,
	username
	 from schools where school_id = $1`, school_id).Scan(
		&schoolData.RegistrationCode,
		&schoolData.School,
		&schoolData.Email,
		&schoolData.Admin,
		&schoolData.Phone,
		&schoolData.Country,
		&schoolData.Curriculam,
		&schoolData.Branch,
		&schoolData.City,
		&schoolData.SchoolId,
		&schoolData.Status,
		&schoolData.Sys_Eamil,
		&schoolData.Username)

	if dberr != nil {
		return Types.SchoolInformation{}, dberr
	}
	return schoolData, nil
}
