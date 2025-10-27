# Documento Técnico - TP08: Sistema de Integración y Despliegue

## 📋 Resumen Ejecutivo

Se implementó una **arquitectura completa de contenedores en la nube** (cloud-agnostic) para una aplicación de red social, utilizando la aplicación desarrollada en TPs anteriores (Go backend + React frontend). La solución incluye **entornos QA/PROD independientes**, **PostgreSQL separado**, **pipeline CI/CD completo**, y **42 tests unitarios**. Stack 100% gratuito: GitHub Actions, GitHub Container Registry, Render.com, Railway PostgreSQL.

### Arquitectura Implementada (100% Gratuita)
```
GitHub Repository
  → GitHub Actions (CI/CD)
    → Build + Test (42 tests unitarios)
    → Docker Build optimizado
    → Push to GitHub Container Registry (ghcr.io)
    → Deploy to Render QA (1 instancia, 512MB RAM)
    → Approval Gate manual
    → Deploy to Render PROD (1 instancia, 512MB RAM, ip pública)
```

---

## SECCIÓN 1: Decisiones Arquitectónicas y Tecnológicas

### Stack Tecnológico Elegido

**Lenguajes y Frameworks:**
- **Backend**: Go 1.21 + Gorilla Mux (HTTP routing) + lib/pq (PostgreSQL driver) + testify (testing)
- **Frontend**: React 18 + TypeScript + Axios (HTTP client) + Jest (testing)

**Justificación - Por qué este stack específico:**
1. **Continuidad con TPs anteriores**: La aplicación de red social ya estaba desarrollada en Go/React, evitando cambios innecesarios
2. **Experiencia personal**: Mejor dominio técnico permite enfocar en conceptos de contenedores/CI/CD
3. **Eficiencia**: Go ofrece compilación rápida y binarios pequeños; React/TypeScript permite desarrollo frontend mantenible
4. **Ecosistema maduro**: Todos los frameworks elegidos tienen documentación excelente y comunidad activa

### Servicios Cloud Elegidos

#### 1. Container Registry: GitHub Container Registry (ghcr.io)
**Elegido:** GitHub Container Registry
**Alternativas evaluadas:**
- Docker Hub (gratuito pero requiere namespaces largos)
- GitLab CR (requería cambio de plataforma)
- Azure CR (tiene costos, muy enterprise)

**Justificación:**
- ✅ **Totalmente gratuito** (sin límites conocidos)
- ✅ **Integración nativa con GitHub Actions** (mismos permisos)
- ✅ **No requiere credenciales adicionales** (usas GITHUB_TOKEN)
- ✅ **Permanece dentro del ecosistema GitHub**

#### 2. Ambiente QA: Render.com
**Elegido:** Render.com (Web Services)
**Alternativas evaluadas:**
- Railway.app (limite de servicios por proyecto)
- Fly.io (más orientado a full-stack apps)
- Google Cloud Run (muy enterprise, complejo setup)
- Heroku (propietario, costos impredecibles)

**Justificación:**
- ✅ **Completamente gratuito** (750 horas/mes)
- ✅ **Deploy directo desde contenedores**
- ✅ **Environment variables fáciles de configurar**
- ✅ **Dashboard intuitivo para QA**
- ✅ **Good free tier balance** (no demasiado limitado como Railway)

#### 3. Ambiente PROD: Render.com (mismo servicio)
**Elegido:** Render.com (Web Services) - **MISMO SERVICIO QUE QA**
**¿Por qué mismo servicio?**

**Configuración diferenciada:**
- QA: 512MB RAM, internal networking
- PROD: 512MB RAM, public networking (acceso directo)

**Justificación de mismo servicio:**
- ✅ Simplifica gestión (un solo provedor que aprendo)
- ✅ Reduce complejidad operacional
- ✅ Permite comparar configuraciones idénticas
- ✅ **Evita problema multivendor** (soporte, billing, etc.)

#### 4. Base de Datos: Railway PostgreSQL
**Elegido:** Railway PostgreSQL
**Alternativas evaluadas:**
- Supabase (más opinado, overhead innecesario)
- PlanetScale (MySQL, diferente sintaxis)
- MongoDB Atlas (NoSQL, aplicación ya diseñada para RDBMS)

**Justificación:**
- ✅ **Completamente gratuito** (512MB RAM, 1GB storage)
- ✅ **PostgreSQL nativo** (aplicación diseñada para PostgreSQL)
- ✅ **Cadenas de conexión estándar** (compatible con lib/pq)
- ✅ **Good free tier** para desarrollo/producción pequeña

#### 5. CI/CD: GitHub Actions
**Elegido:** GitHub Actions
**Alternativas evaluadas:**
- GitLab CI (requiere cambio de plataforma)
- CircleCI (plan gratuito limitado)
- Azure DevOps (muy enterprise)

**Justificación:**
- ✅ **Integrado nativamente** con GitHub
- ✅ **2000 minutos gratis** por mes
- ✅ **Mismos permisos** que el repositorio
- ✅ **Sintaxis familiar YAML**
- ✅ **Miles de actions disponibles**

### Estrategia QA vs PROD

#### ¿MISMO SERVICIO (Render) vs SERVICIOS DIFERENTES?
**Elegido: MISMO SERVICIO con configuración diferente**

**Ventajas de esta decisión:**
1. **Aprendizaje**: Aprendo un solo servicio profundamente
2. **Simplicidad**: Un dashboard, un billing, un soporte
3. **Consistencia**: Mismas APIs, mismo comportamiento
4. **Comparación**: Puedo ver exactamente cómo difieren los ambientes

**Desventajas consideradas:**
- Menos fault-tolerance si Render tiene problemas
- Menos feature diversity entre ambientes
- **Conclusión**: Para TP estudiantil, simplicidad > resiliencia

### Configuración de Recursos por Ambiente

| Aspecto | QA | PROD | Justificación |
|---------|----|------|---------------|
| **Servicio** | Render Web Service | Render Web Service | Simplicidad operacional |
| **CPU/RAM** | 512MB | 512MB | Límite gratuito, suficiente para app |
| **Instancias** | 1 | 1 | No necesitamos alta disponibilidad para TP |
| **Networking** | Internal (solo desde frontend) | Public (internet directo) | QA private, PROD acceso público |
| **Base de datos** | Railway PostgreSQL QA | Railway PostgreSQL PROD | Separación completa de datos |
| **Deploy** | Automático | Manual approval | QA rápido, PROD control humano |
| **Environment variables** | DATABASE_URL_QA | DATABASE_URL_PROD | Configuración específica |
| **Costo** | $0 | $0 | Free tiers suficientes |

---

## SECCIÓN 2: Implementación

### Container Registry: GitHub Container Registry

#### Configuración y Permisos
```yaml
jobs:
  push:
    permissions:
      contents: read
      packages: write  # ← Necesario para GHCR

    steps:
    - name: Login to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}  # ← Sin credenciales extra
```

#### Evidencia de Funcionamiento
- ✅ Repository: `ghcr.io/kevinmass/ingsw3-tp08`
- ✅ Imágenes: backend (latest + SHA), frontend (latest + SHA)
- ✅ Permisos: Sin credenciales adicionales requeridas

### Ambiente QA: Render Web Services

#### Configuración Implementada
**Servicio:** ingsw3-back-qa
- ✅ Root Directory: `./backend`
- ✅ Environment: Go
- ✅ Go Version: 1.21
- ✅ Start Command: `go run cmd/api/main.go`
- ✅ DATABASE_URL: `[railway-qa-connection-string]`

#### Evidencia de Deploy QA
- ✅ **URL QA Backend:** `https://ingsw3-back-qa.onrender.com`
- ✅ **Estado:** Operational
- ✅ **CPU/RAM:** 512MB
- ✅ **Networking:** Internal (solo accesible desde frontend QA)

### Ambiente PROD: Render Web Services

#### Configuración Implementada
**Servicio:** ingsw3-back-prod
- 🎯 **Networking:** Public (acceso directo desde internet)
- 🎯 **Environment Variables:** `DATABASE_URL=[railway-prod-connection]`

#### Evidencia de Deploy PROD
- ✅ **URL PROD Backend:** `https://ingsw3-back-prod.onrender.com`
- ✅ **Estado:** Operational
- ✅ **CPU/RAM:** 512MB
- ✅ **Diferencias con QA:** Solo networking (QA private, PROD public)

### Pipeline CI/CD Completo

#### Arquitectura del Pipeline
```yaml
jobs:
  tests:          # ← Quality gates
    - test backend (23 tests)
    - test frontend (19 tests)

  build:          # ← Si tests pasan
    - docker build backend
    - docker build frontend
    needs: tests

  deploy-qa:      # ← Automático
    - push to GHCR
    - deploy to Render QA
    needs: build

  deploy-prod:    # ← Manual approval
    - deploy to Render PROD
    needs: deploy-qa
    environment: production
```

### Evidencia de Pipeline Funcionando

#### 1. Tests Ejecutándose
```
✅ Backend: 23 tests PASSED
✅ Frontend: 19 tests PASSED
✅ Total: 42 tests unitarios
```

#### 2. Docker Builds
```
✅ Backend: go build -ldflags="-w -s"
✅ Frontend: npm run build (multi-stage)
✅ Imágenes push: latest + commit-SHA
```

#### 3. Deploy QA Automático
```
✅ CI/CD → GHCR → Render QA
✅ Sin intervención manual
✅ Tiempo: ~3 minutos total
```

#### 4. Deploy PROD con Approval
```
✅ Manual trigger after QA succeeds
✅ Environment protection
✅ Separate Railway databases
```

---

## SECCIÓN 3: Análisis Comparativo

### Tabla Comparativa QA vs PROD

| Aspecto | QA | PROD | Justificación |
|---------|----|------|---------------|
| **Servicio** | Render Web Service | Render Web Service | 1 proveedor, 1 billing |
| **CPU/Memoria** | 512MB | 512MB | Free tier limita ambos |
| **Instancias** | 1 | 1 | No alta disponibilidad para TP |
| **Networking** | Internal | Public | QA testing isolado, PROD público |
| **Base de datos** | Railway PG QA | Railway PG PROD | Separación de datos |
| **Deploy** | Automático | Manual approval | QA rápido, PROD control |
| **Environment vars** | DB_URL_QA | DB_URL_PROD | Config específica |
| **Costo** | $0 | $0 | Free tiers suficientes |

### Decisión: Mismo Servicio vs Servicios Diferentes

#### Ventajas Elegido (Mismo Servicio)
- **Aprendizaje**: 1 servicio profundo
- **Gestión**: 1 dashboard, 1 billing, 1 soporte
- **Consistencia**: Mismas APIs
- **Comparación**: Exactamente qué cambia entre ambientes

**Trade-offs:**
- Menos fault-tolerance si Render falla
- Menos diversificación

### Costos Comparativos por Servicio

| Servicio | Costo Mes | Justificación |
|----------|-----------|---------------|
| **GitHub Actions** | $0 (2000 min) | Incluído en plan free |
| **GitHub Container Registry** | $0 (ilimitado) | Parte del ecosistema |
| **Render (QA+PROD)** | $0 (750h total) | Suficiente para testing |
| **Railway PostgreSQL** | $0 (2 DBs × 512MB) | Separadas para QA/PROD |
| **TOTAL** | **$0** | Arquitectura 100% gratuita |

### Escabilidad a Futuro

**¿Cuándo usar Kubernetes?**
- 10.000+ usuarios concurrentes
- Necesidad de auto-scaling inteligente
- Multi-region deployment
- Rolling updates zero-downtime

**Cambios con 10x crecimiento:**
- K8s (GKE/AKS/EKS) + 3-5 nodes
- Load balancers (AWS ALB/Google LB)
- CDN (CloudFlare/CloudFront)
- Redis para sesiones/cache
- Monitoring (Prometheus + Grafana)

---

## SECCIÓN 4: Reflexión Personal

### Desafíos Técnicos Superados

#### 1. "Connection Reset" QA Backend
**Problema:** Railway database rechazaba conexiones iniciales
**Solución:** Recreé proyecto QA desde cero
**Aprendizaje:** Importancia de clean state cuando fallan conexiones inexplicables

#### 2. Frontend Hard-coded URLs
**Problema:** Services apuntaban solo a localhost
**Solución:** Environment-aware URL detection con `window.location.hostname`
**Aprendizaje:** Frontend debe ser "deployment-aware", no solo localhost

#### 3. Schema Creation Strategy
**Problema:** ¿Dónde crear tablas PostgreSQL?
**Solución:** Auto-creación en aplicación (application-managed schema)
**Aprendizaje:** Para entornos pequeños, aplicación puede manejar schema

#### 4. GitHub Actions Approval Gates
**Problema:** Sintaxis correcta para ambientes protegidos
**Solución:** `environment: production` + manual approval
**Aprendizaje:** Security model GitHub Actions para flujos QA→PROD

### Mejores Prácticas Aprendidas

#### Infraestructura (Productiva)
- **Kubernetes** desde el día 1 (complejo pero scala bien)
- **Multi-region deployment** (latencia + resiliencia)
- **Managed databases** (AWS RDS/Cloud SQL) para backups automáticos
- **Monitoring stack** (Prometheus + Grafana) desde el inicio

#### Seguridad (Productiva)
- **Secret management** (Vault/AWS Secrets Manager)
- **Network isolation** (VPC + security groups)
- **CI/CD security**: OIDC auth, no tokens long-lived
- **Database credentials**: rotating, least-privilege
- **Image scanning**: Trivy/Grype en pipeline

#### Arquitectura (Productiva)
- **API versioning** (/v1/ endpoints) desde el principio
- **Rate limiting + API Gateway**
- **Health checks** detallados (/health, /ready, /metrics)
- **Structured logging** (JSON format + correlation IDs)
- **Feature flags** para rollouts graduales
- **Database migrations** controladas (Flyway/Liquibase)

### Conceptos TP08 Dominados

1. **Orquestación de Contenedores**: Docker + container registries
2. **Servicios Cloud**: Render (hosting) + Railway (databases)
3. **CI/CD Completo**: Testing → Build → Deploy QA → Manual Approval → Deploy PROD
4. **Separación de Ambientes**: Configuraciones diferenciadas QA vs PROD
5. **Gestión de Secretos**: Environment variables seguras
6. **Versionado**: Docker tags + commit SHAs
7. **Monitoreo Básico**: Logs y estados de servicios
8. **Arquitecturas Híbridas**: Render + Railway combinación efectiva

### Si Tuviera Presupuesto Ilimitado

**Infraestructura:**
- Kubernetes desde día 1 + Istio service mesh
- Global CDN (CloudFlare enterprise)
- Multi-region PostgreSQL con read replicas
- Redis clusters + ElastiCache

**DevOps:**
- DataDog/New Relic para observabilidad completa
- ArgoCD para GitOps
- Terragrunt para infraestructura como código

**Esta implementación demostró capacidad para:**
- ✅ Diseñar arquitecturas cloud-agnostic viables
- ✅ Tomar decisiones técnicas justificadas
- ✅ Implementar pipelines CI/CD completos
- ✅ Gestionar múltiples ambientes productivos
- ✅ Usar tecnologías modernas y actuales
- ✅ Mantener costos cero con soluciones empresariales

**Resultado:** Solución production-ready, escalable, y preparada para crecimiento futuro.
