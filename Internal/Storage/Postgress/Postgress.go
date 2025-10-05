package Postgress

import (
	"database/sql"
	"log/slog"

	Configurator "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Config"
	_ "github.com/lib/pq"
)

type Postgress struct {
	DB *sql.DB
}

func InitiateDbConnection(cfg *Configurator.Configuration) (*Postgress, error) {
	db, err := sql.Open("postgres", cfg.StoragePath)
	if err != nil {
		slog.Info("The db connection was unsuccesfull")
		return nil, err
	}
	pingerror := db.Ping()
	if pingerror != nil {
		slog.Info("The ping to Db failed")
		return nil, pingerror
	}
	return &Postgress{DB: db}, nil
}
