Listo. Este es tu **README completo, limpio y profesional**, con uso incluido. Solo cópialo y pégalo en `README.md`:

---

```md
# User API Go - Version 2.0

API REST desarrollada en Go para la gestión de usuarios con autenticación JWT, validaciones de seguridad y arquitectura limpia.

Incluye funcionalidades de registro, login, actualización, eliminación de usuarios y subida de archivos.

---

## 🌐 Base URL

```

[http://localhost:8080](http://localhost:8080)

````

---

## 🚀 Versión 2.0

Mejoras implementadas:

- Autenticación segura con JWT
- Encriptación de contraseñas con bcrypt
- Validación de datos (email, password, edad)
- Rate limiting en login (máx. 5 intentos por minuto)
- Manejo seguro de subida de archivos (sin panic)
- Logging centralizado (INFO, WARNING, ERROR)
- Configuración mediante variables de entorno
- Mejora de CORS (soporte PUT y DELETE)
- Update no destructivo (solo actualiza campos enviados)
- Separación de lógica (Clean Architecture)

---

## 🧠 Tecnologías

- Go (Golang)
- SQL Server
- JWT (github.com/golang-jwt/jwt)
- bcrypt
- Gorilla Mux
- CORS middleware

---

## 🔐 Seguridad

- Contraseñas encriptadas con bcrypt
- Autenticación con JWT
- Protección contra fuerza bruta (rate limit)
- Validación de inputs
- Manejo seguro de errores
- Logging para auditoría

---

## 🏗️ Arquitectura

El proyecto sigue principios de Clean Architecture:

- `controllers` → manejo de endpoints HTTP  
- `domain` → lógica de negocio y validaciones  
- `infrastructure` → DB, logging  
- `middleware` → autenticación  

---

## 📡 Endpoints

### POST /createUser
Crea un usuario

### POST /login
Autentica y devuelve token JWT

### PUT /updateUser
Actualiza usuario (requiere token)

### DELETE /deleteUser
Elimina usuario (requiere token)

### POST /uploadFile
Sube archivos al servidor

---

## 🧪 Uso de la API

### 🔹 Crear usuario

```http
POST /createUser
Content-Type: application/json
````

```json
{
  "username": "jesse",
  "email": "jesse@test.com",
  "password": "Password1!",
  "phone_number": "8091234567",
  "birthday": "2000-05-10",
  "address": "Higuey"
}
```

---

### 🔹 Login

```http
POST /login
Content-Type: application/json
```

```json
{
  "email": "jesse@test.com",
  "password": "Password1!"
}
```

Respuesta:

```json
{
  "message": "login successful",
  "access_token": "TOKEN_AQUI"
}
```

---

### 🔹 Usar token

```http
Authorization: Bearer TU_TOKEN
```

---

### 🔹 Actualizar usuario

```http
PUT /updateUser
Authorization: Bearer TU_TOKEN
Content-Type: application/json
```

```json
{
  "email": "jesse@test.com",
  "username": "nuevo_nombre"
}
```

👉 Solo se actualizan los campos enviados

---

### 🔹 Eliminar usuario

```http
DELETE /deleteUser
Authorization: Bearer TU_TOKEN
Content-Type: application/json
```

```json
{
  "email": "jesse@test.com"
}
```

---

### 🔹 Subir archivo

```http
POST /uploadFile
Content-Type: multipart/form-data
```

Body:

```
file: (archivo)
```

---

## ⚙️ Variables de entorno

Antes de ejecutar:

```
DB_CONN=sqlserver://localhost:1433?database=api_users&trusted_connection=yes
JWT_SECRET=tu_clave_segura
ALLOWED_ORIGIN=http://localhost:3000
PORT=8080
```

---

## ▶️ Ejecución

```bash
go run main.go
```

---

## 🧪 Pruebas

Los endpoints fueron probados usando Postman incluyendo:

* Casos exitosos
* Errores controlados
* Seguridad (tokens inválidos)
* Rate limiting

---

## 📂 Estructura del proyecto

```
application/
  controllers/
  middleware/
domain/
infrastructure/
uploads/
```

---

## 📌 Notas

Proyecto desarrollado aplicando buenas prácticas de backend, seguridad y manejo de errores.

---

## 👤 Autor

Jesé Hernández

````





