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
	_, queryerror := p.DB.ExecContext(ctx, "INSERT INTO credentials (email,hashed_password,role,name,id) VALUES ($1,$2,$3,$4,$5)", admin.Email, hashedPassword, "role", admin.Name, admin.Id)
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
	// create an invite token
	// save to db
	// write a message
	// var message string
	token, terr := Tokens.GenerateToken(10)
	var generatedData Types.LinkGenerated
	fmt.Println(adminid)
	if terr != nil {
		slog.Info("There was a session token generation error", "error", terr)
		return &generatedData, terr
	}
	_, qerr := p.DB.ExecContext(ctx, "INSERT INTO school_invites (token,sys_admin,school_name,school_email,status) VALUES($1,$2,$3,$4,$5)", token, adminid, inviteData.Name, inviteData.Email, "pending")

	if qerr != nil {
		slog.Info("There was a session token generation error", "error", qerr)
		return &generatedData, qerr
	}

	// message = fmt.Sprintf("%s has been added to your invitations Contact email: %s", inviteData.Name, inviteData.Email)

	// generatedData.Messsage = message
	generatedData.Token = token
	generatedData.SchoolInvite = inviteData
	return &generatedData, nil
}

func (p *SysAdminStore) GetInviteData(ctx context.Context, token string) (string, error) {
	var email string
	qerr := p.DB.QueryRowContext(ctx, "SELECT school_email from school_invites WHERE token = $1", token).Scan(&email)
	if qerr != nil {
		slog.Info("Query Error", "message", "there as an error querying invite token", "error", qerr)
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

func (p *SysAdminStore) GetInvites(ctx context.Context, limit int, offlimit int) ([]Types.SchoolInformation, error) {
	status := "completed"
	var created_at time.Time

	arrayOfRows := []Types.SchoolInformation{}
	rows, qerr := p.DB.QueryContext(ctx, `SELECT school_name, school_email, admin_name, school_phone, school_country, school_id, school_curriculum, school_branch, school_city, created_at FROM school_invites WHERE created_at >= NOW() - INTERVAL '1 month' AND status = $1 ORDER BY created_at ASC LIMIT $2 OFFSET  $3`, status, limit, offlimit)
	if qerr != nil {
		slog.Info("db Error", "message", "Error in fetching analytics for invite dashboard", "error", qerr)
		return []Types.SchoolInformation{}, qerr
	}

	defer rows.Close()

	for rows.Next() {
		var eachApplication Types.SchoolInformation
		err := rows.Scan(&eachApplication.School, &eachApplication.Email, &eachApplication.Admin, &eachApplication.Phone, &eachApplication.Country, &eachApplication.Id, &eachApplication.Curriculam, &eachApplication.Branch, &eachApplication.City, &created_at)
		if err != nil {
			slog.Info("error in populating type", "message", "row could not be converted into struct", "error", err)
			return []Types.SchoolInformation{}, err
		}
		priority := timecheck.CheckAppPriority(created_at)
		eachApplication.Priority = priority
		arrayOfRows = append(arrayOfRows, eachApplication)
	}

	return arrayOfRows, nil
}

func (p *SysAdminStore) RespondToSchoolInvite(ctx context.Context, token string, status string) (string, error) {
	var schoolInformation Types.SchoolInformation
	now := time.Now().UTC()

	tx, terr := p.DB.BeginTx(ctx, nil)
	if terr != nil {
		slog.Info("Db Error", "message", "there was an error in startin transaction ", "error", terr)
		return "", terr
	}
	qerr := tx.QueryRowContext(ctx, "SELECT school_id,school_email,school_name,admin_name,school_phone,school_country,school_curriculum,school_city,school_branch,token from school_invites WHERE token = $1", token).Scan(&schoolInformation.Id, &schoolInformation.Email, &schoolInformation.School, &schoolInformation.Admin, &schoolInformation.Phone, &schoolInformation.Country, &schoolInformation.Curriculam, &schoolInformation.City, &schoolInformation.Branch, &schoolInformation.Token)
	if qerr != nil {
		tx.Rollback()
		slog.Info("Db Error", "message", "there was an error querying the appication ", "error", qerr)
		return "", qerr
	}
	_, uerr := tx.ExecContext(ctx, "UPDATE school_invites SET status = $1, approved_date = $2 WHERE token = $3", status, now, token)
	if uerr != nil {
		tx.Rollback()
		slog.Info("Db Error", "message", "there was an error updating the status ", "error", uerr)
		return "", uerr
	}
	if status == "approved" {
		_, movedberr := tx.ExecContext(ctx, "INSERT INTO schools (school_id,school_email,school_name,admin_name,school_phone,school_country,school_curriculum,school_city,school_branch,token,status,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", &schoolInformation.Id, &schoolInformation.Email, &schoolInformation.School, &schoolInformation.Admin, &schoolInformation.Phone, &schoolInformation.Country, &schoolInformation.Curriculam, &schoolInformation.City, &schoolInformation.Branch, &schoolInformation.Token, status, now)
		if movedberr != nil {
			slog.Info("Db Error", "message", "there was an error moving to  school db ", "error", movedberr)
			return "", uerr
		}
	}
	txerr := tx.Commit()
	if txerr != nil {
		slog.Info("Transaction Error", "message", "there was an error manipulating db ", "error", txerr)
		return "", txerr
	}

	return schoolInformation.Email, nil
}

func (p *SysAdminStore) SaveSchoolAdminCredentials(ctx context.Context, email string, password string, token string) error {

	_, err := p.DB.ExecContext(ctx, "INSERT INTO school_credentials (email,hashed_password,token) VALUES ($1,$2,$3)", email, password, token)
	if err != nil {
		slog.Info("Db Error", "message", "there was an error saving credential to db ", "error", err)
		return err
	}
	return nil
}
