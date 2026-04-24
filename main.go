package main

import (
	"github.com/CarlosJiZe/vinyl-store-api-team-7/handlers"
	"github.com/CarlosJiZe/vinyl-store-api-team-7/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	//Creamos el router de Gin
	r := gin.Default()

	//Rutas publicas
	r.GET("/login", handlers.Login)

	//Rutas protegidas por autenticacion
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired()) //Aplicamos el middleware de autenticacion a este grupo de rutas
	{
		protected.GET("/logout", handlers.Logout)
		protected.GET("/albums", handlers.GetAlbums)
		protected.GET("/albums/:id", handlers.GetAlbumByID)
		protected.POST("/createAlbum", handlers.CreateAlbum)
		protected.GET("/status", handlers.GetStatus)
	}

	//Iniciamos el servidor en el puerto 8080
	r.Run(":8080")
}
