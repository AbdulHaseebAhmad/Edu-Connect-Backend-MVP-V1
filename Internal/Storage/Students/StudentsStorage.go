package StudentsAppStorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	Helper "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Helpers"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
	"github.com/jackc/pgx/pgtype"
)

type StudentsAppStore struct {
	*Postgress.Postgress
}

func NewStudentsAppStore(pg *Postgress.Postgress) *StudentsAppStore {
	return &StudentsAppStore{pg}
}

func (p *StudentsAppStore) StudentSignup(ctx context.Context, studentSignup Types.StudentsSignupData, studentsDocuments Types.StudentSignupDocuments) error {
	slug, terr := Tokens.GenerateToken(10)
	if terr != nil {
		return terr
	}
	_, err := p.DB.ExecContext(
		ctx, `
    INSERT INTO students_applications
    (first_name, last_name, email, citizenship, passport, transcript,slug,passport_mime_type,transcript_mime_type,status,school_code)
    VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11)
    `,
		studentSignup.First_Name,
		studentSignup.Last_Name,
		studentSignup.Email,
		studentSignup.Citizenship,
		studentsDocuments.Passport,
		studentsDocuments.Transcript,
		slug,
		studentSignup.PassportType,
		studentSignup.TranscriptType,
		"pending",
		studentSignup.SchoolCode,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *StudentsAppStore) StudentsLogin(ctx context.Context, studentLogin Types.StudentsSignIn) (sessionTokens string, csrfTokens string, studentAuth *Types.StudentAuthenticated, err error) {
	fmt.Println("Running storage layer")
	var studentAuths Types.StudentAuthenticated
	var hashedPassword string
	var sessionToken string
	var csrfToken string
	var fname string
	var lname string

	tx, txErr := p.DB.BeginTx(ctx, nil)
	if txErr != nil {
		return "", "", studentAuth, txErr
	}

	qerr := tx.QueryRowContext(ctx, `SELECT 
		sc.student_email, 
		sc.role, 
		sc.student_status, 
		sc.student_id, 
		sc.hashed_password, 
		sc.school_verified,
		
		sa.first_name,
		sa.last_name
		FROM students_credentials sc 
		
		LEFT JOIN students_applications sa on sa.email = sc.student_email 
		
		WHERE student_email = $1`, studentLogin.Email).Scan(&studentAuths.Email, &studentAuths.Role, &studentAuths.Status, &studentAuths.Id, &hashedPassword, &studentAuths.SchoolVerified, &fname, &lname)

	if qerr != nil {
		tx.Rollback()
		return "", "", &studentAuths, qerr
	}

	matched, hasherr := HashPassword.Unhashpassword(studentLogin.Password, hashedPassword)

	if hasherr != nil {
		return "", "", &studentAuths, hasherr
	}

	if matched == true {
		fmt.Println("matched")
		session_id, sessionerr := Tokens.GenerateToken(10)
		if sessionerr != nil {
			return "", "", &studentAuths, sessionerr
		}
		session_token, stokenerr := Tokens.GenerateToken(10)
		if stokenerr != nil {
			return "", "", &studentAuths, sessionerr
		}
		sessionToken = session_token
		csrf_token, csrftokenerr := Tokens.GenerateToken(10)
		if csrftokenerr != nil {
			return "", "", &studentAuths, sessionerr
		}
		csrfToken = csrf_token
		studentAuths.CsrfToken = csrfToken
		_, insertqerr := tx.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, session_id,role,credential_id)  VALUES ($1, $2, $3, $4, $5,$6)", session_token, csrf_token, studentAuths.Email, session_id, studentAuths.Role, studentAuths.Id)
		if insertqerr != nil {
			tx.Rollback()
			return "", "", &studentAuths, insertqerr
		}
		studentAuths.Name = fname + " " + lname
	}
	tx.Commit()
	fmt.Println("st=", sessionToken, "ct=", csrfToken)
	return sessionToken, csrfToken, &studentAuths, nil

}

func (p *StudentsAppStore) GetCountryList(ctx context.Context, countryCode string) ([]Types.Countries, error) {
	if countryCode != "" {
	}
	rows, err := p.DB.QueryContext(ctx,
		`SELECT 
			cl.country_code,
			cl.dialing_code,
			cl.id,
			cl.image_url,
			cl.iso_code,
			cl.iso3_code,
			cl.name,
			COUNT(pr.program_id) AS program_count
		FROM countries cl
		LEFT JOIN universities un ON un.university_country = cl.country_code
		LEFT JOIN programs pr ON pr.university_id = un.university_id
		GROUP BY 
			cl.country_code,
			cl.dialing_code, 
			cl.id,
			cl.image_url,
			cl.iso_code,
			cl.iso3_code,
			cl.name
		ORDER BY program_count DESC`)

	if err != nil {
		return []Types.Countries{}, err
	}
	defer rows.Close()
	var countriesList []Types.Countries
	for rows.Next() {
		var country Types.Countries
		serr := rows.Scan(&country.CountryCode, &country.DiallingCode, &country.Id, &country.Image, &country.IsoCode, &country.IsoCode3, &country.Name, &country.Programs)
		if serr != nil {
			return []Types.Countries{}, serr
		}
		countriesList = append(countriesList, country)
	}
	return countriesList, nil
}

func (p *StudentsAppStore) GetUniversitiesList(ctx context.Context, countryCode string) ([]Types.University, error) {

	rows, err := p.DB.QueryContext(ctx, `
	SELECT 
    un.university_id,
    un.university_name,
    un.university_city,
    un.university_country,
    un.university_acronym,
    un.university_phone,
    un.university_image,
    uc.iso_code AS country_code
	FROM universities un 
	LEFT JOIN countries uc ON uc.country_code = un.university_country  
	WHERE un.university_country = $1`, countryCode)
	if err != nil {
		return []Types.University{}, err
	}
	defer rows.Close()
	var universityList []Types.University
	for rows.Next() {
		var university Types.University
		serr := rows.Scan(&university.Id, &university.Name, &university.City, &university.Country, &university.Acronym, &university.Phone, &university.Image, &university.CountryName)
		if serr != nil {
			return []Types.University{}, serr
		}
		universityList = append(universityList, university)
	}
	return universityList, nil
}

func (p StudentsAppStore) GetUniversityProfile(ctx context.Context, university_id string) (Types.UniversityProfile, error) {
	var universityprofile Types.UniversityProfile
	var programs []byte
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

		json_agg(
        json_build_object(
            'program_id', pr.program_id,
            'program_name', pr.program_name,
            'program_level', pr.program_level,
            'program_duration', pr.program_duration,
            'session_intake', pr.session_intake,
            'program_fee', pr.program_fee,
            'program_status', pr.program_status,
			'application_fee',pr.program_application_fee
        	) 
		)AS programs
		
		from university_profile up
		Left join programs pr on up.university_id = pr.university_id 
		Left join universities un on up.university_id = un.university_id
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
		un.university_name
		`, university_id).Scan(&universityprofile.UniversityId, &universityprofile.UniversityCity, &universityprofile.StudentsCount, &universityprofile.AcceptanceRate, &universityprofile.QSRanking, &universityprofile.AboutUniversity, &universityprofile.FoundedDate, &universityprofile.Type, &universityprofile.Calendar, &universityprofile.GraduationRate, &universityprofile.Employability, &universityprofile.UniversityName, &programs)
	if dberr != nil {
		return Types.UniversityProfile{}, dberr
	}

	err := json.Unmarshal(programs, &universityprofile.Programs)

	if err != nil {
		return Types.UniversityProfile{}, dberr

	}

	return universityprofile, nil
}

func (p *StudentsAppStore) GetUniversityPrograms(ctx context.Context, program_id string) (programsList []Types.Programe, err error) {
	rows, err := p.DB.QueryContext(ctx,
		`SELECT p.program_id, p.program_name, p.program_fee, 
            p.program_duration, p.session_intake, p.program_description, 
            p.program_requirements, p.related_tags, p.possible_careers, 
            p.program_application_fee, p.program_required_documents , p.university_id,

			u.university_name
     FROM programs p
	 left join universities u on p.university_id = u.university_id
     WHERE program_id = $1`,
		program_id)
	if err != nil {
		return []Types.Programe{}, err
	}
	defer rows.Close()
	var programList []Types.Programe
	for rows.Next() {
		var program Types.Programe
		var requirements pgtype.JSONB
		var relatedTags pgtype.JSONB
		var possibleCareers pgtype.JSONB
		var requiredDocuments pgtype.JSONB
		serr := rows.Scan(&program.Id, &program.Name, &program.PFee, &program.Duration, &program.SessionIntake, &program.Description, &requirements, &relatedTags, &possibleCareers, &program.AFee, &requiredDocuments, &program.UniversityCode, &program.UniversityName)
		if serr != nil {
			return []Types.Programe{}, serr
		}

		json.Unmarshal(requirements.Bytes, &program.Requirements)
		json.Unmarshal(relatedTags.Bytes, &program.RelatedTags)
		json.Unmarshal(relatedTags.Bytes, &program.RelatedTags)
		json.Unmarshal(possibleCareers.Bytes, &program.PossibleCareers)
		json.Unmarshal(requiredDocuments.Bytes, &program.RequiredDocs)
		programList = append(programList, program)
	}
	return programList, nil
}

func (p *StudentsAppStore) GetStudentProfileDetails(ctx context.Context, studentId string) (studentprofile Types.StudentProfile, err error) {
	var studentDetails Types.StudentProfile

	qerr := p.DB.QueryRowContext(ctx, `
    SELECT 
        sp.first_name, sp.middle_name, sp.last_name, sp.dob, sp.gender, sp.marrital_status,
        sp.nationality, sp.passport_number, sp.passport_expiry, sp.passport_issue,
        
        sc.email, sc.phone_number, sc.whatsapp_number, 
        sc.permanent_address, sc.street_address, sc.city, sc.state_province, sc.zip_postal_code,
        sc.emergency_contact_name, sc.emergency_relationship, sc.emergency_phone,
        
        se.school_name, se.curriculum, se.graduation_year, se.cummulative_score,
        se.language_type, se.language_overall_score, se.language_reading, 
        se.language_writting, se.language_speaking, se.language_listening,
        
		spref.primary_career_interest, spref.degree_level, spref.preferred_start_date, spref.annual_budget, spref.scholarship_interest

		FROM student_profile sp 
		LEFT JOIN student_contact sc ON sp.student_id = sc.student_id 
		LEFT JOIN student_education se ON sp.student_id = se.student_id
		LEFT JOIN students_preferences spref ON sp.student_id = spref.student_id 
		WHERE sp.student_id = $1

`, studentId).Scan(

		// profile fields
		&studentDetails.First_Name,
		&studentDetails.Middle_Name,
		&studentDetails.Last_Name,
		&studentDetails.Dob,
		&studentDetails.Gender,
		&studentDetails.Marrital_Status,
		&studentDetails.Nationality,
		&studentDetails.Passport_Number,
		&studentDetails.Passport_Expiry,
		&studentDetails.Place_Of_Issue,

		// contact fields
		&studentDetails.Email,
		&studentDetails.Phone,
		&studentDetails.WhatsApp,
		&studentDetails.Address,
		&studentDetails.Street,
		&studentDetails.City,
		&studentDetails.Province,
		&studentDetails.PostCode,
		&studentDetails.EmergencyName,
		&studentDetails.EmergencyRelation,
		&studentDetails.EmergencyPhone,

		// education fields
		&studentDetails.SchoolName,
		&studentDetails.Curriculum,
		&studentDetails.GraduationYear,
		&studentDetails.CummulativeScore,
		&studentDetails.LanguageType,
		&studentDetails.LanguageOverallScore,
		&studentDetails.LanguageReading,
		&studentDetails.LanguageWritting,
		&studentDetails.LanguageSpeaking,
		&studentDetails.LanguageListening,

		//student preferences
		&studentDetails.PrimaryInterest,
		&studentDetails.Degree,
		&studentDetails.PrefferedDate,
		&studentDetails.AnnualBudget,
		&studentDetails.Schalarship,
	)

	if qerr != nil {
		return studentDetails, qerr
	}

	fmt.Println(studentDetails)

	return studentDetails, nil

}

func (p *StudentsAppStore) UpdateProfile(ctx context.Context, student_id string, fieldName string, fieldValue string) error {

	query := fmt.Sprintf("UPDATE student_profile SET %s = $1 WHERE student_id = $2", fieldName)
	_, qerr := p.DB.ExecContext(ctx, query, fieldValue, student_id)

	if qerr != nil {
		return qerr
	}
	return nil
}

func (p *StudentsAppStore) UpdateContact(ctx context.Context, student_id string, fieldName string, fieldValue string) error {

	query := fmt.Sprintf("UPDATE student_contact SET %s = $1 WHERE student_id = $2", fieldName)
	_, qerr := p.DB.ExecContext(ctx, query, fieldValue, student_id)

	if qerr != nil {
		return qerr
	}
	return nil
}
func (p *StudentsAppStore) UpdateEducation(ctx context.Context, student_id string, fieldName string, fieldValue string) error {

	query := fmt.Sprintf("UPDATE student_education SET %s = $1 WHERE student_id = $2", fieldName)
	_, qerr := p.DB.ExecContext(ctx, query, fieldValue, student_id)

	if qerr != nil {
		return qerr
	}
	return nil
}

func (p *StudentsAppStore) UpdatePreferences(ctx context.Context, student_id string, fieldName string, fieldValue string) error {

	query := fmt.Sprintf("UPDATE students_preferences SET %s = $1 WHERE student_id = $2", fieldName)
	_, qerr := p.DB.ExecContext(ctx, query, fieldValue, student_id)

	if qerr != nil {
		return qerr
	}
	return nil
}

func (p *StudentsAppStore) GetstudentsDocuments(ctx context.Context, student_id string) ([]Types.UploadDocument, error) {
	var documentsList []Types.UploadDocument
	// err := p.DB.QueryRowContext(ctx,
	// 	"SELECT cv_mime_type,cv_status,cv_name,passport_mime_type,passport_status,passport_name,identity_mime_type,identity_status,identity_name,high_school_mime_type,high_school_status,high_school_name,language_proefficiency_mime_type,language_proefficiency_status,language_proefficiency_name,cover_letter_mime_type,cover_letter_status,cover_letter_name,motivation_letter_mime_type,motivation_letter_status,motivation_letter_name from students_documents where student_id = $1 ", student_id).Scan(
	// 	&documentsList.Cv.MimeType, &documentsList.Cv.Status, &documentsList.Cv.Name,
	// 	&documentsList.Passport.MimeType, &documentsList.Passport.Status, &documentsList.Passport.Name,
	// 	&documentsList.Identity.MimeType, &documentsList.Identity.Status, &documentsList.Identity.Name,
	// 	&documentsList.Highschool.MimeType, &documentsList.Highschool.Status, &documentsList.Highschool.Name,
	// 	&documentsList.Language.MimeType, &documentsList.Language.Status, &documentsList.Language.Name,
	// 	&documentsList.Coverletter.MimeType, &documentsList.Coverletter.Status, &documentsList.Coverletter.Name,
	// 	&documentsList.Motivationletter.MimeType, &documentsList.Motivationletter.Status, &documentsList.Motivationletter.Name,
	// )
	rows, err := p.DB.QueryContext(ctx,
		`SELECT	
		document_id,
		document,
		document_name,
		document_file_name,
		document_type,
		document_status
		from students_documents where student_id = $1 
		`, student_id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var document Types.UploadDocument
		scanerr := rows.Scan(&document.Id, &document.Data, &document.Name, &document.Document, &document.MimeType, &document.Status)

		if scanerr != nil {
			return nil, scanerr
		}
		documentsList = append(documentsList, document)

	}
	return documentsList, nil
}

func (p *StudentsAppStore) UploadStudentDocuments(ctx context.Context, studentDocs Types.UploadDocument, documentBytes []byte, student_id string) error {
	// fmt.Println(studentDocs, string(documentBytes), student_id)
	// doccolumn := studentDocs.Document
	// docMimecolumn := studentDocs.Document + "_mime_type"
	// docNamecolumn := studentDocs.Document + "_name"
	// docstatuscolumn := studentDocs.Document + "_status"

	// query2 := fmt.Sprintf("UPDATE students_documents SET  %s = $1, %s = $2, %s = $3, %s = $4 where student_id = $5", doccolumn, docNamecolumn, docMimecolumn, docstatuscolumn)
	// fmt.Println(query2)
	// query := fmt.Sprintf("UPDATE students_documents SET  document_type = $1,  = $2, %s = $3, %s = $4 where student_id = $5",  docNamecolumn, docMimecolumn, docstatuscolumn)

	_, err := p.DB.ExecContext(ctx, `insert into students_documents (document, document_file_name,document_name, document_type, document_status, student_id) values ($1,$2,$3,$4,$5,$6)`, documentBytes, studentDocs.Name, studentDocs.Document, studentDocs.MimeType, studentDocs.Status, student_id)
	if err != nil {
		return err
	}
	return nil
}

func (p *StudentsAppStore) GetStudentsDocument(ctx context.Context, studentId string, document_id string) (Types.Documents, error) {
	var Document Types.Documents

	// column := documentname
	// columnmime := documentname + "_mime_type"
	// columnname := documentname + "_name"

	// query := fmt.Sprintf(
	// 	`SELECT %s ,%s, %s FROM students_documents WHERE student_id = $1`,
	// 	columnname, column, columnmime,
	// )

	err := p.DB.QueryRowContext(ctx, `select document_name,document,document_type from students_documents where document_id = $1 and student_id = $2`, document_id, studentId).Scan(&Document.Name, &Document.Data, &Document.MimeType)

	if err != nil {
		return Types.Documents{}, err
	}

	return Document, nil
}

func (p *StudentsAppStore) UploadApplicationReceipt(ctx context.Context, student_id string, university_id string, program_id string, paid_amount string, receipt Types.UploadDocument) error {

	application_id, aterr := Tokens.GenerateToken(10)
	receipt_id, rterr := Tokens.GenerateToken(10)

	if rterr != nil {
		return rterr
	}

	if aterr != nil {
		return aterr
	}

	file, qtberr := Helper.Base64ToBytes(receipt.Data)

	if qtberr != nil {
		return qtberr
	}

	_, qerr := p.DB.ExecContext(ctx, "INSERT INTO university_app_receipts (student_id,university_id,program_id, paid_amount, application_id,receipt_id,receipt,mime_type,receipt_name,receipt_status) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		student_id,
		university_id,
		program_id,
		paid_amount,
		application_id,
		receipt_id,
		file,
		receipt.MimeType,
		receipt.Name,
		receipt.Status)

	if qerr != nil {
		return qerr
	}

	// _, q2err := dbctx.ExecContext(ctx, "INSERT INTO university_applications student_id,application_id,university_id,application_status,decision_status,feepaid",
	// 	student_id,
	// 	application_id,
	// 	university_id,
	// 	"submitted",
	// 	"pending",
	// 	"unpaid",
	// )

	// if q2err != nil {
	// 	dbctx.Rollback()
	// 	return q2err
	// }

	return nil
}

func (p StudentsAppStore) ApplyToUniversity(ctx context.Context, student_id string, university_id string, program_id string) error {

	application_id, tokenerr := Tokens.GenerateToken(10)

	if tokenerr != nil {
		return tokenerr
	}

	_, qerr := p.DB.ExecContext(ctx, "INSERT INTO university_applications (student_id,application_id,university_id,program_id,application_status,decision_status) values ($1,$2,$3,$4,$5,$6)", student_id, application_id, university_id, program_id, "applied", "pending")

	if qerr != nil {
		return qerr
	}

	return nil
}

func (p StudentsAppStore) GetApplicationsData(ctx context.Context, student_id string) ([]Types.ApplicationData, error) {

	var listofapplications []Types.ApplicationData

	rows, dberr := p.DB.QueryContext(ctx,
		`select uap.student_id, uap.application_id, 
	uap.application_status, uap.decision_status,
	uap.university_id, uap.program_id,TO_CHAR(uap.created_at,'DD Month,YYYY'),
	pr.program_name, pr.university_id , un.university_name
	from university_applications uap left join programs pr on pr.program_id = uap.program_id
	left join universities un on pr.university_id = un.university_id where uap.student_id = $1`, student_id)

	if dberr != nil {
		return nil, dberr
	}

	defer rows.Close()

	for rows.Next() {
		var application Types.ApplicationData

		scanerr := rows.Scan(&application.StudentId, &application.ApplicationId, &application.ApplicationStatus,
			&application.DecisionStatus, &application.UniversityId, &application.ProgramId, &application.ApplicationDate,
			&application.ProgramName, &application.UniversityId, &application.UniversityName)

		if scanerr != nil {
			return nil, scanerr
		}
		listofapplications = append(listofapplications, application)
	}
	return listofapplications, nil
}

func (p StudentsAppStore) VerifyApplication(ctx context.Context, student_id string, program_id string, university_id string) (Types.ExistsRow, error) {

	var rowsExist Types.ExistsRow

	dberr := p.DB.QueryRowContext(ctx, `SELECT
		EXISTS (
			SELECT 1
			FROM university_app_receipts
			WHERE student_id = $1
			AND program_id = $2
			AND university_id = $3
		) AS has_receipt,
		EXISTS (
			SELECT 1
			FROM university_applications
			WHERE student_id = $1
			AND program_id = $2
			AND university_id = $3
		) AS has_application`, student_id, program_id, university_id).Scan(&rowsExist.HasReceipt, &rowsExist.HasApplication)

	if dberr != nil {
		return Types.ExistsRow{}, dberr
	}
	return rowsExist, nil
}

func (p StudentsAppStore) ShortListProgram(ctx context.Context, student_id string, program_id string, university_id string) (int, error) {
	var id int
	dberr := p.DB.QueryRowContext(ctx, `INSERT into shortlisted_programs (student_id,program_id,university_id) VALUES ($1,$2,$3) returning id`, student_id, program_id, university_id).Scan(&id)

	if dberr != nil {
		return 0, dberr
	}
	return id, nil
}

func (p StudentsAppStore) GetShortListProgram(ctx context.Context, student_id string) (shortListedPrograms []Types.ShortListProgram, err error) {

	var listOfPrograms []Types.ShortListProgram
	rows, dberr := p.DB.QueryContext(ctx, `SELECT id,student_id,program_id,university_id from shortlisted_programs where student_id = $1`, student_id)

	if dberr != nil {
		return nil, dberr
	}

	for rows.Next() {
		var Program Types.ShortListProgram
		scanerr := rows.Scan(&Program.ShortListId, &Program.Student, &Program.Programe, &Program.University)

		if scanerr != nil {
			return nil, scanerr
		}

		listOfPrograms = append(listOfPrograms, Program)
	}
	return listOfPrograms, nil
}

func (p StudentsAppStore) DeleteShortListProgram(ctx context.Context, student_id string, shortListId string) error {
	_, dberr := p.DB.ExecContext(ctx, `Delete from shortlisted_programs where student_id = $1 and id = $2`, student_id, shortListId)

	if dberr != nil {
		return dberr
	}
	return nil
}
