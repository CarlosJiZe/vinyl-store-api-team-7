package middleware

import (
	"net/http"
	"strings"

	"github.com/CarlosJiZe/vinyl-store-api-team-7/models"
	"github.com/gin-gonic/gin"
)

// AuthRequired verifica que el request tenga un token valido
func AuthRequired() gin.HandlerFunc { //Funcion que se ejecuta antes de llegar al endpoint real
	return func(c *gin.Context) {
		// Obtenemos el header Authorization
		authHeader := c.GetHeader("Authorization")

		// Revisamos que tenga el formato "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort() //Evitamos que el request continue a los handlers si no tiene un token valido
			return
		}

		// Extraemos el token sin el prefijo "Bearer "
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Verificamos que el token exista en nuestro mapa de tokens activos
		models.TokenMutex.RLock() // Bloqueamos para lectura
		username, exists := models.Tokens[token]
		models.TokenMutex.RUnlock() // Desbloqueamos

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido o expirado"})
			c.Abort()
			return
		}

		// Guardamos el nombre de usuario en el contexto para que los handlers puedan acceder a el
		c.Set("username", username)
		c.Next() // Continuamos con el siguiente handler

	}
}
