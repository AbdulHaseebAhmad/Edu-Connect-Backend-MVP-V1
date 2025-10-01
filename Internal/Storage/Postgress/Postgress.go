package Postgress

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	Configurator "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Config"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
	HashPassword "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Hash"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Utils/Tokens"
	_ "github.com/lib/pq"
)

type Postgress struct {
	DB *sql.DB
}

func InitiateDbConnection(cfg *Configurator.Configuration) (*Postgress, error) {
	db, err := sql.Open("postgres", cfg.StoragePath)
	if err != nil {
		slog.Info("The db connection was unsuccesfull")
		return nil, err
	}
	pingerror := db.Ping()
	if pingerror != nil {
		slog.Info("The ping to Db failed")
		return nil, pingerror
	}
	return &Postgress{DB: db}, nil
}

func (p *Postgress) SysAdminLogin(ctx context.Context, admin Types.SysAdminLogin) (sessionToken string, csrfToken string, err error) {
	var hashedPassword string
	var role string
	var email string

	queryerr := p.DB.QueryRowContext(ctx, "SELECT hashed_password,role,email from credentials WHERE email = $1", admin.Email).Scan(&hashedPassword, &role, &email)
	if queryerr != nil {
		slog.Info("There was an error in querying hashed password from db", "error", queryerr)
		return "", "", queryerr
	}

	passwordmatch, matcherr := HashPassword.Unhashpassword(admin.Password, hashedPassword)

	if matcherr != nil {
		slog.Info("There was an internal error", "error", "Hashin algorithim error")
		return "", "", errors.New("authentication Error")
	}
	if !passwordmatch {
		slog.Info("There was an auth error", "error", "Password/Email is wrong")
		return "", "", errors.New("authentication Error")
	}
	session_token, stokenerr := Tokens.GenerateToken(10)
	if stokenerr != nil {
		slog.Info("There was a session token generation error", "error", stokenerr)
		return "", "", stokenerr
	}
	csrf_token, csrftokenerr := Tokens.GenerateToken(10)
	if csrftokenerr != nil {
		slog.Info("There was a csrf token generation error", "error", stokenerr)
		return "", "", csrftokenerr
	}

	_, insertqerr := p.DB.ExecContext(ctx, "INSERT INTO sessions (session_token, csrf_token, email, role)  VALUES ($1, $2, $3, $4)", session_token, csrf_token, email, role)
	if insertqerr != nil {
		slog.Info("There was an error inserting data to db", "error", insertqerr)
		return "", "", nil
	}

	return session_token, csrf_token, nil
}

func (p *Postgress) SysAdminSignup(ctx context.Context, admin Types.SysAdminSignup) (err error) {
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

func (p *Postgress) AuthorizeSysAdmin(ctx context.Context, sessionToken string, csrfToken string) bool {
	var csrf string
	qerr := p.DB.QueryRowContext(ctx, "SELECT csrf_token FROM sessions WHERE session_token = $1", sessionToken).Scan(&csrf)
	if qerr != nil {
		if errors.Is(qerr, sql.ErrNoRows) {
			return false
		}
		return false
	}
	return csrf == csrfToken
}
