package Types

import "github.com/lib/pq"

type SysAdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SchoolAdminLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UniversityLogin struct {
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
	CsrfToken     string `json:"csrf_token"`
}

type StudentAuthenticated struct {
	Role           string `json:"role" validate:"required"`
	Authenticated  bool   `json:"authenticated" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Status         string `json:"student_status" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Id             string `json:"student_id" validate:"required"`
	SchoolVerified string `json:"school_verified" db:"school_verified"`
	CsrfToken      string `json:"csrf_token" validate:"required"`
}

type UniversityAuthenticated struct {
	Role          string `json:"role" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Status        bool   `json:"status" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Id            string `json:"university_id" validate:"required"`
	CsrfToken     string `json:"csrf_token"`
}
type SchoolAdminAuthenticated struct {
	Role          string `json:"role" validate:"required"`
	Authenticated bool   `json:"authenticated" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Status        bool   `json:"status" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Id            string `json:"id" validate:"required"`
	CsrfToken     string `json:"csrf_token"`
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
	RegistrationCode *string `json:"registration_code" validate:"required"`
	School           string  `json:"school_name" validate:"required"`
	Email            string  `json:"school_email" validate:"required,email"`
	Admin            string  `json:"admin_name" validate:"required"`
	Phone            string  `json:"school_phone" validate:"required"`
	Country          string  `json:"school_country" validate:"required"`
	Curriculam       string  `json:"school_curriculum" validate:"required"`
	Branch           string  `json:"school_branch" validate:"required"`
	City             string  `json:"school_city" validate:"required"`
	SchoolId         string  `json:"school_id"`
	Status           string  `json:"status"`
	Sys_Eamil        string  `json:"sys_email"`
	Priority         string  `json:"priority" `
	Username         string  `json:"username"`
	CreatedAt        string  `json:"createdat"`
	CompletedAt      *string `json:"completedat"`
	ApprovedAt       *string `json:"approvedat"`
	Expiry           *string `json:"expiryat"`
	Code             string  `json:"school_code"`
}

type SchoolApplication struct {
	ApplicationId    string  `json:"application_id" validate:"required"`
	SchoolId         string  `json:"school_id" validate:"required"`
	RegistrationCode *string `json:"registration_code" validate:"required"`
	SchoolName       string  `json:"school_name" validate:"required"`
	SchoolEmail      string  `json:"school_email" validate:"required,email"`
	AdminName        string  `json:"admin_name" validate:"required"`
	SchoolPhone      string  `json:"school_phone" validate:"required"`
	SchoolCountry    string  `json:"school_country" validate:"required"`
	SchoolCurriculam string  `json:"school_curriculum" validate:"required"`
	SchoolBranch     string  `json:"school_branch" validate:"required"`
	SchoolCity       string  `json:"school_city" validate:"required"`
}

type InvitesAnalytics struct {
	Total          int     `json:"total" validate:"reuired"`
	Pending        int     `json:"pending" validate:"reuired"`
	Completed      int     `json:"completed" validate:"reuired"`
	Approved       int     `json:"approved" validate:"reuired"`
	ApprovalRate   float64 `json:"approvalRate" validate:"reuired"`
	AcceptanceRate float64 `json:"acceptanceRate" validate:"reuired"`
}

type AnalyticsList struct {
	PendingInvites       []SchoolInformation
	ApprovedApplications []SchoolInformation
	PendingApplications  []SchoolInformation
}

type StudentsSignupData struct {
	First_Name     string `json:"fname" validate:"required"`
	Last_Name      string `json:"lname" validate:"require"`
	Email          string `json:"email" validate:"required,email"`
	SchoolCode     string `json:"school_code" validate:"required"`
	Citizenship    string `json:"citizenship" validate:"required"`
	Passport       string `json:"passport" validate:"required"`
	PassportType   string `json:"passport_mime_type" validate:"required"`
	Transcript     string `json:"transcript" validate:"required"`
	TranscriptType string `json:"transcript_mime_type" validate:"required"`
}

type StudentSignupDocuments struct {
	Passport   []byte
	Transcript []byte
}

type StudentsRegistry struct {
	First_Name     string `json:"fname" validate:"required"`
	Last_Name      string `json:"lname" validate:"require"`
	Email          string `json:"email" validate:"required,email"`
	Citizenship    string `json:"citizenship" validate:"required"`
	Slug           string `json:"slug" validate:"required"`
	Status         string `json:"status" validate:"required"`
	SchoolVerified string `json:"school_verified"`
	Date           string `json:"created_at"`
}

type Documents struct {
	Data     []byte
	MimeType string
	Name     string
	Status   string
}

type UploadDocument struct {
	Id       string `json:"document_id"`
	Data     string `json:"data"`
	Name     string `json:"name"`
	Document string `json:"document_name"`
	MimeType string `json:"type"`
	Status   string `json:"status"`
}

type StudentsSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Countries struct {
	Name         string `json:"name"`
	Id           int    `json:"id"`
	CountryCode  string `json:"country_code"`
	IsoCode      string `json:"iso_code"`
	IsoCode3     string `json:"iso3_code"`
	DiallingCode string `json:"dialling_code"`
	Image        string `json:"image_url"`
	Programs     string `json:"program_count"`
}

type University struct {
	Id string `json:"university_id"`
	// Code    string `json:"university_code"`
	Name        string `json:"university_name"`
	City        string `json:"university_city"`
	Country     string `json:"university_country"`
	CountryName string `json:"country_code"`
	Acronym     string `json:"university_acronym"`
	Image       string `json:"university_image"`
	Phone       string `json:"university_phone"`
}

type UniversityProfile struct {
	UniversityName  string `json:"university_name"`
	UniversityId    string `json:"university_id" db:"university_id"`
	UniversityCity  string `json:"university_city" db:"university_city"`
	StudentsCount   string `json:"students_count" db:"students_count"`
	AcceptanceRate  string `json:"acceptance_rate" db:"acceptance_rate"`
	QSRanking       string `json:"qs_ranking" db:"qs_ranking"`
	AboutUniversity string `json:"about_university" db:"about_university"`
	FoundedDate     string `json:"founded_date" db:"founded_date"`
	Type            string `json:"type" db:"type"`
	Calendar        string `json:"calendar" db:"calendar"`
	GraduationRate  string `json:"graduation_rate" db:"graduation_rate"`
	Employability   string `json:"employability" db:"employability"`
	Programs        []UniProfileProgram
	Media           []UniMedia
	UniContact
}

type UniContact struct {
	UniversityAdmissionEmail string `json:"university_admission_email"`
	UniversityEmail          string `json:"university_email"`
	UniversityAddress        string `json:"university_address"`
	UniversityPhone          string `json:"university_phone"`
	UniversityWebsite        string `json:"university_website"`
	UniversityInstagram      string `json:"university_instagram"`
	UniversityYoutube        string `json:"university_youtube"`
	UniversityLinkedin       string `json:"university_linkedin"`
	UniversityX              string `json:"university_x"`
}

// type UnieMedia struct {
// 	Media         string `json:"media"`
// 	MediaType     string `json:"media_type"`
// 	MediaFileName string `json:"media_file_name"`
// 	// MediaSize     int    `json:"media_size"`
// 	MediaTag string `json:"media_tag"`
// }

type UniMedia struct {
	MediaId       int    `json:"media_id"`
	Media         string `json:"media"`
	MediaFileName string `json:"media_file_name"`
	MediaType     string `json:"media_type"`
	MediaSize     string `json:"media_size"`
	MediaTag      string `json:"media_tag"`
}

type UniProfileProgram struct {
	ProgramId       string `json:"program_id" db:"program_id"`
	ProgramName     string `json:"program_name" db:"program_name"`
	ProgramLevel    string `json:"program_level" db:"program_level"`
	ProgramDuration string `json:"program_duration" db:"program_duration"`
	SessionIntake   string `json:"session_intake" db:"session_intake"`
	ProgramFee      string `json:"program_fee" db:"program_fee"`
	ProgramStatus   string `json:"program_status" db:"program_status"`
	ApplicationFee  string `json:"program_application_fee" db:"application_fee"`
}
type Programe struct {
	Id string `json:"program_id" db:"program_id"`
	// PCode           string   `json:"program_code" db:"program_code"`
	Name            string   `json:"program_name"`
	PFee            string   `json:"program_fee"`
	Duration        string   `json:"program_duration"`
	SessionIntake   string   `json:"session_intake"`
	Description     string   `json:"program_description"`
	Requirements    []string `db:"program_requirements" json:"program_requirements"`
	RelatedTags     []string `db:"related_tags" json:"related_tags"`
	PossibleCareers []string `db:"possible_careers" json:"possible_careers"`
	RequiredDocs    []string `db:"program_required_documents" json:"program_required_documents"`
	AFee            string   `json:"program_application_fee"`
	UniversityCode  string   `json:"university_id"`
	UniversityName  string   `json:"university_name"`
	ProgramLevel    string   `json:"program_level"`
	ProgramCapacity string   `json:"program_capacity"`
}

type ProgrameThumbnail struct {
	Id   string `json:"program_id" db:"program_id"`
	Name string `json:"program_name"`
	// PCode         string `json:"program_code" db:"program_code"`
	SessionIntake string `json:"session_intake"`
	Status        string `json:"program_status"`
	Capacity      string `json:"program_capacity"`
}

type StudentPersonalDetails struct {
	First_Name      string `json:"first_name" db:"first_name"`
	Middle_Name     string `json:"middle_name" db:"middle_name"`
	Last_Name       string `json:"last_name" db:"last_name"`
	Dob             string `json:"dob" db:"dob"`
	Gender          string `json:"gender" db:"gender"`
	Marrital_Status string `json:"marrital_status" db:"marrital_status"`
}

type StudentCitizenshipDetails struct {
	Nationality     string `json:"nationality" db:"nationality"`
	Passport_Number string `json:"passport_number" db:"passport_number"`
	Passport_Expiry string `json:"passport_expiry" db:"passport_expiry"`
	Place_Of_Issue  string `json:"passport_issue" db:"passport_issue"`
}
type StudentProfile struct {
	StudentId      string `json:"student_id"`
	UniversityId   string `json:"university_id"`
	Decision       string `json:"decision_status"`
	Date           string `json:"created_at"`
	ApplicationId  string `json:"application_id"`
	Program        string `json:"program_name"`
	ProgramSession string `json:"session_intake"`
	StudentPersonalDetails
	StudentCitizenshipDetails
	StudentContact
	StudentEducation
	StudentPrefferences
	Documents []UploadDocument
}

type StudentContact struct {
	Email             string `json:"email" db:"email"`
	Phone             string `json:"phone_number" db:"phone_number"`
	WhatsApp          string `json:"whatsapp_number" db:"whatsapp_number"`
	Address           string `json:"permanent_address" db:"permanent_address"`
	Street            string `json:"street_address" db:"street_address"`
	City              string `json:"city" db:"city"`
	Province          string `json:"state_province" db:"state_province"`
	PostCode          string `json:"zip_postal_code" db:"zip_postal_code"`
	EmergencyPhone    string `json:"emergency_phone" db:"emergency_phone"`
	EmergencyName     string `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyRelation string `json:"emergency_relationship" db:"emergency_relationship"`
}

type StudentEducation struct {
	SchoolName       string `json:"school_name" db:"school_name"`
	Curriculum       string `json:"curriculum" db:"curriculum"`
	GraduationYear   string `json:"graduation_year" db:"graduation_year"`
	CummulativeScore string `json:"cummulative_score" db:"cummulative_score"`

	LanguageType         string `json:"language_type" db:"language_type"`
	LanguageOverallScore string `json:"language_overall_score" db:"language_overall_score"`
	LanguageReading      string `json:"language_reading" db:"language_reading"`
	LanguageWritting     string `json:"language_writting" db:"language_writting"`
	LanguageSpeaking     string `json:"language_speaking" db:"language_speaking"`
	LanguageListening    string `json:"language_listening" db:"language_listening"`
}

type StudentPrefferences struct {
	PrimaryInterest string `json:"primary_career_interest" db:"primary_career_interest"`
	Degree          string `json:"degree_level" db:"degree_level"`
	PrefferedDate   string `json:"preferred_start_date" db:"preferred_start_date"`
	AnnualBudget    string `json:"annual_budget" db:"annual_budget"`
	Schalarship     string `json:"scholarship_interest" db:"scholarship_interest"`
}

type StudentDocuments struct {
	Cv               Documents `json:"cv"`
	Passport         Documents `json:"passport"`
	Identity         Documents `json:"identity"`
	Language         Documents `json:"language_proefficiency"`
	Highschool       Documents `json:"high_school"`
	Coverletter      Documents `json:"cover_letter"`
	Motivationletter Documents `json:"motivation_letter"`
}

type UniversityAppReceipt struct {
	StudentID      string `json:"student_id" db:"student_id"`
	FirstName      string `json:"first_name" db:"first_name"`
	StudentEmail   string `json:"student_email" db:"student_email"`
	UniversityID   string `json:"university_id" db:"university_id"`
	ApplicationID  string `json:"application_id" db:"application_id"`
	ReceiptID      string `json:"receipt_id" db:"receipt_id"`
	Receipt        []byte `json:"receipt" db:"receipt"` // BYTEA / blob
	MimeType       string `json:"mime_type" db:"mime_type"`
	ReceiptName    string `json:"receipt_name" db:"receipt_name"`
	ReceiptStatus  string `json:"receipt_status" db:"receipt_status"`
	PaidAmount     string `json:"paid_amount" db:"paid_amount"`
	CreatedDate    string `json:"created_date" db:"created_date"` // DATE in PG → time.Time (date-only)
	UniversityName string `json:"university_name" db:"university_name"`
	ProgramId      string `json:"program_id"`
}

type Invite struct {
	InviteId        string `json:"invitation_id"`
	SchoolName      string `json:"school_name"`
	SchoolEmail     string `json:"school_email"`
	SchoolAdminName string `json:"admin_name"`
	InviteStatus    string `json:"status"`
	Date            string `json:"created_at"`
}

type SchoolApp struct {
	ApplicationId     string `json:"application_id"`
	SchoolName        string `json:"school_name"`
	SchoolEmail       string `json:"school_email"`
	SchoolAdminName   string `json:"admin_name"`
	ApplicationStatus string `json:"status"`
	Date              string `json:"created_at"`
	Priority          string `json:"priority"`
}

type RegisteredStudent struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	StudentId      string `json:"student_id"`
	SchoolVerified string `json:"school_verified"`
	Status         string `json:"status"`
	SchoolId       string `json:"school_id"`
	Email          string `json:"student_email"`
	CreatedDate    string `json:"created_at"`
}

type ApplicationData struct {
	StudentId         string `json:"student_id"`
	ApplicationId     string `json:"application_id"`
	UniversityId      string `json:"university_id"`
	ProgramId         string `json:"program_id"`
	UniversityName    string `json:"university_name"`
	ProgramName       string `json:"program_name"`
	CountryName       string `json:"country_name"`
	ApplicationDate   string `json:"application_date"`
	ApplicationStatus string `json:"application_status"`
	DecisionStatus    string `json:"decision_status"`
}

type ExistsRow struct {
	HasReceipt     bool `db:"has_receipt"`
	HasApplication bool `db:"has_application"`
}

type Scholarship struct {
	ScholarshipID string         `json:"scholarship_id" db:"scholarship_id"`
	Title         string         `json:"title" db:"title"`
	Country       string         `json:"country" db:"country"`
	Region        string         `json:"region" db:"region"`
	Level         string         `json:"level" db:"level"`
	Funding       string         `json:"funding" db:"funding"`
	Status        string         `json:"status" db:"status"`
	Opens         string         `json:"opens" db:"opens"`
	Deadline      string         `json:"deadline" db:"deadline"`
	Description   string         `json:"description" db:"description"`
	Link          string         `json:"link" db:"link"`
	Requirements  pq.StringArray `json:"requirements" db:"requirements"`
}

type ShortListProgram struct {
	ShortListId int    `json:"id"`
	Programe    string `json:"program_id"`
	University  string `json:"university_id"`
	Student     string `json:"student_id"`
}

type Webinar struct {
	ID           int64  `db:"id" json:"id"`
	WebinarCode  string `db:"webinar_code" json:"webinar_code"`
	Title        string `db:"title" json:"title"`
	Speaker      string `db:"speaker" json:"speaker"`
	Link         string `db:"link" json:"link"`
	Date         string `db:"date" json:"date"` // stored as string in Postgres
	Time         string `db:"time" json:"time"` // stored as string in Postgres
	Platform     string `db:"platform" json:"platform"`
	TargetType   string `db:"targettype" json:"targettype"`
	TargetValue  string `db:"targetvalue" json:"targetvalue"`
	Status       string `db:"status" json:"status"`
	Registration int    `db:"registered" json:"registered"`
}
