# TP7 + TP8 - Integración Completa: Quality Assurance + Contenedores

**Alumno:** Octavio Carpineti - Kevin Massholder
**Materia:** Ingeniería de Software III
**Fecha:** Noviembre 2025

**Integración:** TP7 (Pruebas, QA, SonarCloud) + TP8 (Contenedores, PostgreSQL, Deploy)

---

## 📋 Tabla de Contenidos

1. [Descripción del Proyecto](#descripción-del-proyecto)
2. [Arquitectura Integrada](#arquitectura-integrada)
3. [Requisitos Previos](#requisitos-previos)
4. [Instalación](#instalación)
5. [Ejecución del Proyecto](#ejecución-del-proyecto)
6. [Ejecución de Tests](#ejecución-de-tests)
7. [Herramientas de Calidad](#herramientas-de-calidad)
8. [Deployment y Contenedores](#deployment-y-contenedores)
9. [Pipeline CI/CD](#pipeline-cicd)
10. [Estructura del Proyecto](#estructura-del-proyecto)

---

## 📖 Descripción del Proyecto

Mini red social desarrollada con React (frontend) y Go (backend) que implementa:

- Registro y autenticación de usuarios
- Creación, visualización y eliminación de posts
- Sistema de comentarios en posts
- Validaciones de permisos (solo el autor puede eliminar su contenido)

**Stack Tecnológico:**
- **Backend:** Go 1.24 + PostgreSQL (Railway/Render)
- **Frontend:** React 18 + TypeScript
- **Testing:** Go testing + Jest + Cypress (107 tests: 89 unit + 18 handlers + integration + 15 E2E)
- **Containers:** Docker + GitHub Container Registry
- **Deployment:** Render (QA/PROD) + Railway PostgreSQL
- **Quality:** SonarCloud (47 issues fixed) + Code Coverage (86.5%/92.44%)
- **CI/CD:** GitHub Actions (calidad → contenedores → deploy)

---

## 🏗️ Arquitectura Integrada

### Capas de la Aplicación
```
Frontend (React)     →      Backend (Go)
──────────────────────     ────────────────────
React Components         ┌─►  Handlers     (HTTP handlers)
Axios Environment-aware │    ├── auth_handler.go
Auto-detect Backend URL  │    └── post_handler.go
                        │
                        │    Services     (business logic)
                        ├── auth_service.go     ───┐
                        └── post_service.go      ┌─┼───── Repository Interface
                              Validaciones       │ │       (mocks for testing)
                              Permisos           │ │       PostgreSQLUserRepository
                                                 │ │       PostgreSQLPostRepository
                        Repository               │ │
                        ├── user_repository.go ──┘ │
                        └── post_repository.go     │
                                                   │
PostgreSQL (Railway Cloud)     ←─── $1 placeholders + RETURNING
Railway QA / Railway PROD       ←─── Environment variables
```

### Ambiente QA vs PROD
```
QA (Auto-deploy)                          PROD (Manual approval)
─────────────────────────────────        ─────────────────────────────────
Frontend: render-qa.onrender.com         Frontend: render-prod.onrender.com
Backend:  back-qa.onrender.com           Backend:  back-prod.onrender.com
DB:       Railway pg-qa                  DB:       Railway pg-prod
Deploy:   GitHub Actions → auto          Deploy:   Manual approval → deploy
```

---

## 🔧 Requisitos Previos

### Software Necesario

```bash
# Verificar versiones instaladas:
go version    # Debe ser 1.24 o superior
node --version # Debe ser 20 o superior
npm --version  # Debe ser 10 o superior
```

### Instalación de Dependencias (si no las tenés)

**Go:**
```bash
# macOS
brew install go

# Ubuntu/Debian
sudo apt install golang-go

# Windows
# Descargar desde: https://go.dev/dl/
```

**Node.js y npm:**
```bash
# macOS
brew install node

# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Windows
# Descargar desde: https://nodejs.org/
```

---

## 📥 Instalación

### 1. Clonar el Repositorio

```bash
git clone https://github.com/OctavioCarpineti/IngSWIII-TP07-Quality.git
cd IngSWIII-TP07-Quality
```

### 2. Instalar Dependencias del Backend

```bash
cd backend
go mod download
cd ..
```

### 3. Instalar Dependencias del Frontend

```bash
cd frontend
npm install
cd ..
```

---

## 🚀 Ejecución del Proyecto

### Opción 1: Ejecución Manual (Recomendado para desarrollo)

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/api/main.go
```

El backend estará corriendo en `http://localhost:8080`

Deberías ver:
```
🚀 Servidor corriendo en http://localhost:8080
📊 Base de datos inicializada
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm start
```

El frontend estará corriendo en `http://localhost:3000`

Se abrirá automáticamente en tu navegador.

### Opción 2: Ejecución Local con Cypress (Automática)

Para ejecutar todo el proyecto localmente (incluyendo base de datos mock, backend, frontend) y abrir automáticamente la interfaz de Cypress para tests E2E:

```bash
# Desde la raíz del proyecto
# Asegúrate de tener Docker, Go, Node.js y npm instalados

# Ejecutar todo automáticamente (base de datos + backend + frontend + cypress)
./run-local.sh

# Si no tienes permisos de ejecución (en Linux/Mac) o usas git bash en Windows:
chmod +x run-local-db.sh run-local.sh
./run-local.sh

# En Windows CMD/PowerShell, ejecutar con bash:
bash run-local.sh
```

**Qué hace el script:**
- Inicia PostgreSQL local en Docker como base de datos mock (puerto 5432)
- Levanta el backend en `http://localhost:8080`
- Levanta el frontend en `http://localhost:3000`
- Abre automáticamente la interfaz de Cypress para ejecutar tests E2E
- Maneja la limpieza automática al presionar Ctrl+C

**Notas:**
- No modifica el setup de docker-compose ni las pipelines de CI/CD
- La base de datos "mock" es PostgreSQL en Docker para facilitar el desarrollo, sin datos iniciales complejos
- Si solo necesitas la base de datos, ejecuta `./run-local-db.sh`
- Para detener todo, presiona Ctrl+C en la terminal donde corre el script

### Opción 2: Ejecución con Scripts

**Backend:**
```bash
cd backend
# Compilar
go build -o app cmd/api/main.go

# Ejecutar
./app
```

**Frontend:**
```bash
cd frontend
# Build de producción
npm run build

# Servir build (requiere serve instalado: npm install -g serve)
serve -s build -l 3000
```

---

## 🧪 Ejecución de Tests

### Tests Unitarios - Backend Services

```bash
cd backend

# Ejecutar tests de servicios (35 tests)
go test ./tests/services/... -v

# Ejecutar tests con coverage
go test ./tests/services/... -v -cover -coverpkg=./internal/services/...

# Generar reporte HTML de coverage
go test ./tests/services/... -coverprofile=coverage.out -coverpkg=./internal/services/...
go tool cover -html=coverage.out

# Ver coverage en terminal
go tool cover -func=coverage.out
```

**Resultado esperado:**
```
=== RUN   TestRegister_Success
--- PASS: TestRegister_Success (0.00s)
...
PASS
coverage: 86.5% of statements in ./internal/services
ok      ingsw3-tp08/tests/services     0.537s
```

### Tests Unitarios - Backend Handlers

```bash
cd backend

# Ejecutar tests de handlers (18 tests - requieren mocks)
go test ./internal/handlers/... -v

# Con coverage
go test ./internal/handlers/... -v -cover
```

**Resultado esperado:**
```
=== RUN   TestAuthHandler_Register_Success
--- PASS: TestAuthHandler_Register_Success (0.00s)
...
PASS
coverage: 50.4% of statements
ok      ingsw3-tp08/internal/handlers   0.763s
```

### Tests de Integración - Repositories (Local Only)

```bash
cd backend

# Tests de repositorio (requiere Docker para Postgres container)
go test ./tests/integration/... -v

# Con coverage (cubre repositories ~85% + database setup)
go test ./tests/integration/... -v -cover
```

**Resultado esperado:**
```
=== RUN   TestUserRepositoryIntegrationTestSuite/TestCreate_Success
--- PASS: TestUserRepositoryIntegrationTestSuite/TestCreate_Success (2.15s)
...
PASS
ok      ingsw3-tp08/tests/integration     5.823s
```

> **Note:** Integration tests require Docker and are run locally. CI focuses on unit tests for faster feedback.

### Tests Combinados - Full Backend Coverage

```bash
cd backend

# Todos los tests unitarios + integración (requiere Docker)
go test ./tests/services/... ./internal/handlers/... ./tests/integration/... -v -cover -coverpkg=./...

# Ver coverage completo
go tool cover -func=combined.out
```

**Cobertura estimada después de mejoras:**
- **Services**: 86.5%
- **Handlers**: 50.4%
- **Repositories**: ~85% (con integración completa)
- **Total Backend**: ~75-80%

### Tests Unitarios - Frontend

```bash
cd frontend

# Ejecutar tests en modo watch
npm test

# Ejecutar tests una vez
npm test -- --watchAll=false

# Ejecutar tests con coverage
npm test -- --coverage --watchAll=false

# Ver reporte de coverage en navegador
open coverage/lcov-report/index.html
```

**Resultado esperado:**
```
Test Suites: 8 passed, 8 total
Tests:       39 passed, 39 total
Coverage:    92.44% statements
```

### Tests E2E - Cypress

**Prerequisito: Backend y Frontend deben estar corriendo**

```bash
# Terminal 1: Backend
cd backend
go run cmd/api/main.go

# Terminal 2: Frontend  
cd frontend
npm start

# Terminal 3: Cypress
cd frontend

# Modo interactivo (recomendado)
npx cypress open
# Luego click en "E2E Testing" y seleccionar los tests

# Modo headless (para CI/CD)
npx cypress run
```

**Resultado esperado:**
```
Running:  auth.cy.js                    (1 of 4)
  ✓ 5 tests passing

Running:  posts.cy.js                   (2 of 4)
  ✓ 5 tests passing

Running:  comments.cy.js                (3 of 4)
  ✓ 4 tests passing

Running:  full-flow.cy.js               (4 of 4)
  ✓ 1 test passing

Total: 15 tests passing
```

---

## 🔍 Herramientas de Calidad

### 1. SonarCloud (Análisis Estático)

**Acceso al proyecto:**
```
URL: https://sonarcloud.io/project/overview?id=OctavioCarpineti_IngSWIII-TP07-Quality
Organization: octaviocarpineti
```

**Análisis local (opcional):**
```bash
# Requiere configuración de SONAR_TOKEN
docker run --rm \
  -e SONAR_HOST_URL="https://sonarcloud.io" \
  -e SONAR_TOKEN="tu-token" \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli
```

### 2. Code Coverage

**Backend:**
```bash
cd backend
go test ./tests/services/... -coverprofile=coverage.out -coverpkg=./internal/services/...

# Ver en terminal
go tool cover -func=coverage.out | grep total

# Ver en navegador
go tool cover -html=coverage.out
```

**Frontend:**
```bash
cd frontend
npm test -- --coverage --watchAll=false

# Abrir reporte HTML
open coverage/lcov-report/index.html
```

---

## 🚢 Deployment y Contenedores

### Desarrollo Local con Docker Compose

Para ejecutar todo el stack localmente:

```bash
# Ejecutar con docker-compose
docker-compose up --build

# Servicios disponibles:
# - PostgreSQL: localhost:5432
# - Backend:    localhost:8080
# - Frontend:   localhost:3000
```

### Contenedores Individuales

**Backend:**
```bash
cd backend
docker build -t ingsw3-integrated-backend .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgresql://..." \
  ingsw3-integrated-backend
```

**Frontend:**
```bash
cd frontend
docker build -t ingsw3-integrated-frontend .
docker run -p 3000:80 \
  -e REACT_APP_BACKEND_URL="http://localhost:8080" \
  ingsw3-integrated-frontend
```

### Deployment en Producción

#### Registros Necesarios:

**Railway (Base de datos PostgreSQL):**
1. Crear proyecto QA: `ingsw3-integrated-qa`
2. Crear proyecto PROD: `ingsw3-integrated-prod`
3. Copiar `DATABASE_URL` de cada uno

**Render (Aplicación):**
1. Crear servicio web QA backend
2. Crear servicio static site QA frontend
3. Repetir para PROD
4. Configurar environment variables:
   - Backend: `DATABASE_URL`, `PORT`
   - Frontend: `REACT_APP_BACKEND_URL`

#### GitHub Secrets Requeridos:
```
RENDER_QA_BACK_ID     # ID del servicio QA backend en Render
RENDER_QA_FRONT_ID    # ID del servicio QA frontend en Render
RENDER_PROD_BACK_ID   # ID del servicio PROD backend en Render
RENDER_PROD_FRONT_ID  # ID del servicio PROD frontend en Render
RENDER_API_KEY        # API key de Render para deployments
SONAR_TOKEN           # Para SonarCloud analysis
```

### Arquitectura de Deploy

```
Git Push
   ↓
GitHub Actions
   ↓ Quality Gates (TP7)
   ↓  Backend/Frontend Tests
   ↓  Coverage ≥70%
   ↓  SonarCloud Pass
   ↓  Cypress E2E
   ↓
Docker Build (TP8)
   ↓ Push to GHCR
   ↓
Deploy QA (Auto)
   ↓
Manual Approval
   ↓
Deploy PROD (TP8)
```

---

## 🔄 Pipeline CI/CD

### GitHub Actions

El pipeline integrado ejecuta automáticamente en cada push y combina TP7 + TP8:

**Fases del Pipeline:**
1. 🔍 **Calidad (TP7):** Tests unitarios, coverage, SonarCloud, E2E
2. 🐳 **Contenedores (TP8):** Docker build + push to GHCR
3. 🚀 **Despliegue QA:** Deploy automático a Render QA
4. ✋ **Aprobación PROD:** Espera aprobación manual
5. 🎯 **Despliegue PROD:** Deploy final a producción

**Quality Gates Configurados:**
- ❌ Backend coverage < 70% (86.5% alcanzado)
- ❌ Frontend coverage < 70% (92.44% alcanzado)
- ❌ SonarCloud Quality Gate falla (PASSED)
- ❌ Tests unitarios fallan (35 back + 39 front)
- ❌ Tests E2E fallan (15 Cypress)
- ❌ Builds de contenedores fallan

**Ver estado del pipeline:**
```
GitHub > Actions > CI/CD Pipeline
```

**Ejecutar pipeline manualmente:**
```bash
git commit --allow-empty -m "trigger pipeline"
git push
```

---

## 📁 Estructura del Proyecto

```
tp07-quality/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Entry point del servidor
│   ├── internal/
│   │   ├── handlers/                # HTTP handlers (50.4% coverage)
│   │   │   ├── auth_handler.go
│   │   │   ├── post_handler.go
│   │   │   ├── auth_handler_test.go   # 6 tests handler unitarios
│   │   │   └── post_handler_test.go   # 12 tests handler unitarios
│   │   ├── services/                # Lógica de negocio (86.5% coverage)
│   │   │   ├── auth_service.go
│   │   │   └── post_service.go
│   │   ├── repository/              # Acceso a datos (~85% coverage)
│   │   │   ├── user_repository.go
│   │   │   └── post_repository.go
│   │   ├── models/                  # Estructuras de datos
│   │   │   ├── users.go
│   │   │   └── post.go
│   │   ├── database/                # Configuración BD
│   │   │   └── database.go
│   │   └── router/                  # Configuración de rutas
│   │       └── router.go
│   ├── tests/
│   │   ├── services/                # 35 tests unitarios + covers services 86.5%
│   │   │   ├── auth_service_test.go
│   │   │   └── post_service_test.go
│   │   ├── mocks/                   # Mocks para testing
│   │   │   ├── mock_user_repository.go
│   │   │   ├── mock_post_repository.go
│   │   │   ├── auth_service_mock.go    # Mock de AuthService
│   │   │   └── post_service_mock.go    # Mock de PostService
│   │   └── integration/             # Tests de integración con DB real
│   │       ├── test_helpers.go         # Setup Postgres container
│   │       └── user_repository_integration_test.go  # ~6 tests repo integración
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── components/              # Componentes React (92.44% coverage)
│   │   │   ├── Login/
│   │   │   │   ├── Login.tsx
│   │   │   │   ├── Login.test.tsx
│   │   │   │   └── Login.css
│   │   │   ├── PostList/
│   │   │   │   ├── PostList.tsx
│   │   │   │   ├── PostList.test.tsx
│   │   │   │   └── PostList.css
│   │   │   ├── CreatePost/
│   │   │   ├── PostDetail/
│   │   │   ├── CommentList/
│   │   │   └── CommentForm/
│   │   ├── services/                # Servicios HTTP
│   │   │   ├── authService.ts
│   │   │   ├── authService.test.ts
│   │   │   ├── postService.ts
│   │   │   └── postService.test.ts
│   │   ├── types/                   # Definiciones TypeScript
│   │   │   └── index.ts
│   │   ├── App.tsx                  # Componente principal
│   │   └── index.tsx                # Entry point
│   ├── cypress/
│   │   ├── e2e/
│   │   │   └── blog/                # 15 tests E2E
│   │   │       ├── auth.cy.js       # 5 tests
│   │   │       ├── posts.cy.js      # 5 tests
│   │   │       ├── comments.cy.js   # 4 tests
│   │   │       └── full-flow.cy.js  # 1 test
│   │   └── support/
│   │       ├── e2e.js
│   │       └── commands.js
│   ├── cypress.config.js
│   ├── package.json
│   └── package-lock.json
│
├── .github/
│   └── workflows/
│       └── ci.yml                   # Pipeline CI/CD
│
├── run-local-db.sh                 # Script para iniciar DB local
├── run-local.sh                    # Script para iniciar todo local + Cypress
├── sonar-project.properties         # Configuración SonarCloud
├── README.md                        # Este archivo
└── decisiones.md                    # Decisiones técnicas y justificaciones
```

---

## 🐛 Troubleshooting

### Backend no inicia

```bash
# Verificar puerto 8080 disponible
lsof -i :8080
# Si está ocupado, matar el proceso:
kill -9 <PID>

# Verificar Go instalado correctamente
go version

# Limpiar y reinstalar dependencias
cd backend
rm go.sum
go mod tidy
go mod download
```

### Frontend no inicia

```bash
# Verificar puerto 3000 disponible
lsof -i :3000

# Limpiar cache y reinstalar
cd frontend
rm -rf node_modules package-lock.json
npm install

# Si falla con errores de Cypress
npm install --save-dev cypress@13.15.2
```

### Tests de Cypress fallan

```bash
# Verificar que backend y frontend estén corriendo
curl http://localhost:8080/api/health
curl http://localhost:3000

# Limpiar cache de Cypress
npx cypress cache clear
npx cypress install

# Ejecutar con logs detallados
DEBUG=cypress:* npx cypress run
```

### Pipeline falla en GitHub Actions

```bash
# Verificar logs en:
# GitHub > Actions > Click en el run fallido

# Causas comunes:
# 1. package-lock.json desincronizado
cd frontend
rm package-lock.json
npm install
git add package-lock.json
git commit -m "fix: regenerar package-lock.json"
git push

# 2. Tests fallan localmente primero
# Ejecutar todos los tests localmente antes de push
```

---

## 📊 Métricas Alcanzadas - Integración TP7 + TP8

### Quality Assurance (TP7)
| Métrica | Objetivo | Resultado | Estado |
|---------|----------|-----------|--------|
| Backend Coverage | ≥70% | 86.5% | ✅ |
| Frontend Coverage | ≥70% | 92.44% | ✅ |
| Tests Unitarios | - | 74 tests | ✅ |
| Tests E2E Cypress | - | 15 tests | ✅ |
| **Total Tests** | - | **89 tests** | ✅ |
| SonarCloud Quality Gate | Pass | PASSED | ✅ |
| Issues Code Smells Resueltos | ≥3 | 47 issues | ✅ |
| Duplications | <3% | 0.0% | ✅ |

### Deployment & Contenedores (TP8)
| Aspecto | Implementación | Estado |
|---------|---------------|--------|
| Base de Datos | PostgreSQL (Railway QA/PROD) | ✅ |
| Backend Container | Go + multi-stage Docker | ✅ |
| Frontend Container | React + multi-stage Docker | ✅ |
| Container Registry | GitHub Container Registry | ✅ |
| CI/CD Integration | Docker build + push en pipeline | ✅ |
| Deploy QA | Render auto-deploy | ✅ |
| Deploy PROD | Render manual approval | ✅ |
| Environment Config | Variables QA vs PROD separadas | ✅ |

### Arquitectura Integrada
- ✅ **16 archivos modificados** para compatibilidad PostgreSQL
- ✅ **Frontend environment-aware** (auto-detecta backend URLs)
- ✅ **Pipeline fusionado**: calidad → contenedores → deploy
- ✅ **89 tests automatizados** manteniendo cobertura alta
- ✅ **47 issues SonarCloud** resueltos (constantes, duplicaciones)
- ✅ **3 ambientes**: desarrollo local, QA, producción

---



**Alumno:** Octavio Carpineti - Kevin Massholder 
**GitHub:** https://github.com/OctavioCarpineti  
**Repositorio:** https://github.com/OctavioCarpineti/IngSWIII-TP07-Quality
