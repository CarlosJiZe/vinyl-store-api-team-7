# Vinyl Store API

## Project Description

This project is a REST API built in Go using the Gin Framework, designed to manage the inventory of a vinyl record store inspired by platforms like [Merchbar](https://www.merchbar.com).

The API supports token-based authentication, meaning every user must log in to obtain an access token before performing any operation. Multiple users can be authenticated and operating simultaneously thanks to mutex-protected shared state.

The system allows authenticated users to:

- Browse the full album catalog
- Search for a specific album by ID
- Add new albums to the store
- Check the system status and their session info
- Log out and revoke their access token

The API includes full error handling on every endpoint, automatic ID generation for new albums, and case-insensitive duplicate detection by title and artist.

---

## Getting Started

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- Git

### Download the project

```bash
git clone https://github.com/CarlosJiZe/vinyl-store-api-team-7.git
cd vinyl-store-api-team-7
```

### Install dependencies

```bash
go mod tidy
```

### Run the server

```bash
go run main.go
```

The server will start on `http://localhost:8080`. You should see:

```
[GIN-debug] GET    /login
[GIN-debug] GET    /logout
[GIN-debug] GET    /albums
[GIN-debug] GET    /albums/:id
[GIN-debug] POST   /createAlbum
[GIN-debug] GET    /status
[GIN-debug] Listening and serving HTTP on :8080
```

### Connect a client

Open a new terminal and use `curl` to interact with the API. You can open multiple terminals to simulate multiple users connected at the same time.

---

## Default Users

| Username | Password |
|----------|----------|
| carlos   | 1234     |
| ana      | 4321     |
| loco     | 0000     |

---

## Usage of Every Endpoint

### GET /login

Authenticates a user using Basic Auth and returns an access token. This is the only endpoint that does not require a token.

```bash
curl -u carlos:1234 http://localhost:8080/login
```

Response:

```json
{
  "message": "Hi carlos, welcome to the Store System",
  "token": "YEBphr6DYf"
}
```

If the same user logs in again, the previous token is automatically revoked and a new one is issued:

```bash
curl -u carlos:1234 http://localhost:8080/login
```

```json
{
  "message": "Hi carlos, welcome to the Store System",
  "token": "RE0abhPH9h"
}
```

Error cases:

```bash
curl -u carlos:wrongpassword http://localhost:8080/login
```

```json
{ "error": "Usuario o contraseña incorrectos" }
```

```bash
curl http://localhost:8080/login
```

```json
{ "error": "Credenciales requeridas" }
```

---

### GET /logout

Revokes the current user's token. After logout, the token can no longer be used.

```bash
curl -H "Authorization: Bearer YEBphr6DYf" http://localhost:8080/logout
```

Response:

```json
{
  "message": "Bye carlos, your token has been revoked"
}
```

Attempting to use the token after logout:

```bash
curl -H "Authorization: Bearer YEBphr6DYf" http://localhost:8080/albums
```

```json
{ "error": "Token invalido o expirado" }
```

---

### GET /albums

Returns all albums currently in the store.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums
```

Response:

```json
[
  { "id": "1", "title": "Blue Train", "artist": "John Coltrane", "price": 56.99 },
  { "id": "2", "title": "Time Out", "artist": "Dave Brubeck", "price": 37.99 },
  { "id": "3", "title": "Flying Beagle", "artist": "Himiko Kikuchi", "price": 69.99 }
]
```

Attempting to access without a token:

```bash
curl http://localhost:8080/albums
```

```json
{ "error": "Token requerido" }
```

---

### GET /albums/:id

Returns a single album by its ID.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums/2
```

Response:

```json
{ "id": "2", "title": "Time Out", "artist": "Dave Brubeck", "price": 37.99 }
```

If the album does not exist:

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/albums/99
```

```json
{ "error": "Album no encontrado" }
```

---

### POST /createAlbum

Adds a new album to the store. The user only needs to provide title, artist, and price. The ID is assigned automatically based on the current catalog size.

**Mac/Linux:**

```bash
curl -X POST http://localhost:8080/createAlbum \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Kind of Blue", "artist": "Miles Davis", "price": 45.99}'
```

**Windows CMD:**

```cmd
curl -X POST http://localhost:8080/createAlbum -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" -d "{\"title\": \"Kind of Blue\", \"artist\": \"Miles Davis\", \"price\": 45.99}"
```

**Windows PowerShell:**

```powershell
curl -X POST http://localhost:8080/createAlbum `
  -H "Authorization: Bearer <TOKEN>" `
  -H "Content-Type: application/json" `
  -d '{"title": "Kind of Blue", "artist": "Miles Davis", "price": 45.99}'
```

Response:

```json
{ "id": "4", "title": "Kind of Blue", "artist": "Miles Davis", "price": 45.99 }
```

Attempting to add a duplicate (case-insensitive):

```bash
curl -X POST http://localhost:8080/createAlbum \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "blue train", "artist": "john coltrane", "price": 39.99}'
```

```json
{ "error": "Ya existe un album de este artista con ese titulo" }
```

Attempting to add an album with a negative price:

```bash
curl -X POST http://localhost:8080/createAlbum \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Kind of Blue", "artist": "Miles Davis", "price": -10}'
```

```json
{ "error": "El precio debe ser un numero positivo" }
```

Attempting to add an album with missing fields:

```bash
curl -X POST http://localhost:8080/createAlbum \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "", "artist": "", "price": 45.99}'
```

```json
{ "error": "El titulo y el artista son requeridos" }
```

Attempting to add an album without price:

```bash
curl -X POST http://localhost:8080/createAlbum \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Kind of Blue", "artist": "Miles Davis"}'
```

```json
{ "error": "El precio debe ser un numero positivo" }
```

---

### GET /status

Returns the current system status and the username of the authenticated user.

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/status
```

Response:

```json
{
  "message": "Hi carlos, the DPIP System is Up and Running",
  "time": "2026-04-24 13:27:35"
}
```

---

## Authentication

All endpoints except `/login` require a valid Bearer token in the `Authorization` header:

```
Authorization: Bearer <TOKEN>
```

Tokens are revoked on logout or when the same user logs in again. Multiple users can be authenticated and operating simultaneously.

---

## Features for Future Work

### 1. Web UI

Build a frontend interface using React or Vue.js that allows users to browse the vinyl catalog, search albums, and manage inventory visually without needing to use curl commands. This would make the system accessible to non-technical users.

### 2. User Registration

Add a `/register` endpoint that allows new users to create accounts with a username and password. Passwords would be stored securely using bcrypt hashing instead of plain text, and the user data would persist to a JSON file so it survives server restarts.

### 3. Role-Based Access Control

Implement two user roles: customer and admin. Customers can only browse albums (`GET /albums`, `GET /albums/:id`). Admins have full access including adding new albums (`POST /createAlbum`). Roles would be stored per user and enforced by the authentication middleware, returning a `403 Forbidden` response when a customer attempts a restricted operation.

---

## Team

- Carlos Jimenez Zepeda
- [Compañero]

**Professor:** Alfredo Emmanuel Garcia Falcon  
**Course:** Advanced Parallel Programming  
**Institution:** Universidad Panamericana Guadalajara
