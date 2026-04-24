package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CarlosJiZe/vinyl-store-api-team-7/models"
	"github.com/gin-gonic/gin"
)

// GetAlbums regresa todos los albumes de la de tienda
func GetAlbums(c *gin.Context) {
	models.TokenMutex.RLock() //Bloqueamos para lectura
	albums := models.Albums
	models.TokenMutex.RUnlock() //Desbloqueamos

	c.JSON(http.StatusOK, albums) //Respondemos con el slice de albumes en formato JSON
}

// GetAlbumByID regresa un album por su ID
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id") //Obtenemos el ID del album de la URL

	for _, album := range models.Albums { //Iteramos sobre los albumes para encontrar el que tenga el ID solicitado
		if album.ID == id {
			c.JSON(http.StatusOK, album) //Si lo encontramos respondemos con el album en formato JSON
			return
		}
	}

	//Si no encontramos el album respondemos con un error
	c.JSON(http.StatusNotFound, gin.H{"error": "Album no encontrado"})
}

// CreateAlbum agrega un nuevo album
func CreateAlbum(c *gin.Context) {

	//Input es lo que debe mandar el usuario
	var input struct {
		Title  string  `json:"title"`
		Artist string  `json:"artist"`
		Price  float64 `json:"price"`
	}

	//Intentamos hacer el parseo del json
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos ivalidos: " + err.Error()})
		return
	}

	//Validamos que los campos no esten vacios
	if input.Title == "" || input.Artist == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El titulo y el artista son requeridos"})
		return
	}

	//Validamos que el precio sea positivo
	if input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El precio debe ser un numero positivo"})
		return
	}

	//Verificamos duplicados por titulo y artista
	for _, album := range models.Albums {
		if strings.EqualFold(album.Title, input.Title) && strings.EqualFold(album.Artist, input.Artist) {
			c.JSON(http.StatusConflict, gin.H{"error": "Ya existe un albun de este artista con ese titulo"})
			return
		}
	}

	//Generamos el id en base al tamaño del slice
	newID := fmt.Sprintf("%d", len(models.Albums)+1)

	//Creamos el nuevo album con los datos del input y el ID generado
	newAlbum := models.Album{
		ID:     newID,
		Title:  input.Title,
		Artist: input.Artist,
		Price:  input.Price,
	}

	models.Albums = append(models.Albums, newAlbum) //Agregamos el nuevo album al slice de albumes

	c.JSON(http.StatusCreated, newAlbum) //Respondemos con el nuevo album en formato JSON
}
