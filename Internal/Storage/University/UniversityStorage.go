package UniversityStorage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	Helper "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Helpers"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
)

type UniversityStore struct {
	*Postgress.Postgress
}

func NewUniversityStore(pg *Postgress.Postgress) *UniversityStore {
	return &UniversityStore{pg}
}

func (p *UniversityStore) UniversityLogin(ctx context.Context, universityLogin Types.UniversityLogin) (sessionToken string, csrfToken string, universityAuthe *Types.UniversityAuthenticated, err error) {
	var hashedPassword string

	universityAuth := Types.UniversityAuthenticated{
		Authenticated: true,
		Status:        true,
	}

	sessionId, sessionerr := Tokens.GenerateToken(10)
	if sessionerr != nil {
		slog.Info("There was a session ID token generation error", "error", sessionerr)
		return "", "", nil, sessionerr
	}
	queryerr := p.DB.QueryRowContext(ctx, "SELECT uc.university_id,uc.university_email,uc.hashed_password,uc.role,uc.status, ud.university_name from university_credentials uc LEFT  join universities ud on ud.university_id = uc.university_id WHERE uc.university_email = $1", universityLogin.Email).Scan(&universityAuth.Id, &universityAuth.Email, &hashedPassword, &universityAuth.Role, &universityAuth.Status, &universityAuth.Name)
	if queryerr != nil {
		slog.Info("There was an error in querying hashed password from db", "error", queryerr)
		return "", "", nil, queryerr
	}

	passwordmatch, matcherr := HashPassword.Unhashpassword(universityLogin.Password, hashedPassword)

	if matcherr != nil {
		slog.Info("There was an internal error", "error", "Hashin algorithim error")
		return "", "", nil, errors.New("authentication Error")
	}
	if !passwordmatch {
		slog.Info("There was an auth error", "error", "Password/Email is wrong")
		return "", "", nil, errors.New("authentication Error")
	}
	session_token, stokenerr := Tokens.GenerateToken(10)
	if stokenerr != nil {
		slog.Info("There was a session token generation error", "error", stokenerr)
		return "", "", nil, stokenerr
	}
	csrf_token, csrftokenerr := Tokens.GenerateToken(10)
	if csrftokenerr != nil {
		slog.Info("There was a csrf token generation error", "error", stokenerr)
		return "", "", nil, csrftokenerr
	}
	universityAuth.CsrfToken = csrf_token
	_, insertqerr := p.DB.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, session_id, credential_id,role)  VALUES ($1, $2, $3, $4, $5,$6)", session_token, csrf_token, universityAuth.Email, sessionId, universityAuth.Id, universityAuth.Role)
	if insertqerr != nil {
		slog.Info("There was an error inserting data to db", "error", insertqerr)
		return "", "", nil, nil
	}

	return session_token, csrf_token, &universityAuth, nil
}

func (p UniversityStore) GetStudntsApplications(ctx context.Context, university_id string, status string) ([]Types.StudentProfile, error) {

	var studentDetailsList []Types.StudentProfile

	rows, q2err := p.DB.QueryContext(ctx, `SELECT  
    ua.student_id,ua.university_id,ua.decision_status,TO_CHAR(ua.created_at,'Mon DD, YYYY'),ua.application_id,
	p.program_name,p.session_intake,
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
	FROM university_applications AS ua 
	LEFT JOIN programs AS p ON ua.program_id = p.program_id 
	LEFT JOIN student_profile AS sp ON ua.student_id = sp.student_id 
	LEFT JOIN student_contact AS sc ON ua.student_id = sc.student_id
	LEFT JOIN student_education AS se ON ua.student_id = se.student_id
	LEFT JOIN students_preferences AS spr ON ua.student_id = spr.student_id
	WHERE ua.university_id = $1 AND ua.decision_status = $2
`, university_id, status)

	if q2err != nil {
		return nil, q2err
	}
	defer rows.Close()

	for rows.Next() {
		var studentDetails Types.StudentProfile
		scanerr := rows.Scan(
			&studentDetails.StudentId,
			&studentDetails.UniversityId,
			&studentDetails.Decision,
			&studentDetails.Date,
			&studentDetails.ApplicationId,
			&studentDetails.Program,
			&studentDetails.ProgramSession,
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

func (p UniversityStore) GetAllStudntsApplications(ctx context.Context, university_id string, status string) ([]Types.StudentProfile, error) {

	var studentDetailsList []Types.StudentProfile

	rows, q2err := p.DB.QueryContext(ctx, `SELECT  
    ua.student_id,ua.university_id,ua.decision_status,TO_CHAR(ua.created_at,'Mon DD, YYYY'),ua.application_id,
	p.program_name,p.session_intake,
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
	FROM university_applications AS ua 
	LEFT JOIN programs AS p ON ua.program_id = p.program_id 
	LEFT JOIN student_profile AS sp ON ua.student_id = sp.student_id 
	LEFT JOIN student_contact AS sc ON ua.student_id = sc.student_id
	LEFT JOIN student_education AS se ON ua.student_id = se.student_id
	LEFT JOIN students_preferences AS spr ON ua.student_id = spr.student_id
	WHERE ua.decision_status = $1
`, status)

	if q2err != nil {
		return nil, q2err
	}
	defer rows.Close()

	for rows.Next() {
		var studentDetails Types.StudentProfile
		scanerr := rows.Scan(
			&studentDetails.StudentId,
			&studentDetails.UniversityId,
			&studentDetails.Decision,
			&studentDetails.Date,
			&studentDetails.ApplicationId,
			&studentDetails.Program,
			&studentDetails.ProgramSession,
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

func (p UniversityStore) RespondStudntsApplications(ctx context.Context, application_id string, status string) error {
	_, dberr := p.DB.ExecContext(ctx, `UPDATE university_applications set decision_status =$1 where application_id = $2`, status, application_id)

	if dberr != nil {
		return dberr
	}

	return nil
}

func (p UniversityStore) GetUniversityProgramsList(ctx context.Context, university_id string) ([]Types.ProgrameThumbnail, error) {

	var programsList []Types.ProgrameThumbnail
	rows, dber := p.DB.QueryContext(ctx, `SELECT program_id, program_name, session_intake,program_status, program_capacity from programs where university_id = $1`, university_id)

	if dber != nil {
		return nil, dber
	}

	for rows.Next() {
		var program Types.ProgrameThumbnail
		err := rows.Scan(&program.Id, &program.Name, &program.SessionIntake, &program.Status, &program.Capacity)
		if err != nil {
			return nil, err
		}
		programsList = append(programsList, program)
	}

	return programsList, nil
}
func (p UniversityStore) GetProgramDetails(ctx context.Context, program_id string) (Types.Programe, error) {
	var rawRequirements []byte
	var rawTags []byte
	var rawCareers []byte
	var rawDocs []byte

	var program Types.Programe
	dber := p.DB.QueryRowContext(ctx, `SELECT 
	program_id,
	program_name,
	program_fee,
	program_duration,
	session_intake,
	program_description,
	program_requirements,
	related_tags,
	possible_careers,
	program_required_documents,
	program_application_fee,
	university_id,
	program_level,
	program_capacity
	from programs where program_id = $1`, program_id).Scan(&program.Id, &program.Name, &program.PFee, &program.Duration, &program.SessionIntake,
		&program.Description, &rawRequirements, &rawTags, &rawCareers, &rawDocs, &program.AFee, &program.UniversityCode, &program.ProgramLevel, &program.ProgramCapacity)

	if dber != nil {
		return Types.Programe{}, dber
	}

	_ = json.Unmarshal(rawRequirements, &program.Requirements)
	_ = json.Unmarshal(rawTags, &program.RelatedTags)
	_ = json.Unmarshal(rawDocs, &program.RequiredDocs)
	_ = json.Unmarshal(rawCareers, &program.PossibleCareers)

	return program, nil
}

func (p UniversityStore) AddNewProgram(ctx context.Context, programDetails Types.Programe) error {

	program_id_token, _ := Tokens.GenerateToken(5)
	program_id := fmt.Sprintf("%s%s%s", programDetails.UniversityCode, programDetails.ProgramLevel, program_id_token)

	reqJSON, _ := json.Marshal(programDetails.Requirements)
	tagsJSON, _ := json.Marshal(programDetails.RelatedTags)
	careersJSON, _ := json.Marshal(programDetails.PossibleCareers)
	docsJSON, _ := json.Marshal(programDetails.RequiredDocs)

	_, err := p.DB.ExecContext(ctx,
		`INSERT into 
		programs 
		(program_name,
		program_fee,
		program_duration,
		session_intake,
		program_description,
		program_requirements,
		related_tags,
		possible_careers,
		program_application_fee,
		program_required_documents,
		program_id,
		program_status,
		university_id,
		program_capacity,
		program_level)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
		programDetails.Name, programDetails.PFee, programDetails.Duration, programDetails.SessionIntake, programDetails.Description,
		reqJSON, tagsJSON, careersJSON, programDetails.AFee, docsJSON,
		program_id, "active", programDetails.UniversityCode, programDetails.ProgramCapacity, programDetails.ProgramLevel)

	if err != nil {
		return err
	}
	return nil
}
func (p UniversityStore) UpdateProgram(ctx context.Context, programDetails Types.Programe, program_id string) error {

	reqJSON, _ := json.Marshal(programDetails.Requirements)
	tagsJSON, _ := json.Marshal(programDetails.RelatedTags)
	careersJSON, _ := json.Marshal(programDetails.PossibleCareers)
	docsJSON, _ := json.Marshal(programDetails.RequiredDocs)

	_, err := p.DB.ExecContext(ctx,
		`UPDATE programs
		SET
		program_name = $1,
		program_fee = $2,
		program_duration = $3,
		session_intake = $4,
		program_description = $5,
		program_requirements = $6,
		related_tags = $7,
		possible_careers = $8,
		program_application_fee = $9,
		program_required_documents = $10,
		program_status = $11,
		university_id = $12,
		program_capacity = $13,
		program_level = $14
		WHERE program_id = $15;`,
		programDetails.Name, programDetails.PFee, programDetails.Duration, programDetails.SessionIntake, programDetails.Description,
		reqJSON, tagsJSON, careersJSON, programDetails.AFee, docsJSON, "active", programDetails.UniversityCode, programDetails.ProgramCapacity, programDetails.ProgramLevel, program_id)

	if err != nil {
		return err
	}
	return nil
}

func (p UniversityStore) GetUniversityProfile(ctx context.Context, university_id string) (Types.UniversityProfile, error) {
	var universityprofile Types.UniversityProfile
	var media []byte
	dberr := p.DB.QueryRowContext(ctx, `SELECT
		up.university_id,
		up.university_city,
		up.students_count,
		up.acceptance_rate,
		up.qs_ranking,
		up.about_university,
		up.founded_date,
		up.type,
		up.calendar,
		up.graduation_rate,
		up.employability,
		un.university_name,

		uc.university_admission_email,
		uc.university_email,
		uc.university_address,
		uc.university_phone,
		uc.university_website,
		uc.university_instagram,
		uc.university_youtube,
		uc.university_linkedin,
		uc.university_x,

		json_agg(
        json_build_object(
            'media_id', ul.media_id,
            'media', ul.media,
            'media_type', ul.media_type,
            'media_file_name', ul.media_file_name,
            'media_tag', ul.media_tag
        	) 
		) AS media
		
		from university_profile up
		Left join university_life ul on up.university_id = ul.university_id 
		Left join universities un on up.university_id = un.university_id
		Left join university_contact uc on up.university_id = uc.university_id
		where up.university_id = $1
		GROUP BY 
		up.university_id,
		up.university_city,
		up.students_count,
		up.acceptance_rate,
		up.qs_ranking,
		up.about_university,
		up.founded_date,
		up.type,
		up.calendar,
		up.graduation_rate,
		up.employability,
		uc.university_admission_email,
		uc.university_email,
		uc.university_address,
		uc.university_phone,
		uc.university_website,
		uc.university_instagram,
		uc.university_youtube,
		uc.university_linkedin,
		uc.university_x,
		un.university_name
		`, university_id).Scan(
		&universityprofile.UniversityId,
		&universityprofile.UniversityCity,
		&universityprofile.StudentsCount,
		&universityprofile.AcceptanceRate,
		&universityprofile.QSRanking,
		&universityprofile.AboutUniversity,
		&universityprofile.FoundedDate,
		&universityprofile.Type,
		&universityprofile.Calendar,
		&universityprofile.GraduationRate,
		&universityprofile.Employability,
		&universityprofile.UniversityName,
		&universityprofile.UniContact.UniversityAdmissionEmail,
		&universityprofile.UniContact.UniversityEmail,
		&universityprofile.UniContact.UniversityAddress,
		&universityprofile.UniContact.UniversityPhone,
		&universityprofile.UniContact.UniversityWebsite,
		&universityprofile.UniContact.UniversityInstagram,
		&universityprofile.UniContact.UniversityYoutube,
		&universityprofile.UniContact.UniversityLinkedin,
		&universityprofile.UniContact.UniversityX,
		&media)
	if dberr != nil {
		return Types.UniversityProfile{}, dberr
	}

	err := json.Unmarshal(media, &universityprofile.Media)

	if err != nil {
		return Types.UniversityProfile{}, err

	}

	return universityprofile, nil
}

func (p UniversityStore) UploadCampusMedia(ctx context.Context, university_id string, mediaList []Types.UniMedia) error {

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	INSERT INTO university_life 
	(media_type, media_file_name, media_tag, media,university_id)
		VALUES ($1,$2,$3,$4,$5)
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	defer stmt.Close()

	for _, media := range mediaList {
		medias, _ := Helper.Base64ToBytes(media.Media)
		_, err := stmt.Exec(
			media.MediaType,
			media.MediaFileName,
			media.MediaTag,
			medias,
			university_id,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
