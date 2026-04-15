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
	GetInviteData(ctx context.Context, invitation_id string) (string, error)
	GetInvitesAnalytics(ctx context.Context) (Types.InvitesAnalytics, error)
	GetSchoolApplications(ctx context.Context, limit int, offlimit int) ([]Types.SchoolApp, error)
	GetSchoolApplication(ctx context.Context, application_id string) (Types.SchoolApplication, error)
	RespondToSchoolInvite(ctx context.Context, application_id string, status string) (schoolInfo Types.SchoolInformation, generatePassword string, err error)
	GetAnalyticsList(ctx context.Context) (lists Types.AnalyticsList, err error)
	GetStudentsRegistry(ctx context.Context, status string) ([]Types.StudentsRegistry, error)
	GetRegisteredStudents(ctx context.Context) ([]Types.RegisteredStudent, error)
	RespondApplication(ctx context.Context, action string, id string) (string, string, string, error)
	GetStudentsDocument(ctx context.Context, studentId string, documentname string, documentmime string) (Types.Documents, error)
	GetAllReceipts(ctx context.Context) ([]Types.UniversityAppReceipt, error)
	GetReceiptDetails(ctx context.Context, student_id string) (Types.UniversityAppReceipt, error)
	RespondToReceipts(ctx context.Context, receipt_id string, status string) error

	GetAllInvites(ctx context.Context) (invitesList []Types.Invite, err error)
	GetScholarships(ctx context.Context) (scholarShips []Types.Scholarship, err error)
	AddScholarship(ctx context.Context, scholarship Types.Scholarship) error
	UpdateScholarship(ctx context.Context, scholarship Types.Scholarship, scholarshipId string) error
	DeleteScholarship(ctx context.Context, scholarshipId string) error

	CreateWebinar(ctx context.Context, webinar Types.Webinar) (string, error)
	GetWebinars(ctx context.Context) (webinars []Types.Webinar, err error)
	UpdateWebinar(ctx context.Context, webinarId string, webinar Types.Webinar) error
	DeleteWebinar(ctx context.Context, webinarId string) error

	GetUniversities(ctx context.Context) (listOfUniversities []Types.FeaturedUniversity, err error)

	AddFeaturedPartners(ctx context.Context, partners []Types.FeaturedPartner) (err error)
	GetFeaturedPartners(ctx context.Context) (listOfPartners []Types.FeaturedUniversity, err error)
	DeleteFeaturedPartner(ctx context.Context, partner_id string) error

	GetUniversitiesCommissions(ctx context.Context) (commisions []Types.Commision, err error)
}

type SchoolAdmin interface {
	ValidateLink(ctx context.Context, invitation_id string) (string, error)
	SubmitInvite(ctx context.Context, schoolInfo Types.SchoolInformation) (string, error)
	SchoolAdminLogin(ctx context.Context, schooladmin Types.SchoolAdminLogin) (sessionToken string, csrfToken string, sysadminauth *Types.SchoolAdminAuthenticated, err error)
	GetUnProcessedStudentsList(ctx context.Context, school_id string) ([]Types.StudentProfile, error)
	VerifyStudentAccount(ctx context.Context, school_id string, student_id string, status string) error
	GetProcessedStudentsList(ctx context.Context, school_id string, status string) ([]Types.StudentProfile, error)
	GetSchoolProfileData(ctx context.Context, school_id string) (Types.SchoolInformation, error)
}

type StudentsApp interface {
	StudentSignup(ctx context.Context, studentSignup Types.StudentsSignupData, studentsDocuments Types.StudentSignupDocuments) error
	StudentsLogin(ctx context.Context, studentLogin Types.StudentsSignIn) (sessionToken string, csrfToken string, studentAuth *Types.StudentAuthenticated, err error)
	GetCountryList(ctx context.Context, countrycode string) ([]Types.Countries, error)
	GetUniversitiesList(ctx context.Context, countrycode string) ([]Types.University, error)
	GetUniversityProfile(ctx context.Context, university_id string) (Types.UniversityProfile, error)
	GetUniversityPrograms(ctx context.Context, program_id string) (programsList []Types.Programe, err error)
	GetStudentProfileDetails(ctx context.Context, studentId string) (studentprofile Types.StudentProfile, err error)
	UpdateProfile(ctx context.Context, student_id string, fieldName string, fieldValue string) error
	UpdateContact(ctx context.Context, student_id string, fieldName string, fieldValue string) error
	UpdateEducation(ctx context.Context, student_id string, fieldName string, fieldValue string) error
	UpdatePreferences(ctx context.Context, student_id string, fieldName string, fieldValue string) error
	GetstudentsDocuments(ctx context.Context, student_id string) ([]Types.UploadDocument, error)
	UploadStudentDocuments(ctx context.Context, studentDocs Types.UploadDocument, documentBytes []byte, student_id string) error
	GetStudentsDocument(ctx context.Context, studentId string, document_id string) (Types.Documents, error)
	UploadApplicationReceipt(ctx context.Context, student_id string, university_id string, program_id string, paid_amount string, receipt Types.UploadDocument) error
	ApplyToUniversity(ctx context.Context, student_id string, university_id string, program_id string) error
	GetApplicationsData(ctx context.Context, student_id string) ([]Types.ApplicationData, error)
	VerifyApplication(ctx context.Context, student_id string, program_id string, university_id string) (Types.ExistsRow, error)

	ShortListProgram(ctx context.Context, student_id string, program_id string, university_id string) (int, error)
	GetShortListProgram(ctx context.Context, student_id string) (shortListedPrograms []Types.ShortListProgram, err error)
	DeleteShortListProgram(ctx context.Context, student_id string, shortListId string) error

	RegisterForEvent(ctx context.Context, student_id string, webinar_code string) (email string, webinarData Types.Webinar, err error)
	EventRegisterationCheck(ctx context.Context, student_id string, webinar_code string) (bool, error)
}

type UniversityPortal interface {
	UniversityLogin(ctx context.Context, universityLogin Types.UniversityLogin) (sessionToken string, csrfToken string, universityAuth *Types.UniversityAuthenticated, err error)
	GetStudntsApplications(ctx context.Context, university_id string, status string) ([]Types.StudentProfile, error)
	RespondStudntsApplications(ctx context.Context, application_id string, status string) error
	GetUniversityProgramsList(ctx context.Context, university_id string) ([]Types.ProgrameThumbnail, error)
	GetProgramDetails(ctx context.Context, program_id string) (Types.Programe, error)
	AddNewProgram(ctx context.Context, programDetails Types.Programe) error
	UpdateProgram(ctx context.Context, programDetails Types.Programe, program_id string) error
	GetUniversityProfile(ctx context.Context, university_id string) (Types.UniversityProfile, error)
	UploadCampusMedia(ctx context.Context, university_id string, uploadedMidaArray []Types.UniMedia) error
	GetAllStudntsApplications(ctx context.Context, university_id string, status string) ([]Types.StudentProfile, error)
}
