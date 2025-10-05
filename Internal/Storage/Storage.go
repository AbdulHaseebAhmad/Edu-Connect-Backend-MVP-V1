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
}

type SchoolAdmin interface {
	ValidateLink(ctx context.Context, inviteToken string) (string, error)
}
