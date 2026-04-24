package models

// Album representa un disco de vini en la tienda
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Almbums es la "base de datos" en memoria con albumes de inicio
var Albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Time Out", Artist: "Dave Brubeck", Price: 37.99},
	{ID: "3", Title: "Flying Beagle", Artist: "Himiko Kikuchi", Price: 69.99},
}
