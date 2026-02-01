package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"torque-dms/core/identity/ports/input"
)

type PermissionMiddleware struct {
	permissionService input.PermissionService
}

func NewPermissionMiddleware(permissionService input.PermissionService) *PermissionMiddleware {
	return &PermissionMiddleware{permissionService: permissionService}
}

func (m *PermissionMiddleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		entityID, exists := c.Get("entity_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "entity not found in context"})
			c.Abort()
			return
		}

		// Obtener el resource basado en método y path
		method := c.Request.Method
		path := c.FullPath()

		// Por ahora, permitir todo si está autenticado
		// Después implementar la búsqueda del resource y verificación
		_ = entityID
		_ = method
		_ = path

		c.Next()
	}
}