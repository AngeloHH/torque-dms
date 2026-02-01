package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"torque-dms/adapters/input/http"
	"torque-dms/adapters/output/postgres/repositories"
	"torque-dms/core/identity/services"
	sharedDomain "torque-dms/core/shared/domain"
	"torque-dms/models"
)

func main() {
	// Cargar configuración desde .env
	dbHost := getEnv("DB_HOST", "db")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	webPort := getEnv("WEB_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-key-change-in-production")

	// Construir DATABASE_URL
	databaseURL := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	// Cargar reglas de validación
	if err := sharedDomain.LoadValidationRules("settings/validation_rules.yml"); err != nil {
		log.Fatal("Failed to load validation rules:", err)
	}

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to database")

	// Auto-migrar modelos
	if err := autoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated")

	// Crear repositories
	entityRepo := repositories.NewEntityRepository(db)
	userRepo := repositories.NewUserRepository(db)
	phoneRepo := repositories.NewPhoneRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	resourceRepo := repositories.NewResourceRepository(db)

	// Crear services
	entityService := services.NewEntityService(entityRepo, phoneRepo)
	authService := services.NewAuthService(entityRepo, userRepo, phoneRepo, jwtSecret)
	permissionService := services.NewPermissionService(roleRepo, resourceRepo)

	// Crear router
	router := http.NewRouter(authService, entityService, permissionService, jwtSecret)

	// Iniciar servidor
	log.Printf("Server starting on port %s", webPort)
	if err := router.Run(":" + webPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Geo
		&models.Country{},
		&models.Location{},
		&models.Route{},

		// Users
		&models.Entity{},
		&models.UserAccount{},
		&models.EntityPhone{},
		&models.Resource{},
		&models.Role{},
		&models.RoleResource{},
		&models.EntityResource{},
		&models.EntityRole{},

		// Inventory
		&models.VehicleModel3D{},
		&models.VehicleModelZone{},
		&models.Vehicle{},
		&models.VehicleLocationHistory{},
		&models.VehicleTracking{},
		&models.VehiclePhoto{},
		&models.VehicleZoneMark{},

		// Leads
		&models.LeadSource{},
		&models.LeadStepPreset{},
		&models.LeadStep{},
		&models.Lead{},
		&models.LeadStepProgress{},
		&models.LeadAssignment{},
		&models.LeadNote{},
		&models.LeadActivity{},
	)
}