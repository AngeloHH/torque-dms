package http

import (
	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/handlers"
	"torque-dms/adapters/input/http/middleware"
	"torque-dms/core/identity/ports/input"
)

type Router struct {
	engine            *gin.Engine
	authService       input.AuthService
	entityService     input.EntityService
	permissionService input.PermissionService
}

func NewRouter(
	authService input.AuthService,
	entityService input.EntityService,
	permissionService input.PermissionService,
	jwtSecret string,
) *Router {
	r := &Router{
		engine:            gin.Default(),
		authService:       authService,
		entityService:     entityService,
		permissionService: permissionService,
	}

	r.setupRoutes(jwtSecret)
	return r
}

func (r *Router) setupRoutes(jwtSecret string) {
	// Handlers
	authHandler := handlers.NewAuthHandler(r.authService)
	entityHandler := handlers.NewEntityHandler(r.entityService)

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
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}