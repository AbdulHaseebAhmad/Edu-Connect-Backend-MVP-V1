package Types

type SysAdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SysAdminSignup struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Id       string `json:"id" validate:"required,min=6,alphanum"`
}

type SysAdminAuthenticated struct {
	Role          string `json:"role" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Status        bool   `json:"status" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
}
