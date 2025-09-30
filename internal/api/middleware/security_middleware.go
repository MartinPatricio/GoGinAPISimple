package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders agrega headers de seguridad a cada respuesta.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		// Puedes agregar un Content-Security-Policy más estricto si es necesario
		// c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
