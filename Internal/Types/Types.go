package Types

type SysAdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SchoolAdminLogin struct {
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
	Id            string `json:"id" validate:"required"`
}

type SchoolAdminAuthenticated struct {
	Role          string `json:"role" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Status        bool   `json:"status" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Id            string `json:"id" validate:"required"`
}

type SysAdminKey string
type SysAdminId string

type SchoolInvite struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type LinkGenerated struct {
	Messsage string `json:"message" validate:"required"`
	Token    string `json:"token" validate:"required"`
	SchoolInvite
}

type SchoolInformation struct {
	School     string `json:"schoolName" validate:"required"`
	Admin      string `json:"adminName" validate:"required"`
	Email      string `json:"contactEmail" validate:"required,email"`
	Phone      string `json:"contactPhone" validate:"required"`
	Country    string `json:"country" validate:"required"`
	Id         string `json:"schoolID" validate:"required"`
	Curriculam string `json:"curriculum" validate:"required"`
	Branch     string `json:"branch" validate:"required"`
	City       string `json:"city" validate:"required"`
	Priority   string `json:"priority" `
	Token      string `json:"token"`
	Sys_Eamil  string `json:"sys_email"`
	Username   string `json:"username"`
}

type InvitesAnalytics struct {
	Total          int     `json:"total" validate:"reuired"`
	Pending        int     `json:"pending" validate:"reuired"`
	Completed      int     `json:"completed" validate:"reuired"`
	Approved       int     `json:"approved" validate:"reuired"`
	ApprovalRate   float64 `json:"approvalRate" validate:"reuired"`
	AcceptanceRate float64 `json:"acceptanceRate" validate:"reuired"`
}
