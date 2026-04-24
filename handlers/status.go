package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStatus regresa el estdo del sistema y el usuario loggeado
func GetStatus(c *gin.Context) {
	//Obtenemos el nombre de usuario del contexto, lo pusimos ahi en el middleware de autenticacion
	username, _ := c.Get("username")

	//Respondemos con un mensaje de estado y el nombre del usuario loggeado
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username.(string) + ", the DPIP Sysrem is UP and Running",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}
