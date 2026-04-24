package models

import "sync"

// User representa un usuario en el sistema
type User struct {
	Username string
	Password string
}

// User es la "base de datos" en memoria para los usuarios
var Users = []User{
	{Username: "carlos", Password: "1234"},
	{Username: "ana", Password: "4321"},
	{Username: "loco", Password: "0000"},
}

// Tokens hace un mapeo de los tokens de autenticación a los nombres de usuario
var Tokens = make(map[string]string)

// TokenMutex protege el acceso concurrente al mapa de Tokens
var TokenMutex sync.RWMutex
