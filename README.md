# TP08 - Sistema de Integración y Despliegue

**Materia:** Ingeniería de Software 3
**Alumno:** Octavio Carpineti - Kevin Massholder
**Año:** 2025

Mini red social completa con PostgreSQL, entornos QA/PROD separados, Railway databases, Render deployment, y suite completa de pruebas unitarias (42 tests).

---

## 📋 Tabla de Contenidos

- [Tecnologías](#tecnologías)
- [Arquitectura](#arquitectura)
- [Funcionalidades](#funcionalidades)
- [Prerequisitos](#prerequisitos)
- [Instalación](#instalación)
- [Ejecución](#ejecución)
- [Testing](#testing)
- [CI/CD](#cicd)
- [Estructura del Proyecto](#estructura-del-proyecto)

---

## 🛠️ Tecnologías

### Backend
- **Go 1.21+**
- **PostgreSQL** (Railway cloud databases)
- **Gorilla Mux** (routing)
- **lib/pq** (PostgreSQL driver)
- **testify** (testing + mocking)

### Frontend
- **React 18** con **TypeScript**
- **Axios** (HTTP client)
- **Jest** + **React Testing Library** (testing)

### Infraestructura
- **Railway** (PostgreSQL databases)
- **Render** (deployment platform)
- **GitHub Actions** (CI/CD)
- **Docker** (containerization)

---

## 🚀 Despliegue y Arquitectura

### Entornos de Despliegue
```
┌─────────────────┐    ┌─────────────────┐
│   Frontend QA   │    │  Frontend PROD  │
│ Render Service  │    │ Render Service  │
└─────────┬───────┘    └───────┬─────────┘
          │                    │
          │                    │
          ▼                    ▼
┌─────────────────┐    ┌─────────────────┐
│   Backend QA    │    │  Backend PROD   │
│ Render Service  │    │ Render Service  │
│                 │    │                 │
│ DATABASE_URL →  │    │ DATABASE_URL →  │
└─────────┬───────┘    └───────┬─────────┘
          │                    │
          ▼                    ▼
┌─────────────────┐    ┌─────────────────┐
│ PostgreSQL QA   │    │ PostgreSQL PROD │
│   Railway DB    │    │   Railway DB    │
└─────────────────┘    └─────────────────┘
```

### Arquitectura por Capas

```
Frontend (React)     →      Backend (Go)
──────────────────────     ────────────────────
Login/PostList         ┌─►  Handlers     (HTTP handlers)
React Components       │    ├── auth_handler.go
API Calls (axios)      │    └── post_handler.go
                        │
                        │    Services     (business logic)
                        ├── auth_service.go      ───┐
                        └── post_service.go           │
                                                      │ MOCK repository
                        Repository   (data access)   │ (for testing)
                        ├── user_repository.go ──┐   │
                        └── post_repository.go ──┐┼───┘
                                                  │
PSQL Repository                  PostgreSQL
(SELECT/INSERT/UPDATE)          (Railway Cloud)
```

### Configuración de Base de Datos

**Esquema PostgreSQL:**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
```

---

## ✨ Funcionalidades

### Autenticación
- ✅ Registro de usuarios con validación
- ✅ Login con email/password
- ✅ JWT-like session handling (headers)
- ✅ CORS configurado para cross-origin

### Posts y Comentarios
- ✅ Crear post con título y contenido
- ✅ Listar posts de todos los usuarios
- ✅ Ver detalle de post con comentarios
- ✅ Eliminar post (solo autor)
- ✅ Comentar en posts
- ✅ Eliminar comentarios (solo autor)

### Validaciones de Negocio
- 🔒 **Autorización**: Solo el autor puede eliminar posts/comentarios
- ✉️ **Email**: Validación de formato y unicidad
- 🔑 **Password**: Mínimo 6 caracteres
- 📝 **Posts**: Título mínimo 3 caracteres
- 🗃️ **Base de Datos**: Constraints a nivel DB (foreign keys, serial IDs)

### Separación de QA/PROD
- ✅ **Bases de datos independientes**: QA y PROD no comparten datos
- ✅ **URLs separadas**: Cada entorno tiene su propia URL
- ✅ **Variables de entorno**: Configuración por entorno

---

## 📦 Prerequisitos

### Cuentas y Servicios Externos

#### Railway (Base de Datos PostgreSQL)
1. Registrarse en [Railway.app](https://railway.app)
2. Agregar método de pago (requerido para PostgreSQL)
3. Crear proyecto: **ingsw3-tp08-qa** y **ingsw3-tp08-prod**

#### Render (Despliegue)
1. Registrarse en [Render.com](https://render.com)
2. Conectar repositorio de GitHub
3. Crear servicios separados para QA y PROD

### Instalación de Herramientas Locales

#### Go (Backend)
```bash
go version  # Debe ser 1.21+
```

#### Node.js (Frontend)
```bash
node --version  # Debe ser 18+
npm --version
```

#### Git
```bash
git --version
```

---

## 🗄️ Configuración de Base de Datos (Railway)

### 1. Crear Base de Datos QA
1. **Railway Dashboard** → **New Project** → **Provision PostgreSQL**
2. Nombre: `ingsw3-tp08-qa`
3. Plan: **Hobby** (512MB RAM, 1GB storage)
4. Crear y esperar configuración (~2-3 minutos)

### 2. Crear Base de Datos PROD
1. Repetir proceso para PROD
2. Nombre: `ingsw3-tp08-prod`
3. Plan: **Hobby** (libre para uso básico)

### 3. Configurar Esquema
**Para cada base de datos:**
1. Ir a → **Variables** → **Query** tab
2. Ejecutar el esquema de arriba (users, posts, comments)

### 4. Obtener URLs de Conexión
**Para cada DB:**
- Ir a **"Variables"** tab
- Copiar **`DATABASE_URL`** value

Ejemplo: `postgresql://postgres:abcd1234@us-west1-postgres-xyz.railway.app:5432/railway`

---

## 🎪 Despliegue en Render

### 1. Servicios Backend (QA y PROD)

#### Backend QA:
1. **Render Dashboard** → **New** → **Web Service**
2. **Conectar GitHub repo**: `Kevinmass/IngSWIII-TP08`
3. **Configurar servicio:**
   - **Name**: `ingsw3-back-qa`
   - **Root Directory**: `./backend`
   - **Environment**: `Go`
   - **Go Version**: `1.21`
   - **Build Command**: `go mod download`
   - **Start Command**: `go run cmd/api/main.go`

4. **Environment Variables:**
   - **DATABASE_URL**: `[tu QA Railway DATABASE_URL]`

#### Backend PROD:
- Repetir con nombre: `ingsw3-back-prod`
- Usar PROD Railway DATABASE_URL

### 2. Servicios Frontend (QA y PROD)

#### Frontend QA:
1. **Render Dashboard** → **New** → **Static Site**
2. **Conectar repo**: `Kevinmass/IngSWIII-TP08`
3. **Configurar:**
   - **Name**: `ingsw3-front-qa`
   - **Root Directory**: `./frontend`
   - **Build Command**: `npm install && npm run build`
   - **Publish Directory**: `build`

#### Frontend PROD:
- Repetir con nombre: `ingsw3-front-prod`

### 3. Variables de Entorno Frontend
**Los frontend services necesitan variables de entorno definidas por Render:**

#### Frontend QA:
- **REACT_APP_BACKEND_URL**: `https://ingsw3-back-qa.onrender.com`

#### Frontend PROD:
- **REACT_APP_BACKEND_URL**: `https://ingsw3-back-prod.onrender.com`

**NOTA:** Las URLs de Render se generan automáticamente. Reemplazar con URLs reales una vez creados los servicios backend.

---

## 🖥️ Desarrollo Local

### 1. Instalar Dependencias
```bash
git clone https://github.com/Kevinmass/IngSWIII-TP08.git
cd IngSWIII-TP08

# Backend
cd backend
go mod download

# Frontend
cd ../frontend
npm install
```

### 2. Ejecución Local
**Backend (Terminal 1):**
```bash
cd backend
# Agregar DATABASE_URL si quieres usar PostgreSQL local
DATABASE_URL="postgresql://..." go run cmd/api/main.go
# O usar valor por defecto (error si no se configura)
```

**Frontend (Terminal 2):**
```bash
cd frontend
npm start
```

---

## ▶️ Ejecución

### Opción A: Ejecutar Backend y Frontend por separado

#### Terminal 1 - Backend
```bash
cd backend
go run cmd/api/main.go
```

Deberías ver:
```
Base de datos inicializada correctamente
🚀 Servidor corriendo en http://localhost:8080
```

#### Terminal 2 - Frontend
```bash
cd frontend
npm start
```

Se abrirá automáticamente en: `http://localhost:3000`

### Opción B: Script para ejecutar ambos (Linux/Mac)

```bash
# Crear script
cat > run.sh << 'EOF'
#!/bin/bash
cd backend && go run cmd/api/main.go &
BACKEND_PID=$!
cd ../frontend && npm start
kill $BACKEND_PID
EOF

chmod +x run.sh
./run.sh
```

---

## 🧪 Testing

### Backend Tests (Go)

```bash
cd backend

# Ejecutar todos los tests
go test ./tests/services/... -v

# Con cobertura
go test ./tests/services/... -v -cover

# Solo un test específico
go test ./tests/services/ -v -run TestRegister_Success
```

**Resultado esperado:**
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
...
PASS
ok      tp06-testing/tests/services     0.582s
```

**Total: 23 tests** ✅

### Frontend Tests (React)

```bash
cd frontend

# Ejecutar todos los tests
npm test

# Con cobertura
npm test -- --coverage

# Sin modo watch
npm test -- --watchAll=false
```

**Resultado esperado:**
```
PASS  src/components/Login/Login.test.tsx
PASS  src/components/PostList/PostList.test.tsx
PASS  src/components/CommentList/CommentList.test.tsx
PASS  src/services/authService.test.ts

Test Suites: 4 passed, 4 total
Tests:       19 passed, 19 total
```

**Total: 19 tests** ✅

### Ejecutar TODOS los tests (Backend + Frontend)

```bash
# Desde la raíz del proyecto
cd backend && go test ./... && cd ../frontend && npm test -- --watchAll=false
```

---

## 🔄 CI/CD

### GitHub Actions

El proyecto incluye un pipeline de CI/CD que se ejecuta automáticamente en cada push.

**Archivo:** `.github/workflows/ci.yml`

**Workflow:**
1. ✅ **Backend Tests** - Ejecuta `go test`
2. ✅ **Frontend Tests** - Ejecuta `npm test`
3. ✅ **Backend Build** - Compila con `go build`
4. ✅ **Frontend Build** - Compila con `npm run build`
5. ✅ **Summary** - Resumen final

**Ver resultados:**
1. Ir a: `https://github.com/TU-USUARIO/tp06-testing/actions`
2. Seleccionar el workflow más reciente
3. Ver logs detallados de cada job

---

## 📁 Estructura del Proyecto

```
IngSWIII-TP08/
├── .github/
│   └── workflows/
│       └── ci-cd.yml               # Pipeline CI/CD con Render deployment
│
├── backend/
│   ├── cmd/api/
│   │   └── main.go                 # Punto de entrada (PostgreSQL-only)
│   ├── internal/
│   │   ├── database/
│   │   │   └── database.go         # PostgreSQL initialization + auto-schema
│   │   ├── models/                 # Structs (User, Post, Comment)
│   │   │   ├── users.go
│   │   │   ├── post.go
│   │   ├── repository/             # PostgreSQL data access
│   │   │   ├── user_repository.go  # PostgreSQL with $1, $2 placeholders
│   │   │   └── post_repository.go
│   │   ├── services/               # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   └── post_service.go
│   │   ├── handlers/               # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   └── post_handler.go
│   │   └── router/
│   │       └── router.go           # Routes + CORS middleware
│   ├── tests/                      # Unit tests with mocks
│   │   ├── mocks/
│   │   │   ├── user_repository_mock.go
│   │   │   └── post_repository_mock.go
│   │   └── services/
│   │       ├── auth_service_test.go
│   │       └── post_service_test.go
│   ├── Dockerfile                  # Go 1.21 + PostgreSQL
│   ├── go.mod                      # PostgreSQL-only dependencies
│   └── go.sum                      # Lockfile checksums
│
├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── components/             # React components
│   │   │   ├── Login/
│   │   │   │   ├── Login.tsx
│   │   │   │   ├── Login.test.tsx
│   │   │   │   └── Login.css
│   │   │   ├── PostList/
│   │   │   │   ├── PostList.tsx
│   │   │   │   ├── PostList.test.tsx
│   │   │   │   └── PostList.css
│   │   │   ├── CreatePost/
│   │   │   ├── CommentList/
│   │   │   ├── CommentForm/
│   │   │   └── PostDetail/
│   │   ├── services/               # API services (env-aware)
│   │   │   ├── authService.ts      # Auto-detect backend URL
│   │   │   ├── postService.ts
│   │   │   └── authService.test.ts
│   │   ├── __mocks__/
│   │   │   └── axios.ts            # HTTP mocking
│   │   ├── types/
│   │   │   └── index.ts            # TypeScript definitions
│   │   ├── App.tsx
│   │   └── setupTests.ts
│   ├── Dockerfile                  # Multi-stage Node.js build
│   ├── package.json
│   └── tsconfig.json
│
├── error-log.txt                   # Deployment troubleshooting logs
├── decisiones.md                   # Technical documentation
└── README.md                       # This file
```

### 🔗 URLs y Endpoints

#### Despliegue Actual (Render):
- **Frontend QA**: `https://ingsw3-front-qa.onrender.com`
- **Backend QA**: `https://ingsw3-back-qa.onrender.com`
- **Frontend PROD**: `https://ingsw3-front-prod.onrender.com`
- **Backend PROD**: `https://ingsw3-back-prod.onrender.com`

#### API Endpoints:
```
POST   /api/auth/register     # User registration
POST   /api/auth/login        # User login
GET    /api/posts             # List all posts
POST   /api/posts             # Create new post
GET    /api/posts/:id         # Get post details
DELETE /api/posts/:id         # Delete post (author only)
GET    /api/posts/:id/comments    # Get post comments
POST   /api/posts/:id/comments    # Add comment
DELETE /api/posts/:postId/comments/:commentId  # Delete comment (author only)
```

### 🚀 Desarrollo vs Producción

**Desarrollo Local:**
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- Base de Datos: Railway PostgreSQL (ambos entornos)

**Entorno de Producción:**
- Frontend: Static site served by Render
- Backend: Go server on Render
- Base de Datos: Railway PostgreSQL (QA y PROD separados)

---

## 📊 Cobertura de Tests

### Backend (23 tests)

| Componente  | Tests |                 Descripción                     |
|-------------|-------|-------------------------------------------------|
| AuthService | 11    | Register (6), Login (5)                         |
| PostService | 12    | CreatePost (5), DeletePost (3), DeleteComment(4)|

### Frontend (19 tests)

| Componente | Tests |            Descripción             |
|------------|-------|------------------------------------|
| Login      | 5     | Renderizado, validaciones, estados |
| PostList   | 5     | Renderizado, eliminación, permisos |
| CommentList| 5     | Renderizado, eliminación, permisos |
| authService| 4     | Login/Register con mocks HTTP      |

**Total: 42 tests automatizados** ✅

---

## 🎯 Conceptos Implementados

### Testing
- ✅ **Pruebas Unitarias** (backend + frontend)
- ✅ **Patrón AAA** (Arrange, Act, Assert)
- ✅ **Mocking** (Repository + HTTP)
- ✅ **Aislamiento** de dependencias
- ✅ **Casos edge** y validaciones

### Arquitectura
- ✅ **Separación de concerns** (capas)
- ✅ **Dependency Injection** (interfaces)
- ✅ **Repository Pattern**
- ✅ **RESTful API**

### DevOps
- ✅ **CI/CD** con GitHub Actions
- ✅ **Automatización** de tests
- ✅ **Build automático**

---

## 🔍 Comandos Útiles

### Backend
```bash
# Compilar
go build ./...

# Tests (con mocks, no requieren DB)
go test ./tests/services/... -v

# Tests de integración (requieren PostgreSQL)
go test ./... -v

# Verificar dependencias
go mod verify
go mod tidy
```

### Frontend
```bash
# Desarrollo
npm start

# Tests
npm test

# Build producción
npm run build

# Limpiar node_modules
rm -rf node_modules && npm install
```

### Git
```bash
# Status
git status

# Commit
git add .
git commit -m "mensaje"

# Push
git push origin main
```

---

## 📚 Documentación Adicional

- **[decisiones.md](./decisiones.md)** - Decisiones técnicas y justificaciones
- **[backend/tests/desc.md](./backend/tests/desc.md)** - Explicación de tests backend
- **[backend/internal/database/desc.md](./backend/internal/database/desc.md)** - Explicación de base de datos
- **[frontend/src/services/desc.md](./frontend/src/services/desc.md)** - Explicación de servicios HTTP

---

## 🐛 Troubleshooting

### El backend no arranca
```bash
# Verificar que no esté corriendo en otro lado
lsof -i :8080
kill -9 PID_DEL_PROCESO

# Verificar dependencias
cd backend
go mod tidy
```

### El frontend no arranca
```bash
# Reinstalar dependencias
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Los tests fallan
```bash
# Backend: Verificar que no dependa de BD
rm backend/database.db
go test ./tests/services/... -v  # Deben pasar igual

# Frontend: Limpiar cache de Jest
npm test -- --clearCache
npm test
```

### CORS errors
Verificar que el backend tenga el middleware CORS configurado en `router/router.go`

---

## 👥 Autores:
**Carpineti Octavio - Kevin Massholder**  
Ingenieria en sistemas de informacion - UCC
Materia: Ingeniería de Software 3  
Año: 2025
