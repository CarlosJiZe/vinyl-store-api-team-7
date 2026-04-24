package handlers

import (
	"math/rand"
	"net/http"

	"github.com/CarlosJiZe/vinyl-store-api-team-7/models"
	"github.com/gin-gonic/gin"
)

// generateToken genera un token aleatorio de 10 caracteres
func generateToken() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" //Caracteres posibles
	token := make([]byte, 10)                                                 //Slice de bytes para el token, son 10 espacios vacios
	for i := range token {                                                    //Va generando numeros de 0 al 62 y selecciona un caracter  de chars hasta llenar el token
		token[i] = chars[rand.Intn(len(chars))]
	}
	return string(token)
}

// Login autentica al usuario y le da un token
func Login(c *gin.Context) {
	//Gin nos ayuda a extraer el usuario y contraseña del header
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales requeridas"})
		return
	}

	//Buscamos el usuario en nuestra "base de datos"
	validUser := false
	for _, user := range models.Users {
		if user.Username == username && user.Password == password {
			validUser = true
			break
		}
	}

	//Si no es valido, respondemos con un error
	if !validUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// Si el usuario ya tiene un token activo, lo revocamos antes de generar uno nuevo
	models.TokenMutex.Lock() //Bloqueamos para escritura
	for token, user := range models.Tokens {
		if user == username {
			delete(models.Tokens, token) //Eliminamos el token del mapa
			break
		}
	}
	models.TokenMutex.Unlock() //Desbloqueamos

	//Generamos el token y lo guardamos en el mapa
	token := generateToken()
	models.TokenMutex.Lock() //Bloqueamos para escritura
	models.Tokens[token] = username
	models.TokenMutex.Unlock() //Desbloqueamos

	//Respondemos con el token
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi " + username + ", welcome to the Store System",
		"token":   token,
	})
}

// Logout con el que se revoca el token al usuario
func Logout(c *gin.Context) {
	//Obtenemos el username del contexto, que fue guardado por el middleware de autenticacion
	username, _ := c.Get("username")

	//Buscamos y eliminamos su token del mapa
	models.TokenMutex.Lock() //Bloqueamos
	for token, user := range models.Tokens {
		if user == username {
			delete(models.Tokens, token) //Eliminamos el token del mapa
			break
		}
	}
	models.TokenMutex.Unlock() //Desbloqueamos

	//Respondemos con un menssaje de despedida
	c.JSON(http.StatusOK, gin.H{
		"message": "Bye " + username.(string) + ", your token has been revoked",
	})
}
