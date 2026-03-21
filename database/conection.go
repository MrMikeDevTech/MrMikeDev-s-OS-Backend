package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MrMikeDevTech/mrmikedevs-os/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "15801")
	user := getEnv("DB_USER", "mrmikedev")
	pass := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "mrmikedevsos")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, pass, name, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("No se pudo conectar a la DB: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error al obtener la instancia de sql.DB:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Conexión a PostgreSQL establecida en puerto", port)

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("❌ Error en AutoMigrate:", err)
	}
	log.Println("🚀 Migración completada")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
