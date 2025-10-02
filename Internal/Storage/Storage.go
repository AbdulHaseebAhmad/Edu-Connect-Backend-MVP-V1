package Storage

import (
	"context"

	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Types"
)

type Storage interface {
	SysAdminLogin(ctx context.Context, admin Types.SysAdminLogin) (sessionToken string, csrfToken string, sysadminauth *Types.SysAdminAuthenticated, err error)
	SysAdminSignup(ctx context.Context, admin Types.SysAdminSignup) (err error)
	AuthorizeSysAdmin(ctx context.Context, sessionToken string, csrfToken string) (isAuthorize bool)
}
