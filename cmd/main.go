package main

import (
	"banner"
	"banner/pkg/handler"
	"banner/pkg/repository"
	"banner/pkg/service"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "db",
		Port:     "5432",
		Username: "postgres",
		Password: "admin",
		DBName:   "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	//migrations
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			"postgres", "admin", "db", "5432", "postgres"))
	if err != nil {
		logrus.Fatalf("Error creating migration instance: %v", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error applying migrations: %v", err)
	}
	logrus.Debug("Migrations applied successfully")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(banner.Server)
	go func() {
		if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Printf("Banner App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Banner App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
