package SchoolAdminStorage

import (
	"context"
	"errors"
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
	queryerr := p.DB.QueryRowContext(ctx, "SELECT id,sys_email,hashed_password,role,name from school_credentials WHERE sys_email = $1", schooladmin.Email).Scan(&SchoolAdminAut.Id, &SchoolAdminAut.Email, &hashedPassword, &SchoolAdminAut.Role, &SchoolAdminAut.Name)
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

	_, insertqerr := p.DB.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, session_id, credential_id,role)  VALUES ($1, $2, $3, $4, $5,$6)", session_token, csrf_token, SchoolAdminAut.Email, sessionId, SchoolAdminAut.Id, SchoolAdminAut.Role)
	if insertqerr != nil {
		slog.Info("There was an error inserting data to db", "error", insertqerr)
		return "", "", &Types.SchoolAdminAuthenticated{}, nil
	}

	return session_token, csrf_token, &SchoolAdminAut, nil
}

func (p *SchoolAdminStore) ValidateLink(ctx context.Context, inviteToken string) (string, error) {
	var returnedToken string
	var returneTimeStamp string
	var status string
	// layout := "2006-01-02 15:04:05"

	qerr := p.DB.QueryRowContext(ctx, "SELECT token, created_at, status from school_invites WHERE token = $1", inviteToken).Scan(&returnedToken, &returneTimeStamp, &status)

	if qerr != nil {
		slog.Info("There was an error queryying", "error", qerr)
		return "", qerr
	}

	parsedTime, terr := time.Parse(time.RFC3339Nano, returneTimeStamp)

	if terr != nil {
		slog.Info("There was an error parsing time", "error", terr)
		return "", terr
	}

	expired := time.Now().Before(parsedTime)

	if !expired {
		return "", errors.New("token is expired")
	}
	return status, nil
}

func (p *SchoolAdminStore) SubmitInvite(ctx context.Context, schoolInfo Types.SchoolInformation, token string) (string, error) {
	var status string
	currentTime := time.Now()
	qerr := p.DB.QueryRowContext(ctx, "UPDATE school_invites SET status = $1, admin_name = $2, school_phone = $3, school_country = $4, school_id = $5, school_curriculum = $6, school_branch = $7, school_city = $8, completed_date = $9 WHERE token = $10 RETURNING status", "completed", schoolInfo.Admin, schoolInfo.Phone, schoolInfo.Country, schoolInfo.Id, schoolInfo.Curriculam, schoolInfo.Branch, schoolInfo.City, currentTime, token).Scan(&status)
	if qerr != nil {
		slog.Info("Query Error", "message", "there was an error saving school info to db", "error", qerr)
		return "", qerr
	}
	return "completed", nil
}
