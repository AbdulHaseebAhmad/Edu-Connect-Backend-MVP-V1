package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	Configurator "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Config"
	Smtp "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Email/SMTP"
	SchoolAdminHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/SchoolAdministration"
	StudentAppHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/StudentsApp"
	SysAdminHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/SystemAdministration"
	UniversityHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/UniversityPortal"
	Middlewares "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Middleware"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	SchoolAdminStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/SchoolAdmin"
	StudentsAppStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Students"
	SysAdminStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/SysAdmins"
	UniversityStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/University"
	"github.com/rs/cors"
)

func main() {

	// setup cnfigurations
	cfg := Configurator.LoadConfiguration()
	fmt.Println(*cfg)

	smtp := Smtp.NewSMTPSender(cfg)
	//setup database

	db, dberror := Postgress.InitiateDbConnection(cfg)
	if dberror != nil {
		log.Fatal(dberror)
		return
	}
	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("Path", cfg.StoragePath))

	// Wrap the DB connection in your role-based store that implements Storage.
	studentStore := StudentsAppStorage.NewStudentsAppStore(db)

	// Wrap the DB connection in your role-based store that implements Storage.
	sysAdminStore := SysAdminStorage.NewSysAdminStore(db)

	// Wrap the DB connection in your role-based store that implements Storage.
	schoolAdminStore := SchoolAdminStorage.NewSchoolAdminStore(db)

	universityStore := UniversityStorage.NewUniversityStore(db)

	//setup router
	router := http.NewServeMux()

	//---->  Routes Starrt <-----
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome To the App"))
	})

	//University Routes

	router.HandleFunc("POST /api/university/login", UniversityHandler.Login(universityStore))

	router.Handle("GET /api/university/app/applications", Middlewares.Authorizer(sysAdminStore, UniversityHandler.GetStudntsApplications(universityStore)))
	router.Handle("GET /api/university/app/applications/respond", Middlewares.Authorizer(sysAdminStore, UniversityHandler.RespondStudntsApplications(universityStore)))
	router.Handle("GET /api/university/app/programs/list", Middlewares.Authorizer(sysAdminStore, UniversityHandler.GetUniversityProgramsList(universityStore)))
	router.Handle("GET /api/university/app/program/details", Middlewares.Authorizer(sysAdminStore, UniversityHandler.GetProgramDetails(universityStore)))
	router.Handle("POST /api/university/app/program/add", Middlewares.Authorizer(sysAdminStore, UniversityHandler.AddNewProgram(universityStore)))
	router.Handle("POST /api/university/app/program/update", Middlewares.Authorizer(sysAdminStore, UniversityHandler.UpdateProgram(universityStore)))

	router.Handle("GET /api/university/app/profile/get", Middlewares.Authorizer(sysAdminStore, UniversityHandler.GetUniversityProfile(universityStore)))
	router.Handle("POST /api/university/app/media/upload", Middlewares.Authorizer(sysAdminStore, UniversityHandler.UploadCampusMedia(universityStore)))
	//---> School Admin Routes Start <----
	router.HandleFunc("POST /api/schooladmin/login", SchoolAdminHandler.Login(schoolAdminStore))
	router.HandleFunc("GET /api/schooladmin/invite/validate", SchoolAdminHandler.LinkValidation(schoolAdminStore))
	router.HandleFunc("POST /api/schooladmin/invite/{invitation_id}/accept", SchoolAdminHandler.SubmitInviteData(schoolAdminStore))

	router.HandleFunc("GET /api/schooladmin/unprocessed/students", Middlewares.Authorizer(sysAdminStore, SchoolAdminHandler.GetUnProcessedStudentsList(schoolAdminStore)))
	router.HandleFunc("GET /api/schooladmin/verify/students", Middlewares.Authorizer(sysAdminStore, SchoolAdminHandler.VerifyStudentAccount(schoolAdminStore)))
	router.HandleFunc("GET /api/schooladmin/processed/students", Middlewares.Authorizer(sysAdminStore, SchoolAdminHandler.GetProcessedStudentsList(schoolAdminStore)))
	router.HandleFunc("GET /api/schooladmin/profile", Middlewares.Authorizer(sysAdminStore, SchoolAdminHandler.GetSchoolProfileData(schoolAdminStore)))
	//----> School Admin Routes Ennd <----

	//->> Sys Admin Auth Routes Start
	router.HandleFunc("POST /api/sysadmin/login", SysAdminHandler.Login(sysAdminStore))
	router.HandleFunc("POST /api/sysadmin/signup", SysAdminHandler.Signup(sysAdminStore))
	// router.HandleFunc("POST /api/sysadmin/invite/create", SysAdminHandler.CreateInvite(db))
	// ->> Sys Admin Auth Routes End

	// ->> Protected Routes Start <--

	// router.Handle("GET /api/sysadmin/testing", Middlewares.Authorizer(sysAdminStore, func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome To the Protected Route"))
	// }))

	// --> Sys Admin  Protected Routes
	router.Handle("POST /api/sysadmin/invite/create", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.CreateInvite(sysAdminStore)))
	router.Handle("POST /api/sysadmin/invite/send/{invitation_id}", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.SendInvite(sysAdminStore, smtp)))
	router.Handle("GET /api/sysadmin/invite/analytics", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetInvitesAnalytics(sysAdminStore)))
	router.Handle("GET /api/sysadmin/invite/applications", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetSchoolApplications(sysAdminStore)))
	router.Handle("GET /api/sysadmin/invite/application", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetSchoolApplication(sysAdminStore)))
	router.Handle("GET /api/sysadmin/invite/respond", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.RespondToSchoolApplication(sysAdminStore, smtp)))
	router.Handle("GET /api/sysadmin/invite/lists", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetAnalyticalLists(sysAdminStore)))

	router.Handle("GET /api/sysadmin/schools/get/invites", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetAllInvites(sysAdminStore)))
	router.Handle("GET /api/sysadmin/student/get/registered", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetRegisteredStudents(sysAdminStore)))

	router.Handle("GET /api/sysadmin/student/app/registry", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetStudentsRegistry(sysAdminStore)))
	router.Handle("GET /api/sysadmin/student/app/respond", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.RespondApplication(sysAdminStore, smtp)))
	router.Handle("GET /api/sysadmin/student/documents", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetStudentDocuments(sysAdminStore)))

	router.Handle("GET /api/sysadmin/get/receipts", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetReceipts(sysAdminStore)))
	router.Handle("GET /api/sysadmin/update/receipt/status", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.RespondToReceipts(sysAdminStore)))
	router.Handle("GET /api/sysadmin/scholarships/get", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetScholarships(sysAdminStore)))
	router.Handle("POST /api/sysadmin/scholarships/add", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.AddScholarships(sysAdminStore)))
	router.Handle("PUT /api/sysadmin/scholarships/update", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.UpdateScholarship(sysAdminStore)))
	router.Handle("DELETE /api/sysadmin/scholarships/delete", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.DeleteScholarship(sysAdminStore)))

	router.Handle("POST /api/sysadmin/webinar/create", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.CreateWebinar(sysAdminStore)))
	router.Handle("GET /api/sysadmin/webinar/get", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetWebinars(sysAdminStore)))
	router.Handle("PUT /api/sysadmin/webinar/update", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.UpdateWebinar(sysAdminStore)))
	router.Handle("DELETE /api/sysadmin/webinar/delete", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.DeleteWebinar(sysAdminStore)))

	// router.Handle("GET /api/sysadmin/review/applications/university", Middlewares.Authorizer(sysAdminStore,SysAdminHandler.))

	// ->> Protected Routes End <--

	// --> student App Routes

	router.Handle("POST /api/students/app/signup", StudentAppHandler.StudentSignup(studentStore))
	router.Handle("POST /api/students/app/login", StudentAppHandler.StudentSignin(studentStore))
	// router.Handle("GET /api/students/app/countries", StudentAppHandler.GetCountriesList(studentStore))

	// --> Student app Protected Routes

	router.Handle("GET /api/students/app/countries", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetCountriesList(studentStore)))
	router.Handle("GET /api/students/app/universities", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetUniversitiesList(studentStore)))
	router.Handle("GET /api/students/app/university", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetUniversityProfile(studentStore)))
	router.Handle("GET /api/students/app/programs", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetUniversityPrograms(studentStore)))
	router.Handle("GET /api/students/app/profile", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetStudentProfileDetails(studentStore)))
	router.Handle("GET /api/students/app/profile/update", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.UpdateStudentProfileDetails(studentStore)))
	router.Handle("GET /api/students/app/documents", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetstudentsDocuments(studentStore)))
	router.Handle("POST /api/students/app/documents/upload", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.UploadStudentDocuments(studentStore)))
	router.Handle("GET /api/students/app/documents/get", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetStudentsDocument(studentStore)))
	router.Handle("POST /api/students/app/receipt/upload", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.UploadApplicationReceipt(studentStore)))

	router.Handle("GET /api/students/app/university/apply", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.ApplyToUniversity(studentStore)))
	router.Handle("GET /api/students/app/university/applications", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetApplicationsData(studentStore)))
	router.Handle("GET /api/students/app/application/check", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.VerifyApplication(studentStore)))

	router.Handle("GET /api/students/app/programs/shortlist/add", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.ShortListProgram(studentStore)))
	router.Handle("GET /api/students/app/programs/shortlist/list", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.GetShortListProgram(studentStore)))
	router.Handle("GET /api/students/app/programs/shortlist/delete", Middlewares.Authorizer(sysAdminStore, StudentAppHandler.DeleteShortListProgram(studentStore)))

	//---->   Routes End   <-----

	// configure CORS options
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:5173/", "https://geosedutest.web.app", "https://geosedutest.web.app/", "https://pigeos.com",
			"https://www.pigeos.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	// wrap the router with CORS handler
	handler := c.Handler(router)
	//setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: handler,
	}
	fmt.Println("Server starting at:", cfg.HttpServer.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed To Start Server", err)
		} else {
		}
	}()

	<-done
	slog.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Info("failed to shhutdown", slog.String("erro", err.Error()))
	}
	slog.Info("Server Shut Down Successfully")
}
