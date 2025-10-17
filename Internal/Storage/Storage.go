package Storage

import (
	"context"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
)

type SysAdmin interface {
	SysAdminLogin(ctx context.Context, admin Types.SysAdminLogin) (sessionToken string, csrfToken string, sysadminauth *Types.SysAdminAuthenticated, err error)
	SysAdminSignup(ctx context.Context, admin Types.SysAdminSignup) (err error)
	AuthorizeSysAdmin(ctx context.Context, sessionToken string, csrfToken string) (string, bool)
	GenerateInvite(ctx context.Context, adminid Types.SysAdminId, inviteData Types.SchoolInvite) (*Types.LinkGenerated, error)
	GetInviteData(ctx context.Context, token string) (string, error)
	GetInvitesAnalytics(ctx context.Context) (Types.InvitesAnalytics, error)
	GetInvites(ctx context.Context, limit int, offlimit int) ([]Types.SchoolInformation, error)
	RespondToSchoolInvite(ctx context.Context, token string, status string) (email string, err error)
	SaveSchoolAdminCredentials(ctx context.Context, email string, password string, token string) error
}

type SchoolAdmin interface {
	ValidateLink(ctx context.Context, inviteToken string) (string, error)
	SubmitInvite(ctx context.Context, schoolInfo Types.SchoolInformation, token string) (string, error)
}
