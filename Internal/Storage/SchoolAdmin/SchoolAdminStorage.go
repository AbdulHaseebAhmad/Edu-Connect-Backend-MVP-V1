package SchoolAdminStorage

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
)

type SchoolAdminStore struct {
	*Postgress.Postgress
}

func NewSchoolAdminStore(pg *Postgress.Postgress) *SchoolAdminStore {
	return &SchoolAdminStore{pg}
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
