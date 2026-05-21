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

		method := c.Request.Method
		path := c.FullPath()

		result, err := m.permissionService.CheckRouteAccess(entityID.(uint), method, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission check failed"})
			c.Abort()
			return
		}

		if !result.Allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("access_scope", string(result.Scope))
		if result.OwnershipField != "" {
			c.Set("ownership_field", result.OwnershipField)
		}

		c.Next()
	}
}