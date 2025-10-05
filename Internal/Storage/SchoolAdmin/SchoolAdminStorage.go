package SchoolAdminStorage

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
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

	parsedTime, err := time.Parse(time.RFC3339Nano, returneTimeStamp)

	if err != nil {
		slog.Info("There was an error queryying", "error", err)
		return "", err
	}

	expired := time.Now().Before(parsedTime)

	if !expired {
		return "", errors.New("token is expired")
	}
	return status, nil
}
