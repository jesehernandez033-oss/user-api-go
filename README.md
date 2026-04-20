# 🚀 User Management API (Go)

API REST desarrollada en Go para la gestión de usuarios, con autenticación JWT, validaciones de negocio y conexión a base de datos SQL Server.

---

## 📌 Características

* 🔐 Autenticación con JWT
* 👤 Crear, actualizar y eliminar usuarios
* ✅ Validaciones:

  * Campos obligatorios
  * Contraseña segura
  * Edad mínima (14 años)
  * Email único
* 🗄️ Conexión a base de datos SQL Server
* 🧪 Probado con Postman

---

## 🧱 Arquitectura

El proyecto sigue una estructura simple inspirada en Clean Architecture:

```
/application
    /controllers   → Manejo de endpoints
/application
    token.go       → Lógica de JWT
/infrastructure
    db.go          → Conexión a base de datos
main.go            → Entry point
```

---

## ⚙️ Requisitos

* Go 1.18+
* SQL Server
* Postman (para pruebas)

---

## ▶️ Cómo ejecutar el proyecto

### 1. Clonar repositorio

```bash
git clone <repo-url>
cd <project-folder>
```

---

### 2. Configurar base de datos

Crear base de datos:

```sql
CREATE DATABASE api_users;
```

Tabla:

```sql
CREATE TABLE users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    username VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100),
    phone_number VARCHAR(20),
    birthday DATE,
    address VARCHAR(255)
);
```

---

### 3. Configurar conexión

En `infrastructure/db.go`:

```go
connString := "sqlserver://localhost?database=api_users&trusted_connection=yes&encrypt=disable"
```

---

### 4. Ejecutar API

```bash
go run main.go
```

Servidor corriendo en:

```text
http://localhost:8080
```

---

## 🔐 Autenticación

### Login

```http
POST /login
```

Respuesta:

```json
{
  "Message": "loginSuccesful",
  "AccessToken": "TOKEN"
}
```

---

### Uso del token

En los endpoints protegidos:

```http
Authorization: Bearer TU_TOKEN
```

---

## 📡 Endpoints

---

### 🟢 Health Check

```http
GET /ping
```

Respuesta:

```
pong
```

---

### 🟢 Crear usuario

```http
POST /createUser
```

Body:

```json
{
  "username": "jese",
  "email": "test@test.com",
  "password": "Password1!",
  "phone_number": "8091234567",
  "birthday": "2000-01-01",
  "address": "Santo Domingo"
}
```

---

### 🟡 Actualizar usuario

```http
PUT /updateUser
```

Body:

```json
{
  "email": "test@test.com",
  "username": "nuevo",
  "password": "Password1!",
  "phone_number": "8099999999",
  "birthday": "2000-01-01",
  "address": "Nueva direccion"
}
```

---

### 🔴 Eliminar usuario

```http
DELETE /deleteUser?email=test@test.com
```

---

## ✅ Validaciones implementadas

* Username, email y password obligatorios
* Password:

  * 8-16 caracteres
  * Al menos 1 mayúscula
  * 1 número
  * 1 carácter especial
* Edad mínima: 14 años
* Email único en base de datos
* Formato de fecha: `YYYY-MM-DD`

---

## ❌ Manejo de errores

| Código | Descripción           |
| ------ | --------------------- |
| 400    | Datos inválidos       |
| 401    | Token inválido        |
| 404    | Usuario no encontrado |
| 500    | Error interno         |

---

## 🧪 Pruebas

Se realizaron pruebas con:

* Datos válidos
* Password inválido
* Usuario menor de edad
* Email duplicado
* Token inválido o ausente

---

## 🚀 Mejoras futuras

* 🔐 Encriptar contraseñas con bcrypt
* 🧱 Separar capas (services, repositories)
* 🧪 Tests automatizados más completos
* 🐳 Docker completo con DB integrada

---

## 👨‍💻 Autor

Proyecto desarrollado como ejercicio técnico de API en Go.

---
