package http

import (
	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/handlers"
	"torque-dms/adapters/input/http/middleware"
	identityInput "torque-dms/core/identity/ports/input"
	inventoryInput "torque-dms/core/inventory/ports/input"
)

type Router struct {
	engine            *gin.Engine
	authService       identityInput.AuthService
	entityService     identityInput.EntityService
	permissionService identityInput.PermissionService
	vehicleService    inventoryInput.VehicleService
	locationService   inventoryInput.LocationService
}

func NewRouter(
	authService identityInput.AuthService,
	entityService identityInput.EntityService,
	permissionService identityInput.PermissionService,
	vehicleService inventoryInput.VehicleService,
	locationService inventoryInput.LocationService,
	jwtSecret string,
) *Router {
	r := &Router{
		engine:            gin.Default(),
		authService:       authService,
		entityService:     entityService,
		permissionService: permissionService,
		vehicleService:    vehicleService,
		locationService:   locationService,
	}

	r.setupRoutes(jwtSecret)
	return r
}

func (r *Router) setupRoutes(jwtSecret string) {
	// Handlers
	authHandler := handlers.NewAuthHandler(r.authService)
	entityHandler := handlers.NewEntityHandler(r.entityService)
	vehicleHandler := handlers.NewVehicleHandler(r.vehicleService)
	locationHandler := handlers.NewLocationHandler(r.locationService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)
	permissionMiddleware := middleware.NewPermissionMiddleware(r.permissionService)

	// Global middleware
	r.engine.Use(middleware.CORS())

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	public := r.engine.Group("/api")
	{
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
	}

	// Protected routes
	protected := r.engine.Group("/api")
	protected.Use(authMiddleware.Authenticate())
	protected.Use(permissionMiddleware.Check())
	{
		// Auth
		protected.POST("/auth/logout", authHandler.Logout)
		protected.POST("/auth/change-password", authHandler.ChangePassword)

		// Entities
		protected.GET("/entities", entityHandler.List)
		protected.GET("/entities/:id", entityHandler.GetByID)
		protected.POST("/entities", entityHandler.Create)
		protected.PUT("/entities/:id", entityHandler.Update)
		protected.DELETE("/entities/:id", entityHandler.Delete)
		protected.POST("/entities/:id/suspend", entityHandler.Suspend)
		protected.POST("/entities/:id/activate", entityHandler.Activate)

		// Locations
		protected.GET("/locations", locationHandler.List)
		protected.GET("/locations/active", locationHandler.ListActive)
		protected.GET("/locations/:id", locationHandler.GetByID)
		protected.POST("/locations", locationHandler.Create)
		protected.PUT("/locations/:id", locationHandler.Update)
		protected.DELETE("/locations/:id", locationHandler.Delete)
		protected.POST("/locations/:id/deactivate", locationHandler.Deactivate)
		protected.POST("/locations/:id/activate", locationHandler.Activate)

		// Vehicles
		protected.GET("/vehicles", vehicleHandler.List)
		protected.GET("/vehicles/available", vehicleHandler.ListAvailable)
		protected.GET("/vehicles/:id", vehicleHandler.GetByID)
		protected.GET("/vehicles/vin/:vin", vehicleHandler.GetByVIN)
		protected.POST("/vehicles", vehicleHandler.Create)
		protected.PUT("/vehicles/:id", vehicleHandler.Update)
		protected.DELETE("/vehicles/:id", vehicleHandler.Delete)
		protected.POST("/vehicles/:id/sold", vehicleHandler.MarkAsSold)
		protected.POST("/vehicles/:id/ready", vehicleHandler.MarkAsReadyForSale)
		protected.POST("/vehicles/:id/recon", vehicleHandler.SendToRecon)
		protected.POST("/vehicles/:id/location", vehicleHandler.ChangeLocation)

		// Vehicle Photos
		protected.GET("/vehicles/:id/photos", vehicleHandler.GetPhotos)
		protected.POST("/vehicles/:id/photos", vehicleHandler.AddPhoto)
		protected.POST("/vehicles/:id/photos/primary", vehicleHandler.SetPrimaryPhoto)
		protected.DELETE("/vehicles/:id/photos/:photoId", vehicleHandler.DeletePhoto)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}