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
	SysAdminHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/SystemAdministration"
	Middlewares "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Middleware"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
	SchoolAdminStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/SchoolAdmin"
	SysAdminStorage "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/SysAdmins"
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
	sysAdminStore := SysAdminStorage.NewSysAdminStore(db)

	// Wrap the DB connection in your role-based store that implements Storage.
	schoolAdminStore := SchoolAdminStorage.NewSchoolAdminStore(db)

	//setup router
	router := http.NewServeMux()

	//---->  Routes Starrt <-----
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome To the App"))
	})

	//->> Sys Admin Auth Routes Start
	router.HandleFunc("POST /api/sysadmin/login", SysAdminHandler.Login(sysAdminStore))
	router.HandleFunc("POST /api/sysadmin/signup", SysAdminHandler.Signup(sysAdminStore))
	// router.HandleFunc("POST /api/sysadmin/invite/create", SysAdminHandler.CreateInvite(db))
	// ->> Sys Admin Auth Routes End

	//---> School Admin Routes Start <---
	router.HandleFunc("GET /api/schooladmin/invite/validate", SchoolAdminHandler.LinkValidation(schoolAdminStore))
	router.HandleFunc("POST /api/schooladmin/invite/{token}/accept", SchoolAdminHandler.SubmitInviteData(schoolAdminStore))
	//----> School Admin Routes Ennd <----

	// ->> Protected Routes Start <--

	// router.Handle("GET /api/sysadmin/testing", Middlewares.Authorizer(sysAdminStore, func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Welcome To the Protected Route"))
	// }))

	// --> Sys Admin  Protected Routes
	router.Handle("POST /api/sysadmin/invite/create", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.CreateInvite(sysAdminStore)))
	router.Handle("POST /api/sysadmin/invite/send/{token}", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.SendInvite(sysAdminStore, smtp)))
	router.Handle("GET /api/sysadmin/invite/analytics", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetInvitesAnalytics(sysAdminStore)))
	router.Handle("GET /api/sysadmin/invite/applications", Middlewares.Authorizer(sysAdminStore, SysAdminHandler.GetInvitesApplications(sysAdminStore)))
	// ->> Protected Routes End <--

	//---->   Routes End   <-----

	// configure CORS options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5173/"},
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
