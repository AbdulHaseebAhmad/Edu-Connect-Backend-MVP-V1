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
	SysAdminHandler "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Handlers/SystemAdministration"
	"github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Storage/Postgress"
)

func main() {

	// setup cnfigurations
	cfg := Configurator.LoadConfiguration()
	fmt.Println(*cfg)

	//setup database

	db, dberror := Postgress.InitiateDbConnection(cfg)
	if dberror != nil {
		log.Fatal(dberror)
		return
	}
	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("Path", cfg.StoragePath))

	//setup router
	router := http.NewServeMux()

	//---->  Routes Starrt <-----
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome To the App"))
	})

	//->> Sys Admin Auth Routes Start
	router.HandleFunc("POST /api/sysadmin/login", SysAdminHandler.Login(db))
	router.HandleFunc("POST /api/sysadmin/signup", SysAdminHandler.Signup(db))
	// ->> Sys Admin Auth Routes End
	//---->   Routes End   <-----
	//setup server
	server := http.Server{
		Addr:    cfg.HttpServer.Address,
		Handler: router,
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
