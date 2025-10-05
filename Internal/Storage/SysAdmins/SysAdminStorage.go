package SysAdminStorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	_ "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
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
	_, qerr := p.DB.ExecContext(ctx, "INSERT INTO school_invites (token,admin,name,email,status) VALUES($1,$2,$3,$4,$5)", token, adminid, inviteData.Name, inviteData.Email, "pending")

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
