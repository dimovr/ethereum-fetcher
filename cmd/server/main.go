package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ethereum_fetcher/db"
	"ethereum_fetcher/internal/config"
	"ethereum_fetcher/internal/routes"
	"ethereum_fetcher/internal/services"
)

func main() {
	Run()
}

func Run() {
	cfg := config.Load()
	db := initDb(cfg.DBConnectionURL)

	services, err := services.Init(db, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize services:  %v", err)
	}

	r := gin.Default()
	routes.SetupRoutes(services, r)

	if err := r.Run(":" + cfg.APIPort); err != nil {
		log.Fatalf("Failed to start server:  %v", err)
	}
}

func initDb(dbUrl string) *gorm.DB {
	dbConn, err := db.InitDB(dbUrl)
	if err != nil {
		log.Fatalf("Database connection failed:  %v", err)
	}
	return dbConn
}
