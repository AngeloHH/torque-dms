package http

import (
	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/handlers"
	"torque-dms/adapters/input/http/middleware"
	identityInput "torque-dms/core/identity/ports/input"
	inventoryInput "torque-dms/core/inventory/ports/input"
	salesInput "torque-dms/core/sales/ports/input"
)

type Router struct {
	engine            *gin.Engine
	authService       identityInput.AuthService
	entityService     identityInput.EntityService
	permissionService identityInput.PermissionService
	vehicleService    inventoryInput.VehicleService
	locationService   inventoryInput.LocationService
	leadService       salesInput.LeadService
	stepService       salesInput.StepService
}

func NewRouter(
	authService identityInput.AuthService,
	entityService identityInput.EntityService,
	permissionService identityInput.PermissionService,
	vehicleService inventoryInput.VehicleService,
	locationService inventoryInput.LocationService,
	leadService salesInput.LeadService,
	stepService salesInput.StepService,
	jwtSecret string,
) *Router {
	r := &Router{
		engine:            gin.Default(),
		authService:       authService,
		entityService:     entityService,
		permissionService: permissionService,
		vehicleService:    vehicleService,
		locationService:   locationService,
		leadService:       leadService,
		stepService:       stepService,
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
	leadHandler := handlers.NewLeadHandler(r.leadService, r.stepService)
	stepHandler := handlers.NewStepHandler(r.stepService)

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

		// Lead Sources
		protected.GET("/lead-sources", leadHandler.GetSources)
		protected.GET("/lead-sources/active", leadHandler.GetActiveSources)
		protected.POST("/lead-sources", leadHandler.CreateSource)
		protected.POST("/lead-sources/:id/deactivate", leadHandler.DeactivateSource)
		protected.POST("/lead-sources/:id/activate", leadHandler.ActivateSource)

		// Leads
		protected.GET("/leads", leadHandler.List)
		protected.GET("/leads/:id", leadHandler.GetByID)
		protected.POST("/leads", leadHandler.Create)
		protected.PUT("/leads/:id", leadHandler.Update)
		protected.DELETE("/leads/:id", leadHandler.Delete)

		// Lead Assignments
		protected.GET("/leads/:id/assignments", leadHandler.GetAssignments)
		protected.POST("/leads/:id/assignments", leadHandler.Assign)
		protected.POST("/leads/:id/assignments/primary", leadHandler.SetPrimaryAssignment)
		protected.DELETE("/leads/:id/assignments/:assignmentId", leadHandler.RemoveAssignment)

		// Lead Notes
		protected.GET("/leads/:id/notes", leadHandler.GetNotes)
		protected.POST("/leads/:id/notes", leadHandler.AddNote)
		protected.PUT("/leads/:id/notes/:noteId", leadHandler.UpdateNote)
		protected.DELETE("/leads/:id/notes/:noteId", leadHandler.DeleteNote)

		// Lead Activities
		protected.GET("/leads/:id/activities", leadHandler.GetActivities)
		protected.POST("/leads/:id/activities", leadHandler.AddActivity)
		protected.POST("/leads/:id/activities/:activityId/complete", leadHandler.CompleteActivity)
		protected.GET("/activities/scheduled", leadHandler.GetMyScheduledActivities)
		protected.GET("/activities/overdue", leadHandler.GetOverdueActivities)

		// Lead Progress
		protected.GET("/leads/:id/progress", leadHandler.GetProgress)
		protected.PUT("/leads/:id/progress/:stepId", leadHandler.UpdateProgress)

		// Step Presets
		protected.GET("/presets", stepHandler.GetPresets)
		protected.GET("/presets/public", stepHandler.GetPublicPresets)
		protected.GET("/presets/mine", stepHandler.GetMyPresets)
		protected.GET("/presets/:id", stepHandler.GetPreset)
		protected.POST("/presets", stepHandler.CreatePreset)
		protected.DELETE("/presets/:id", stepHandler.DeletePreset)
		protected.POST("/presets/:id/public", stepHandler.MakePresetPublic)
		protected.POST("/presets/:id/shared", stepHandler.MakePresetShared)
		protected.POST("/presets/:id/private", stepHandler.MakePresetPrivate)

		// Preset Steps
		protected.GET("/presets/:presetId/steps", stepHandler.GetSteps)
		protected.POST("/presets/:presetId/steps", stepHandler.CreateStep)
		protected.POST("/presets/:presetId/steps/:stepId/deactivate", stepHandler.DeactivateStep)
		protected.POST("/presets/:presetId/steps/:stepId/activate", stepHandler.ActivateStep)
		protected.DELETE("/presets/:presetId/steps/:stepId", stepHandler.DeleteStep)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}