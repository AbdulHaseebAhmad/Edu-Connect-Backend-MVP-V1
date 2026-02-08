package SysAdminStorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	_ "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	Emailhelper "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Emails"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	timecheck "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Time"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
)

type SysAdminStore struct {
	*Postgress.Postgress
}

func NewSysAdminStore(pg *Postgress.Postgress) *SysAdminStore {
	return &SysAdminStore{pg}
}

func (p *SysAdminStore) SysAdminLogin(ctx context.Context, admin Types.SysAdminLogin) (sessionToken string, csrfToken string, sysadminauth *Types.SysAdminAuthenticated, err error) {
	var hashedPassword string

	SysAdminAut := Types.SysAdminAuthenticated{
		Authenticated: true,
		Status:        true,
	}

	sessionId, sessionerr := Tokens.GenerateToken(10)
	if sessionerr != nil {
		slog.Info("There was a session ID token generation error", "error", sessionerr)
		return "", "", &Types.SysAdminAuthenticated{}, sessionerr
	}
	queryerr := p.DB.QueryRowContext(ctx, "SELECT hashed_password,role,email,name,id from credentials WHERE email = $1", admin.Email).Scan(&hashedPassword, &SysAdminAut.Role, &SysAdminAut.Email, &SysAdminAut.Name, &SysAdminAut.Id)
	if queryerr != nil {
		slog.Info("There was an error in querying hashed password from db", "error", queryerr)
		return "", "", &Types.SysAdminAuthenticated{}, queryerr
	}

	passwordmatch, matcherr := HashPassword.Unhashpassword(admin.Password, hashedPassword)

	if matcherr != nil {
		slog.Info("There was an internal error", "error", "Hashin algorithim error")
		return "", "", &Types.SysAdminAuthenticated{}, errors.New("authentication Error")
	}
	if !passwordmatch {
		slog.Info("There was an auth error", "error", "Password/Email is wrong")
		return "", "", &Types.SysAdminAuthenticated{}, errors.New("authentication Error")
	}
	session_token, stokenerr := Tokens.GenerateToken(10)
	if stokenerr != nil {
		slog.Info("There was a session token generation error", "error", stokenerr)
		return "", "", &Types.SysAdminAuthenticated{}, stokenerr
	}
	csrf_token, csrftokenerr := Tokens.GenerateToken(10)
	if csrftokenerr != nil {
		slog.Info("There was a csrf token generation error", "error", stokenerr)
		return "", "", &Types.SysAdminAuthenticated{}, csrftokenerr
	}

	SysAdminAut.CsrfToken = csrf_token

	_, insertqerr := p.DB.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, role,credential_id,session_id)  VALUES ($1, $2, $3, $4, $5, $6)", session_token, csrf_token, SysAdminAut.Email, SysAdminAut.Role, SysAdminAut.Id, sessionId)
	if insertqerr != nil {
		slog.Info("There was an error inserting data to db", "error", insertqerr)
		return "", "", &Types.SysAdminAuthenticated{}, nil
	}

	return session_token, csrf_token, &SysAdminAut, nil
}

func (p *SysAdminStore) SysAdminSignup(ctx context.Context, admin Types.SysAdminSignup) (err error) {
	hashedPassword, hasherror := HashPassword.Hashpassword(admin.Password)
	if hasherror != nil {
		return hasherror
	}
	_, queryerror := p.DB.ExecContext(ctx, "INSERT INTO credentials (email,hashed_password,role,name,id) VALUES ($1,$2,$3,$4,$5)", admin.Email, hashedPassword, "sys_admin", admin.Name, admin.Id)
	if queryerror != nil {
		slog.Info("There was an error in querying db", "error", queryerror)
		return queryerror
	}
	return nil
}

func (p *SysAdminStore) AuthorizeSysAdmin(ctx context.Context, sessionToken string, csrfToken string) (string, bool) {
	var csrf string
	var id string
	qerr := p.DB.QueryRowContext(ctx, "SELECT csrf_token ,credential_id FROM sessions WHERE session_token = $1 ", sessionToken).Scan(&csrf, &id)
	if qerr != nil {
		fmt.Println(qerr)
		if errors.Is(qerr, sql.ErrNoRows) {
			return "", false
		}
		return "", false
	}
	fmt.Println("id", id)
	return id, csrf == csrfToken
}

func (p *SysAdminStore) GenerateInvite(ctx context.Context, adminid Types.SysAdminId, inviteData Types.SchoolInvite) (*Types.LinkGenerated, error) {
	var generatedData Types.LinkGenerated
	var invitaion_id string
	qerr := p.DB.QueryRowContext(ctx, "INSERT INTO school_invites (sys_admin,school_name,school_email,invitation_status) VALUES($1,$2,$3,$4) returning invitation_id", adminid, inviteData.Name, inviteData.Email, "pending").Scan(&invitaion_id)

	if qerr != nil {
		slog.Info("There was a session token generation error", "error", qerr)
		return &generatedData, qerr
	}

	generatedData.Token = invitaion_id
	generatedData.SchoolInvite = inviteData
	return &generatedData, nil
}

func (p *SysAdminStore) GetInviteData(ctx context.Context, invitation_id string) (string, error) {
	var email string
	qerr := p.DB.QueryRowContext(ctx, "SELECT school_email from school_invites WHERE invitation_id = $1", invitation_id).Scan(&email)
	if qerr != nil {
		slog.Info("Query Error", "message", "there as an error querying invite school_id", "error", qerr)
		return "", qerr
	}
	return email, nil
}

func (p *SysAdminStore) GetInvitesAnalytics(ctx context.Context) (Types.InvitesAnalytics, error) {
	var analytics Types.InvitesAnalytics
	pending := "pending"
	completed := "completed"
	approved := "approved"
	err := p.DB.QueryRowContext(ctx, "SELECT COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '1 month' ),COUNT(*) FILTER (WHERE status = $1 AND created_at >= NOW() - INTERVAL '1 month'), COUNT(*) FILTER (WHERE status = $2 AND created_at >= NOW() - INTERVAL '1 month'), COUNT(*) FILTER (WHERE status = $3 AND created_at >= NOW() - INTERVAL '1 month') FROM school_invites", pending, completed, approved).Scan(&analytics.Total, &analytics.Pending, &analytics.Completed, &analytics.Approved)

	analytics.ApprovalRate = (float64(analytics.Approved) / float64(analytics.Total) * 100)
	analytics.AcceptanceRate = (float64(analytics.Completed) / float64(analytics.Total) * 100)

	if err != nil {
		slog.Info("db Error", "message", "Error in fetching analytics for invite dashboard", "error", err)
		return Types.InvitesAnalytics{}, err
	}

	return analytics, nil
}

func (p *SysAdminStore) GetSchoolApplications(ctx context.Context, limit int, offlimit int) ([]Types.SchoolApp, error) {
	var created_at time.Time

	listOfApplications := []Types.SchoolApp{}
	rows, qerr := p.DB.QueryContext(ctx, `SELECT school_name, school_email, application_id, status,to_char(completed_at,'DD MM YYYY') FROM school_applications WHERE created_at >= NOW() - INTERVAL '3 month'  ORDER BY created_at ASC LIMIT $1 OFFSET  $2`, limit, offlimit)
	if qerr != nil {
		slog.Info("db Error", "message", "Error in fetching analytics for invite dashboard", "error", qerr)
		return nil, qerr
	}

	defer rows.Close()

	for rows.Next() {
		var eachApplication Types.SchoolApp
		err := rows.Scan(&eachApplication.SchoolName, &eachApplication.SchoolEmail, &eachApplication.ApplicationId, &eachApplication.ApplicationStatus, &eachApplication.Date)
		if err != nil {
			slog.Info("error in populating type", "message", "row could not be converted into struct", "error", err)
			return nil, err
		}
		priority := timecheck.CheckAppPriority(created_at)
		eachApplication.Priority = priority
		listOfApplications = append(listOfApplications, eachApplication)
	}

	return listOfApplications, nil
}

func (p *SysAdminStore) GetSchoolApplication(ctx context.Context, application_id string) (Types.SchoolApplication, error) {
	var applicationDetail Types.SchoolApplication

	qerr := p.DB.QueryRowContext(ctx, `SELECT  application_id,school_id,registration_code,school_name, school_email,admin_name,school_phone,school_country,school_curriculum,school_branch,school_city FROM school_applications WHERE application_id = $1`, application_id).Scan(&applicationDetail.ApplicationId, &applicationDetail.SchoolId, &applicationDetail.RegistrationCode, &applicationDetail.SchoolName, &applicationDetail.SchoolEmail, &applicationDetail.AdminName, &applicationDetail.SchoolPhone, &applicationDetail.SchoolCountry, &applicationDetail.SchoolCurriculam, &applicationDetail.SchoolBranch, &applicationDetail.SchoolCity)
	if qerr != nil {
		slog.Info("db Error", "message", "Error in fetching analytics for invite dashboard", "error", qerr)
		return Types.SchoolApplication{}, qerr
	}

	return applicationDetail, nil
}

func (p *SysAdminStore) GetAllInvites(ctx context.Context) (invitesList []Types.Invite, err error) {
	var invites []Types.Invite
	rows, dberr := p.DB.QueryContext(ctx, `SELECT invitation_id,invitation_status, school_email, school_name,TO_CHAR(created_at, 'DD Month YYYY') from school_invites`)

	if dberr != nil {
		return nil, dberr
	}

	for rows.Next() {
		var invite Types.Invite
		scanerr := rows.Scan(&invite.InviteId, &invite.InviteStatus, &invite.SchoolEmail, &invite.SchoolName, &invite.Date)

		if scanerr != nil {
			return nil, scanerr
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

func (p *SysAdminStore) RespondToSchoolInvite(ctx context.Context, application_id string, status string) (Types.SchoolInformation, string, error) {
	var schoolInformation Types.SchoolInformation
	now := time.Now().UTC()
	var generatedPassword string
	tx, terr := p.DB.BeginTx(ctx, nil)
	if terr != nil {
		slog.Info("Db Error", "message", "there was an error in startin transaction ", "error", terr)
		return Types.SchoolInformation{}, "", terr
	}

	qerr := tx.QueryRowContext(ctx, "SELECT school_id,registration_code,school_name,school_email,admin_name,school_phone,school_country,school_curriculum,school_branch,school_city from school_applications WHERE application_id = $1", application_id).Scan(
		&schoolInformation.SchoolId, &schoolInformation.RegistrationCode, &schoolInformation.School, &schoolInformation.Email, &schoolInformation.Admin, &schoolInformation.Phone, &schoolInformation.Country, &schoolInformation.Curriculam, &schoolInformation.Branch, &schoolInformation.City)
	if qerr != nil {
		tx.Rollback()
		slog.Info("Db Error", "message", "there was an error querying the appication ", "error", qerr)
		return Types.SchoolInformation{}, "", qerr
	}

	_, uerr := tx.ExecContext(ctx, "UPDATE school_applications SET status = $1, responded_at = $2 WHERE application_id = $3", status, now, application_id)

	if uerr != nil {
		tx.Rollback()
		slog.Info("Db Error", "message", "there was an error updating the status ", "error", uerr)
		return Types.SchoolInformation{}, "", uerr
	}

	if status == "approved" {
		username := schoolInformation.School + schoolInformation.Branch
		sys_email := Emailhelper.GenerateEmails(schoolInformation.Email, "school admin")
		schoolInformation.Sys_Eamil = schoolInformation.Email
		schoolInformation.Username = username

		_, movedberr := tx.ExecContext(ctx, "INSERT INTO schools (registration_code,school_email,school_name,admin_name,school_phone,school_country,school_curriculum,school_city,school_branch,school_id,status,created_at,username,sys_email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)", &schoolInformation.RegistrationCode, &schoolInformation.Email, &schoolInformation.School, &schoolInformation.Admin, &schoolInformation.Phone, &schoolInformation.Country, &schoolInformation.Curriculam, &schoolInformation.City, &schoolInformation.Branch, &schoolInformation.SchoolId, status, now, username, sys_email)
		if movedberr != nil {
			slog.Info("Db Error", "message", "there was an error moving to  school db ", "error", movedberr)
			return Types.SchoolInformation{}, "", uerr
		}

		generatePassword, gerr := Tokens.GenerateToken(10)
		if gerr != nil {
			slog.Info("Password error", "message", "There was an error creating password", "error", gerr)
			return Types.SchoolInformation{}, "", gerr
		}

		generatedPassword = generatePassword
		securePassword, herr := HashPassword.Hashpassword(generatePassword)
		if herr != nil {
			slog.Info("Password error", "message", "There was an error hashing password", "error", herr)
			return Types.SchoolInformation{}, "", herr
		}

		_, serr := p.DB.ExecContext(ctx, "INSERT INTO school_credentials (school_id,sys_email,hashed_password,name,role) VALUES ($1,$2,$3,$4,$5)", schoolInformation.SchoolId, sys_email, securePassword, schoolInformation.School, "schooladmin")
		if serr != nil {
			slog.Info("Db Error", "message", "there was an error saving credential to db ", "error", serr)
			return Types.SchoolInformation{}, "", serr
		}

		sserr := p.DB.QueryRowContext(ctx, "INSERT INTO school_codes (school_id) VALUES ($1) returning school_code", schoolInformation.SchoolId).Scan(&schoolInformation.Code)
		if sserr != nil {
			slog.Info("Db Error", "message", "there was an error saving school codes to db ", "error", sserr)
			return Types.SchoolInformation{}, "", sserr
		}

	}
	txerr := tx.Commit()
	if txerr != nil {
		slog.Info("Transaction Error", "message", "there was an error manipulating db ", "error", txerr)
		return Types.SchoolInformation{}, "", txerr
	}

	return schoolInformation, generatedPassword, nil
}

func (p *SysAdminStore) GetAnalyticsList(ctx context.Context) (lists Types.AnalyticsList, err error) {
	var analyticsList Types.AnalyticsList
	rows, qerr := p.DB.QueryContext(ctx, `SELECT school_name, school_email, admin_name, school_phone, school_country, registration_code, school_curriculum, school_branch, school_city, created_at,school_id,completed_date,approved_date,status FROM school_invites WHERE created_at >= NOW() - INTERVAL '1 month' ORDER BY created_at ASC LIMIT $1 OFFSET  $2`, 100, 0)

	if qerr != nil {
		slog.Info("There was an error in query", "message", "Could not query analytics list", "error", qerr)
		return Types.AnalyticsList{}, qerr

	}

	defer rows.Close()

	for rows.Next() {
		var eachApplication Types.SchoolInformation
		err := rows.Scan(&eachApplication.School, &eachApplication.Email, &eachApplication.Admin, &eachApplication.Phone, &eachApplication.Country, &eachApplication.RegistrationCode, &eachApplication.Curriculam, &eachApplication.Branch, &eachApplication.City, &eachApplication.CreatedAt, &eachApplication.SchoolId, &eachApplication.CompletedAt, &eachApplication.ApprovedAt, &eachApplication.Status)
		if err != nil {
			slog.Info("error in populating type", "message", "row could not be converted into struct", "error", err)
			return Types.AnalyticsList{}, err
		}
		switch eachApplication.Status {
		case "pending":
			analyticsList.PendingInvites = append(analyticsList.PendingInvites, eachApplication)
		case "approved":
			analyticsList.ApprovedApplications = append(analyticsList.ApprovedApplications, eachApplication)
		case "completed":
			analyticsList.PendingApplications = append(analyticsList.PendingApplications, eachApplication)
		default:
			slog.Info("Unknown status found", "school", eachApplication.School, "status", eachApplication.Status)
		}
	}

	if err := rows.Err(); err != nil {
		return Types.AnalyticsList{}, err
	}

	return analyticsList, nil

}

func (p *SysAdminStore) GetStudentsRegistry(
	ctx context.Context,
	status string,
) ([]Types.StudentsRegistry, error) {

	var (
		rows *sql.Rows
		err  error
	)

	if status == "all" {
		rows, err = p.DB.QueryContext(
			ctx,
			`
             SELECT sap.first_name, sap.last_name, sap.email, sap.citizenship, sap.slug, sap.status,
			 TO_CHAR(sap.created_at,'DD Montth, YYYY')
            FROM students_applications sap
            `,
		)
	} else {
		rows, err = p.DB.QueryContext(
			ctx,
			`
            SELECT sap.first_name, sap.last_name, sap.email, sap.citizenship, sap.slug, sap.status,
			TO_CHAR(sap.created_at,'DD Montth, YYYY')
            FROM students_applications sap
            WHERE sap.status = $1
            `,
			status,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listOfStudents []Types.StudentsRegistry

	for rows.Next() {
		var student Types.StudentsRegistry

		if err := rows.Scan(
			&student.First_Name,
			&student.Last_Name,
			&student.Email,
			&student.Citizenship,
			&student.Slug,
			&student.Status,
			&student.Date,
		); err != nil {
			return nil, err
		}

		listOfStudents = append(listOfStudents, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listOfStudents, nil
}

func (p *SysAdminStore) RespondApplication(ctx context.Context, action string, id string) (string, string, error) {
	var email string
	var fname string
	var lname string
	var citizenship string
	var passport []byte
	var transcript []byte
	var passportType string
	var transcriptType string
	var school_code string
	tx, terr := p.DB.BeginTx(ctx, nil)
	if terr != nil {
		slog.Info("Db Error", "message", "there was an error in startin transaction ", "error", terr)
		return "", "", terr
	}
	err := tx.QueryRowContext(ctx, "UPDATE students_applications SET status = $1 where slug = $2 returning email,first_name,last_name,citizenship,passport,transcript,passport_mime_type,transcript_mime_type,school_code", action, id).Scan(&email, &fname, &lname, &citizenship, &passport, &transcript, &passportType, &transcriptType, &school_code)

	if err != nil {
		tx.Rollback()
		return "", "", nil
	}

	rpassword := ""
	hashedPassword := ""
	studentId := ""

	if action == "approved" {
		password, tokenerr := Tokens.GenerateToken(8)
		if tokenerr != nil {
			return "", "", nil
		}
		hpassword, perr := HashPassword.Hashpassword(password)
		if perr != nil {
			return "", "", nil
		}
		studentid, sterr := Tokens.GenerateToken(12)
		if sterr != nil {
			return "", "", nil
		}
		studentId = studentid
		hashedPassword = hpassword
		rpassword = password
		_, sqerr := tx.ExecContext(ctx,
			`INSERT INTO students_credentials (student_email, hashed_password, role, student_id,school_id) 
		 VALUES ($1, $2, $3, $4,$5)`,
			email, hashedPassword, "student", studentId, school_code)

		if sqerr != nil {
			tx.Rollback()
			return "", "", sqerr
		}
		_, tqerr := tx.ExecContext(ctx, `INSERT into student_profile (student_id,first_name,last_name,nationality) values ($1,$2,$3,$4)`, studentId, fname, lname, citizenship)

		if tqerr != nil {
			tx.Rollback()
			return "", "", tqerr
		}

		_, fqerr := tx.ExecContext(ctx, `INSERT into student_contact (student_id,email) values ($1,$2)`, studentId, email)

		if fqerr != nil {
			tx.Rollback()
			return "", "", fqerr
		}

		// _, fiqerr := tx.ExecContext(ctx, `INSERT into students_documents (student_id,passport,passport_name,passport_status,passport_mime_type,high_school,high_school_name,high_school_status,high_school_mime_type) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		// 	studentid, passport, "passport", "uploaded", passportType, transcript, "high school Diploma", "uploaded", transcriptType)

		_, fiqerr := tx.ExecContext(ctx, `
			INSERT INTO students_documents 
			(student_id, document, document_type, document_name, document_file_name, document_status)
			VALUES 
			($1,$2,$3,$4,$5,$6),
			($1,$7,$8,$9,$10,$11)
			`,
			studentid,
			passport, passportType, "passport", "passport", "uploaded",
			transcript, transcriptType, "high school diploma", "high school diploma", "uploaded",
		)

		if fiqerr != nil {
			tx.Rollback()
			return "", "", fiqerr
		}
		_, siqerr := tx.ExecContext(ctx, `INSERT into student_education (student_id) values ($1)`,
			studentid)

		if siqerr != nil {
			tx.Rollback()
			return "", "", siqerr
		}
		_, seqerr := tx.ExecContext(ctx, `INSERT into students_preferences (student_id) values ($1)`,
			studentid)

		if seqerr != nil {
			tx.Rollback()
			return "", "", seqerr
		}

	}

	txerr := tx.Commit()
	if txerr != nil {
		slog.Info("Transaction Error", "message", "there was an error manipulating db ", "error", txerr)
		return "", "", txerr
	}

	return email, rpassword, nil
}

func (p *SysAdminStore) GetStudentsDocument(ctx context.Context, studentId string, documentname string, documentmime string) (Types.Documents, error) {
	var Document Types.Documents

	var columnname string
	var columnmime string

	if documentname == "passport" {
		columnname = "passport"
		columnmime = "passport_mime_type"
	} else {
		columnname = "transcript"
		columnmime = "transcript_mime_type"
	}

	query := fmt.Sprintf(
		`SELECT %s ,%s FROM students_applications WHERE slug = $1`,
		columnname, columnmime,
	)

	err := p.DB.QueryRowContext(ctx, query, studentId).Scan(&Document.Data, &Document.MimeType)

	if err != nil {
		return Types.Documents{}, err
	}

	return Document, nil
}

func (p *SysAdminStore) GetAllReceipts(ctx context.Context) ([]Types.UniversityAppReceipt, error) {

	var receiptsList []Types.UniversityAppReceipt

	rows, dberr := p.DB.QueryContext(ctx, `select 
		uar.student_id, 
		sp.first_name,
		sc.student_email,
		uar.university_id,
		uar.application_id,
		uar.receipt_id,
		uar.receipt,
		uar.mime_type,
		uar.receipt_name,
		uar.receipt_status,
		uar.paid_amount,
		uar.program_id,
		TO_CHAR(uar.created_date, 'MM DD, YYYY'),

		u.university_name
	from university_app_receipts uar
	LEFT JOIN students_credentials sc on uar.student_id = sc.student_id
	LEFT JOIN universities u on uar.university_id = u.university_id
	LEFT JOIN student_profile sp on uar.student_id = sp.student_id
	`)

	if dberr != nil {
		return nil, dberr
	}

	for rows.Next() {
		var receipt Types.UniversityAppReceipt
		err := rows.Scan(&receipt.StudentID, &receipt.FirstName, &receipt.StudentEmail, &receipt.UniversityID, &receipt.ApplicationID, &receipt.ReceiptID,
			&receipt.Receipt, &receipt.MimeType, &receipt.ReceiptName, &receipt.ReceiptStatus, &receipt.PaidAmount, &receipt.ProgramId, &receipt.CreatedDate, &receipt.UniversityName)
		if err != nil {
			return nil, err
		}
		receiptsList = append(receiptsList, receipt)
	}
	return receiptsList, nil
}
func (p *SysAdminStore) GetReceiptDetails(ctx context.Context, student_id string) (Types.UniversityAppReceipt, error) {

	var receipt Types.UniversityAppReceipt

	dberr := p.DB.QueryRowContext(ctx, `select 
		uar.student_id, 
		sp.first_name,
		sc.student_email,
		uar.university_id,
		uar.application_id,
		uar.receipt_id,
		uar.receipt,
		uar.mime_type,
		uar.receipt_name,
		uar.receipt_status,
		uar.paid_amount,
		TO_CHAR(uar.created_date, 'MM DD, YYYY'),

		u.university_name
		from university_app_receipts uar
		LEFT JOIN students_credentials sc on uar.student_id = sc.student_id
		LEFT JOIN universities u on uar.university_id = u.university_id
		LEFT JOIN student_profile sp on uar.student_id = sp.student_id
		where uar.receipt_id = $1
	`, student_id).Scan(&receipt.StudentID, &receipt.FirstName, &receipt.StudentEmail, &receipt.UniversityID, &receipt.ApplicationID, &receipt.ReceiptID,
		&receipt.Receipt, &receipt.MimeType, &receipt.ReceiptName, &receipt.ReceiptStatus, &receipt.PaidAmount, &receipt.CreatedDate, &receipt.UniversityName)

	if dberr != nil {
		return Types.UniversityAppReceipt{}, dberr
	}

	return receipt, nil
}

func (p *SysAdminStore) RespondToReceipts(ctx context.Context, receipt_id string, status string) error {

	_, dberr := p.DB.ExecContext(ctx, `UPDATE university_app_receipts SET receipt_status = $1 where receipt_id = $2`, status, receipt_id)
	if dberr != nil {
		return nil

	}
	return nil
}

func (p *SysAdminStore) GetRegisteredStudents(ctx context.Context) ([]Types.RegisteredStudent, error) {

	var list []Types.RegisteredStudent
	rows, dberr := p.DB.QueryContext(ctx, `select sc.student_email, TO_CHAR(sc.created_date,'DD Month, YYYY'), sc.student_status,sc.student_id, sc.school_verified, sc.school_id,
	spd.first_name,spd.last_name 
	from students_credentials sc LEFT JOIN student_profile spd on sc.student_id = spd.student_id `)

	if dberr != nil {
		return nil, dberr
	}

	defer rows.Close()

	for rows.Next() {
		var registeredStudent Types.RegisteredStudent
		scanerr := rows.Scan(&registeredStudent.Email, &registeredStudent.CreatedDate, &registeredStudent.Status, &registeredStudent.StudentId, &registeredStudent.SchoolVerified, &registeredStudent.SchoolId, &registeredStudent.FirstName, &registeredStudent.LastName)
		if scanerr != nil {
			return nil, scanerr
		}
		list = append(list, registeredStudent)
	}

	return list, nil
}
